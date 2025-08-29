# NetPulse

A lightweight Linux system tray application that monitors network interface speeds in real-time.

## Features

- Real-time network speed monitoring in system tray
- Support for all network interfaces
- Interface selection (All interfaces or specific interface)
- Automatic unit scaling (B/s, KB/s, MB/s)
- Clean, minimal interface

## Usage

```bash
go build -o netpulse main.go
./run.sh
```

## Requirements

- Linux with GTK3
- Go 1.24+
- System tray support

## Interface Selection

Right-click the tray icon to select:
- All Interfaces (combined stats)
- Individual interfaces (eth0, wlan0, etc.)