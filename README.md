# fleet-management-sync-action

A GitHub Action to sync pipeline configurations to [Grafana Fleet Management](https://grafana.com/docs/grafana-cloud/send-data/fleet-management/introduction/).

## Overview

`fleet-management-sync-action` discovers pipeline configuration files in your repository and syncs them to Grafana Fleet Management. It recursively searches for YAML files containing pipeline definitions and uploads them using the Fleet Management API.

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

Pipeline definitions require two files in the same directory:

1.  A YAML file (`.yaml` or `.yml`) for metadata.
2.  An Alloy file (`.alloy`) for the pipeline's contents.

The YAML and Alloy files must share the same base name (e.g., `my-pipeline.yaml` and `my-pipeline.alloy`).

### Example

**`my-pipeline.yaml` (Metadata)**
```yaml
name: my-pipeline    # Optional - defaults to filename without extension
enabled: true
matchers:
  - environment=production
  - service=api
```

**`my-pipeline.alloy` (Contents)**
```alloy
logging {
  level = "info"
}
```

## Examples

### Directory Structure

```
.
└── pipelines/
    ├── frontend.yaml
    ├── frontend.alloy
    ├── backend.yaml
    ├── backend.alloy
    └── monitoring/
        ├── o11y.yaml
        └── o11y.alloy
```
