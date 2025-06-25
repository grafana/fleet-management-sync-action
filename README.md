# fm-sync

A GitHub Action to sync pipeline configurations to Grafana Fleet Management.

## Overview

`fm-sync` discovers pipeline configuration files in your repository and syncs them to Grafana Fleet Management. It recursively searches for YAML files containing pipeline definitions and uploads them using the Fleet Management API.

## Usage

```yaml
name: Sync Pipelines to Fleet Management
on:
  push:
    branches: [main]
    paths:
      - 'pipelines/**'
      - '.github/workflows/fm-sync.yml'

jobs:
  sync:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      
      - uses: grafana/fleet-management-sync-action@v1  # Replace with actual version/digest
        with:
          root-path: 'pipelines'
          fm-username: ${{ secrets.FM_USERNAME }}
          fm-token: ${{ secrets.FM_TOKEN }}
```

## Inputs

| Input | Description | Required | Default |
|-------|-------------|----------|---------|
| `root-path` | Root path to start searching for pipeline YAML files | No | `.` |
| `fm-username` | Fleet Management username for authentication | Yes | - |
| `fm-token` | Fleet Management API token for authentication | Yes | - |

## Pipeline Configuration

Pipeline files must be YAML files (`.yaml` or `.yml`) with the following structure:

### Option 1: Inline Configuration

```yaml
name: my-pipeline    # Optional - defaults to filename without extension
contents: |          # Inline pipeline configuration
  logging {
    level = "info"
  }
enabled: true
matchers:
  - environment=production
  - service=api
```

### Option 2: External Configuration File

```yaml
name: my-pipeline
contents_file: configs/logging.alloy  # Path to external config file
enabled: true
matchers:
  - environment=production
  - service=api
```

**Important**: You must specify either `contents` OR `contents_file`, but not both.

### Path Resolution for `contents_file`

The `contents_file` path is resolved in the following order:

1. Absolute path - used as-is
2. Relative to the YAML file's directory
3. Relative to the `root-path`

## Examples

### Directory Structure

```
.
├── configs/
│   ├── logging.alloy
│   ├── tracing.alloy
│   └── metrics.alloy
└── pipelines/
    ├── frontend.yaml      # uses inline contents
    ├── backend.yaml       # uses relative path
    └── monitoring/
        └── o11y.yaml      # uses path from the repo root
```

### frontend.yaml - Inline Configuration

```yaml
# No `name` field means that this pipeline will be synced with the filename
# without the extension (frontend)
contents: |
  prometheus.scrape "app_metrics" {
    targets = [{__address__ = "localhost:8080"}]
  }
enabled: true
matchers:
  - app=frontend
```

### backend.yaml - External Configuration (Relative)

```yaml
name: api-monitoring
contents_file: ../configs/logging.alloy
enabled: true
matchers:
  - app=backend
  - tier=critical
```

### o11y.yaml - External Configuration (From root-path)

```yaml
name: infra-o11y
contents_file: configs/tracing.alloy
enabled: true
matchers:
  - app=infra-o11y
  - tier=critical
```
