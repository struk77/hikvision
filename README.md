# Overview
This repository contents Prometheus Exporter for Hikvision Cameras.

# Building
```bash
go build
```

# Usage
```bash
./hikvision_exporter --cameras=cameras.yml --listen=":19101" --period=60
```
