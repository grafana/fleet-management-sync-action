# fleet-management-sync-action

A GitHub Action to sync pipeline configurations to [Grafana Fleet Management](https://grafana.com/docs/grafana-cloud/send-data/fleet-management/introduction/).

## Overview

`fleet-management-sync-action` discovers pipeline configuration files in your repository and syncs them to Grafana Fleet Management. It recursively searches for YAML and/or .alloy files containing pipeline definitions and uploads them using the Fleet Management API.

## Usage

```yaml
name: Sync Pipelines to Fleet Management
on:
  push:
    branches: [main]
    paths:
      - 'pipelines/**'
      - '.github/workflows/fleet-management-sync-action.yml'

jobs:
  sync:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      
      - uses: grafana/fleet-management-sync-action@v1  # Replace with actual version
        with:
          pipelines-root-path: './pipelines'
          fm-username: ${{ secrets.FM_USERNAME }}
          fm-token: ${{ secrets.FM_TOKEN }}
```

## Inputs

| Input | Description | Required | Default |
|-------|-------------|----------|---------|
| `pipelines-root-path` | Root path to start searching for pipeline YAML files | Yes | - |
| `fm-username` | Fleet Management username for authentication | Yes | - |
| `fm-token` | Fleet Management API token for authentication | Yes | - |

## Pipeline Configuration

Pipeline files can be defined in two ways:

### Option 1: YAML Configuration File

YAML files (`.yaml` or `.yml`) with the following structure:

```yaml
name: my_pipeline    # Optional - defaults to filename without extension
contents: |          # Inline pipeline configuration
  logging {
    level = "info"
  }
enabled: true
matchers:
  - "environment=production"
  - "service=api"
```

### Option 2: Alloy Configuration File

Alloy files (`.alloy`) with a metadata header:

```alloy
/* fleet-management
name: my_pipeline
enabled: true
matchers:
  - "environment=production"
  - "service=api"
*/

logging {
  level = "info"
}
```

The metadata block must be at the top of the file, but can be preceded by whitespace.

## Examples

### Example 1: YAML-based Pipelines

This example shows how to define pipelines using YAML files.

#### Directory Structure

```
.
└── pipelines/
    ├── frontend.yaml
    └── another-pipeline.yaml
```

#### frontend.yaml

```yaml
# No `name` field means that this pipeline will be synced with the filename
# without the extension (frontend)
contents: |
  prometheus.scrape "app_metrics" {
    targets = [{__address__ = "localhost:8080"}]
  }
enabled: true
matchers:
  - "app=frontend"
```

### Example 2: Alloy-based Pipelines

This example shows how to define pipelines using Alloy files.

#### Directory Structure

```
.
└── pipelines/
    ├── backend.alloy
    └── another-pipeline.alloy
```

#### backend.alloy

```alloy
/* fleet-management
name: backend-pipeline
enabled: true
matchers:
  - "app=backend"
*/

logging {
  level = "info"
}
```
