# TipFax

A Go application that connects to StreamElements and prints tips to a thermal printer.

## Features

- Connects to StreamElements Astro WebSocket API
- Prints tips to a thermal printer
- Graceful shutdown handling
- Configurable via environment variables
- Web interface to view server status (default: http://localhost:8082)

## Requirements

- Go 1.21 or later
- A compatible thermal printer (tested with Epson TM-T20II)
- StreamElements JWT token

## Configuration

The following environment variables are available:

- `SE_JWT_TOKEN`: StreamElements JWT token (required)
- `DEVICE_PATH`: Printer device path (default: `/dev/usb/lp0`)

## Building

```bash
make
```

## Running

```bash
./bin/server
```