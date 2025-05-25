package web

import (
	"html/template"
	"net/http"

	"github.com/DaniruKun/tipfax/internal/config"
)

type PageData struct {
	ServerStatus string
	DevicePath   string
	JWTPreview   string
}

func StatusHandler(cfg *config.Config, devicePath string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get first few characters of JWT token
		jwtPreview := "Not configured"
		if cfg.SeJWTToken != "" {
			if len(cfg.SeJWTToken) > 20 {
				jwtPreview = cfg.SeJWTToken[:20] + "..."
			} else {
				jwtPreview = cfg.SeJWTToken
			}
		}

		data := PageData{
			ServerStatus: "Running",
			DevicePath:   devicePath,
			JWTPreview:   jwtPreview,
		}

		tmpl := `<!DOCTYPE html>
<html>
<head>
    <title>TipFax Server Status</title>
    <style>
        body {
            font-family: Arial, sans-serif;
            max-width: 800px;
            margin: 50px auto;
            padding: 20px;
            background-color: #f5f5f5;
        }
        .container {
            background-color: white;
            padding: 30px;
            border-radius: 10px;
            box-shadow: 0 2px 10px rgba(0,0,0,0.1);
        }
        h1 {
            color: #333;
            text-align: center;
            margin-bottom: 30px;
        }
        .status-item {
            margin: 15px 0;
            padding: 15px;
            background-color: #f8f9fa;
            border-left: 4px solid #007bff;
            border-radius: 4px;
        }
        .status-label {
            font-weight: bold;
            color: #495057;
        }
        .status-value {
            color: #28a745;
            font-family: monospace;
        }
        .footer {
            text-align: center;
            margin-top: 30px;
            color: #6c757d;
            font-size: 14px;
        }
    </style>
</head>
<body>
    <div class="container">
        <h1>TipFax Server Status</h1>
        
        <div class="status-item">
            <div class="status-label">Server Status:</div>
            <div class="status-value">{{.ServerStatus}}</div>
        </div>
        
        <div class="status-item">
            <div class="status-label">Configured Device Path:</div>
            <div class="status-value">{{.DevicePath}}</div>
        </div>
        
        <div class="status-item">
            <div class="status-label">JWT Token Preview:</div>
            <div class="status-value">{{.JWTPreview}}</div>
        </div>
        
        <div class="footer">
            TipFax Server - StreamElements Tip Printer
        </div>
    </div>
</body>
</html>`

		t, err := template.New("status").Parse(tmpl)
		if err != nil {
			http.Error(w, "Template error", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "text/html")
		t.Execute(w, data)
	}
}
