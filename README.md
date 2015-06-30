# Prometheus Aurora exporter

This is an exporter for Prometheus to get instrumentation data for [Apache Aurora](http://aurora.apache.org/)

## Build and run

    make
    ./aurora_exporter <flags>

### Flags

Name                           | Description
-------------------------------|------------
web.listen-address             | Address to listen on for web interface and telemetry.
web.telemetry-path             | Path under which to expose metrics.
exporter.aurora-url            | URL to a host running the Aurora scheduler

---
