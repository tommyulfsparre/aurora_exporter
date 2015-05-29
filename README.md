# Prometheus Aurora exporter

This is an exporter for Prometheus to get instrumentation data for [Apache Aurora](http://aurora.apache.org/)

## Run

    ./aurora_exporter <flags>

### Flags

Name                           | Description
-------------------------------|------------
web.listen-address             | Address to listen on for web interface and telemetry.
web.telemetry-path             | Path under which to expose metrics.
exporter.aurora-url            | [URL](#aurora-url)

#### Aurora URL
Can be either a single ``http://host:port`` URL or a comma-separated ``zk://host1:port,zk://host2:port`` URL.

---
