# store-exporter
_Utility to extract metrics from arbitary data stores in Prometheus format_

## Overview

Export your custom app metrics from external data _stores_ like PostgreSQL, MySQL, Redis(coming soon!)

## Features

- Extract column names from results and expose them as custom metric labels.
- Ability to register multiple jobs with different stores.

## Table of Contents

- [Getting Started](#getting-started)
  - [Motivation](#motivation)
  - [Installation](#installation)
  - [Quickstart](#quickstart)
  - [Sending a sample scrape request](#testing-a-sample-alert)

- [Advanced Section](#advanced-section)
  - [Configuration options](#configuation-options)
  - [Setting up Prometheus](#setting-up-prometheus)


### Motivation

`store-exporter` loads SQL query file and fetches the data from DB and transforms the result in Prometheus text format. A lot of times, it is undesirable to add instrumentation right in your app for the following reasons:

- Your app doesn't have any HTTP server, but to just extract metrics you've to invoke HTTP server.
- Your app cares about being _fast_ in which case adding any external library penalises performance.
- You don't want to mix the app logic with the metric collection/exposition logic.

In all the above cases, it is more suitable to take a [Sidecar approach](https://docs.microsoft.com/en-us/azure/architecture/patterns/sidecar), where you query for metrics from an external persistent store your app maintains. This utility just makes it easier for anyone to write custom SQL queries and expose metrics without having to worry about Prometheus format/exposition logic. You can run a single binary anywhere in your cluster environment which has access to the external store which exposes the metrics on an HTTP server confirming to Prometheus metric format.


### Installation

There are multiple ways of installing `store-exporter`.

### Running as docker container

[mrkaran/store-exporter](https://hub.docker.com/r/mrkaran/store-exporter)

`docker run -p 9609:9609 -v /etc/store-exporter/config.toml:/etc/store-exporter/config.toml mrkaran/store-exporter:latest`

### Precompiled binaries

Precompiled binaries for released versions are available in the [_Releases_ section](https://github.com/mr-karan/store-exporter/releases/).

### Compiling the binary

You can checkout the source code and build manually:

```bash
git clone https://github.com/mr-karan/store-exporter.git
cd store-exporter
make build
cp config.sample.toml config.toml
./store-exporter
```

### Quickstart

```sh
mkdir store-exporter && cd store-exporter/ # copy the binary and config.sample in this folder
cp config.toml.sample config.toml # change the settings like server address, job metadata, aws credentials etc.
./store-exporter # this command starts a web server and is ready to collect metrics from EC2.
```

### Testing a sample scrape request

You can send a `GET` request to `/metrics` and see the following metrics in Prometheus format:

```bash
# HELP job_name_basicname this is such a great help text
# TYPE job_name_basicname gauge
job_name_basicname{job="myjob",pg_db_blks_hit="74400",pg_db_tup_inserted="120"} 13713
# HELP job_name_verybasic_name this is such a great help text again
# TYPE job_name_verybasic_name gauge
job_name_verybasic_name{job="myjob",pg_db_conflicts="0",pg_db_temp_bytes="0"} 40
# HELP version Version of store-exporter
# TYPE version gauge
version{build="846771f (2019-08-28 10:28:07 +0530)"} 1
```

## Advanced Section

### Configuration Options

- **[server]**
  - **address**: Port which the server listens to. Default is *9608*
  - **name**: _Optional_, human identifier for the server.
  - **read_timeout**: Duration (in milliseconds) for the request body to be fully read) Read this [blog](https://blog.cloudflare.com/the-complete-guide-to-golang-net-http-timeouts/) for more info.
  - **write_timeout**: Duration (in milliseconds) for the response body to be written.

- **[app]**
  - **log_level**: "production" for all `INFO` level logs. If you want to enable verbose logging use "debug".
  - **jobs**
    - **name**: Unique identifier for the job.
    - **query**: Path to SQL file.
    - **db**: Type of SQL DB. Supported values: [postgres, mysql].
    - **dsn**: Connection URL to the DB.
    - **metrics**:
      - **name**: Name of the metric.
      - **help**: Helptext for the metric
      - **query**: Name of the query mapped in `sql` file, used to query the db for this metric.
      - **labels**: List of column names fetched from the DB, to be used in metric as key/value pairs.
      - **value**: Column name, for which the value should be used for the metric.

**NOTE**: You can use `--config` flag to supply a custom config file path while running `store-exporter`.

### Setting up Prometheus

You can add the following config under `scrape_configs` in Prometheus' configuration.

```yaml
  - job_name: 'store-exporter'
    metrics_path: '/metrics'
    static_configs:
    - targets: ['localhost:9610']
      labels:
        service: my-app-metrics
```

Validate your setup by querying `version` to check if store-exporter is discovered by Prometheus:

```plain
`version{build="846771f (2019-08-28 10:28:07 +0530)"} 1`
```

## Contribution

PRs on Feature Requests, Bug fixes are welcome. Feel free to open an issue and have a discussion first. Contributions on more external stores are also welcome and encouraged.

Read [CONTRIBUTING.md](CONTRIBUTING.md) for more details.

## License

[MIT](license)