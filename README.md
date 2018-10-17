# Prometheus Aurora exporter

This is an exporter for Prometheus to get instrumentation data for [Apache Aurora](http://aurora.apache.org/)

## Build and run

    go build .
    ./aurora_exporter <flags>

### Flags

Name                            | Description
--------------------------------|------------
web.listen-address              | Address to listen on for web interface and telemetry.
web.telemetry-path              | Path under which to expose metrics.
exporter.aurora-url             | [URL](#aurora-url) to an Aurora scheduler or ZooKeeper ensemble.
zk.path                         | The path for the aurora scheduler znode, this defaults to `/aurora/scheduler`
exporter.bypass-leader-redirect | Don't follow redirects to the leader instance. Ignored for quotas scraping because it results in an internal server error on non-leading masters.

#### Aurora URL
Can be either a single ``http://host:port`` or a comma-separated ``zk://host1:port,zk://host2:port`` URL.

## Console Dashboard

Copy the content of `consoles` to the consoles folder used by your Prometheus master. Your Aurora
dashboard will then be available at `http://your-prometheus:9000/consoles/aurora.html`.

---
