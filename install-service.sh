#!/bin/bash

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Configuration
SERVICE_NAME="tipfax"
SERVICE_USER="tipfax"
SERVICE_GROUP="tipfax"
INSTALL_DIR="/opt/tipfax"
BINARY_NAME="tipfax-server"

echo -e "${GREEN}TipFax Service Installation Script${NC}"
echo "=================================="

# Check if running as root
if [[ $EUID -ne 0 ]]; then
   echo -e "${RED}This script must be run as root (use sudo)${NC}"
   exit 1
fi

# Check if SE_JWT_TOKEN is provided
if [[ -z "$SE_JWT_TOKEN" ]]; then
    echo -e "${YELLOW}Warning: SE_JWT_TOKEN environment variable not set${NC}"
    echo "You'll need to set it manually in the service file later"
fi

echo "Step 1: Creating service user and group..."
if ! id "$SERVICE_USER" &>/dev/null; then
    useradd --system --no-create-home --shell /bin/false "$SERVICE_USER"
    echo -e "${GREEN}Created user: $SERVICE_USER${NC}"
else
    echo "User $SERVICE_USER already exists"
fi

echo "Step 2: Creating installation directory..."
mkdir -p "$INSTALL_DIR"
chown "$SERVICE_USER:$SERVICE_GROUP" "$INSTALL_DIR"
chmod 755 "$INSTALL_DIR"

echo "Step 3: Installing the application binary..."
if [[ ! -f "go.mod" ]]; then
    echo -e "${RED}Error: go.mod not found. Please run this script from the project root directory.${NC}"
    exit 1
fi

# Check if binary already exists (pre-built)
if [[ -f "$BINARY_NAME" ]]; then
    echo -e "${GREEN}Found pre-built binary: $BINARY_NAME${NC}"
    cp "$BINARY_NAME" "$INSTALL_DIR/$BINARY_NAME"
else
    echo "Building the application..."
    # Check if Go is available
    if ! command -v go &> /dev/null; then
        echo -e "${YELLOW}Go not found in PATH. Trying common Go installation paths...${NC}"
        
        # Common Go installation paths
        GO_PATHS=(
            "/usr/local/go/bin/go"
            "/opt/go/bin/go"
            "$HOME/go/bin/go"
            "/snap/bin/go"
						"$HOME/.nix-profile/bin/go"
        )
        
        GO_CMD=""
        for path in "${GO_PATHS[@]}"; do
            if [[ -x "$path" ]]; then
                GO_CMD="$path"
                echo -e "${GREEN}Found Go at: $GO_CMD${NC}"
                break
            fi
        done
        
        if [[ -z "$GO_CMD" ]]; then
            echo -e "${RED}Error: Go compiler not found and no pre-built binary available.${NC}"
            echo "Please either:"
            echo "1. Build the binary first: go build -o $BINARY_NAME ./cmd/server"
            echo "2. Install Go and ensure it's in your PATH"
            exit 1
        fi
    else
        GO_CMD="go"
    fi
    
    $GO_CMD build -o "$INSTALL_DIR/$BINARY_NAME" ./cmd/server
fi
chown "$SERVICE_USER:$SERVICE_GROUP" "$INSTALL_DIR/$BINARY_NAME"
chmod 755 "$INSTALL_DIR/$BINARY_NAME"
echo -e "${GREEN}Built and installed binary to $INSTALL_DIR/$BINARY_NAME${NC}"

echo "Step 4: Installing systemd service..."
if [[ ! -f "tipfax.service" ]]; then
    echo -e "${RED}Error: tipfax.service file not found in current directory${NC}"
    exit 1
fi

# Update SE_JWT_TOKEN in service file if provided
if [[ -n "$SE_JWT_TOKEN" ]]; then
    sed -i "s/Environment=SE_JWT_TOKEN=/Environment=SE_JWT_TOKEN=$SE_JWT_TOKEN/" tipfax.service
    echo -e "${GREEN}Updated SE_JWT_TOKEN in service file${NC}"
fi

cp tipfax.service /etc/systemd/system/
chmod 644 /etc/systemd/system/tipfax.service
echo -e "${GREEN}Installed service file to /etc/systemd/system/tipfax.service${NC}"

echo "Step 5: Adding user to lp group for printer access..."
usermod -a -G lp "$SERVICE_USER"

echo "Step 6: Reloading systemd and enabling service..."
systemctl daemon-reload
systemctl enable "$SERVICE_NAME"
echo -e "${GREEN}Service enabled${NC}"

echo ""
echo -e "${GREEN}Installation completed successfully!${NC}"
echo ""
echo "Next steps:"
echo "1. If you didn't set SE_JWT_TOKEN, edit the service file:"
echo "   sudo systemctl edit tipfax --full"
echo "   And update the Environment=SE_JWT_TOKEN= line"
echo ""
echo "2. Start the service:"
echo "   sudo systemctl start tipfax"
echo ""
echo "3. Check service status:"
echo "   sudo systemctl status tipfax"
echo ""
echo "4. View logs:"
echo "   sudo journalctl -u tipfax -f"
echo ""
echo "5. To uninstall:"
echo "   sudo systemctl stop tipfax"
echo "   sudo systemctl disable tipfax"
echo "   sudo rm /etc/systemd/system/tipfax.service"
echo "   sudo rm -rf $INSTALL_DIR"
echo "   sudo userdel $SERVICE_USER" 