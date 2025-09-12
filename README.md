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
          namespace: ${{ github.event.repository.name }}
```

## Inputs

| Input | Description | Required | Default |
|-------|-------------|----------|---------|
| `pipelines-root-path` | Root path to start searching for pipeline YAML files | Yes | - |
| `fm-username` | Fleet Management username for authentication | Yes | - |
| `fm-token` | Fleet Management API token for authentication | Yes | - |
| `namespace` | Namespace for the pipelines. See "Configuring the Namespace" for examples. | Yes | - |

## Configuring the Namespace

The `namespace` input is used to scope sync cleanup, so a sync with a given namespace will only clean up pipelines no longer present in the repository, and not any pipelines created from other sources (i.e., with a different namespace).

Here are some examples of how you can configure the namespace, though you can use any value that you want, within length limits:

### Hardcoded Namespace

You can use a simple hardcoded string for the namespace.

```yaml
- uses: grafana/fleet-management-sync-action@v1
  with:
    namespace: "my-production-pipelines"
    # ... other inputs
```

### Namespace from Repository Name

You can use GitHub Actions expressions to dynamically set the `namespace` to the name of the repository.

```yaml
- uses: grafana/fleet-management-sync-action@v1
  with:
    namespace: ${{ github.event.repository.name }}
    # ... other inputs
```

## Pipeline Configuration

Pipeline definitions require two files in the same directory:

1. A YAML file (`.yaml` or `.yml`) for metadata.
2. An Alloy file (`.alloy`) for the pipeline's contents.

The YAML and Alloy files must share the same base name (e.g., `my-pipeline.yaml` and `my-pipeline.alloy`).

You can use any directory structure you want, the action will recursively search the `pipelines-root-path` for YAML and Alloy files.

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

### Directory Structure

```text
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
