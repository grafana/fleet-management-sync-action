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
          fm-url: ${{ secrets.FM_URL }}
          fm-username: ${{ secrets.FM_USERNAME }}
          fm-token: ${{ secrets.FM_TOKEN }}
          namespace: ${{ github.event.repository.name }}
```

## Inputs

| Input | Description | Required | Default |
|-------|-------------|----------|---------|
| `pipelines-root-path` | Root path to start searching for pipeline YAML files | Yes | - |
| `fm-url` | [Fleet Management API URL](https://grafana.com/docs/grafana-cloud/send-data/fleet-management/api-reference/pipeline-api/#find-the-base-url) | Yes | - |
| `fm-username` | Fleet Management username for authentication | Yes | - |
| `fm-token` | Fleet Management API token for authentication | Yes | - |
| `namespace` | Namespace for the pipelines. See "Configuring the Namespace" for examples. | Yes | - |
| `global-matcher` | Global matcher added to all pipelines. Useful for targeting specific collector groups across multiple repositories. See "Using Global Matchers" for examples. | No | - |

## Configuring the Namespace

The `namespace` input is used to scope sync cleanup. During a sync, the action will delete pipelines in Fleet Management that match the same source type (Git) and namespace but are NOT included in the current sync request. This ensures your API state matches your repository while leaving pipelines from other sources or namespaces untouched.

**Note:** Namespace and source type are used to determine which pipelines to delete during cleanup. However, pipelines are identified by name during the sync itself. If a pipeline with the same name already exists in Fleet Management (even from a different source or namespace), it will be overwritten. Ensure pipeline names are unique across all sources to avoid conflicts.

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

## Using Global Matchers

The `global-matcher` input allows you to add a matcher to all pipelines synced from a repository. This is useful if you only want a particular group of collectors to receive the pipelines from a particular repo.

For example, if you manage separate repositories for different environments or teams:

```yaml
- uses: grafana/fleet-management-sync-action@v1
  with:
    global-matcher: "namespace=prod"
    # ... other inputs
```

or

```yaml
- uses: grafana/fleet-management-sync-action@v1
  with:
    global-matcher: "team=platform"
    # ... other inputs
```

## Pipeline Configuration

Pipeline definitions require two files in the same directory:

1. A YAML file (`.yaml` or `.yml`) for metadata.
2. An Alloy file (`.alloy`) for the pipeline's contents.

The YAML and Alloy files must share the same base name (e.g., `my-pipeline.yaml` and `my-pipeline.alloy`). If a YAML file does not have a corresponding Alloy file with the same base name, the action will fail with an error.

You can use any directory structure you want, the action will recursively search the `pipelines-root-path` for YAML and Alloy files.

### Example

**`my-pipeline.yaml` (Metadata)**

```yaml
name: my-pipeline    # Optional - defaults to the base filename without extension
enabled: true
matchers:
  - environment=production
  - service=api
```

**`my-pipeline.alloy` (Contents)**

```alloy
local.file "example" {
  filename = "example.txt"
}
```

**`monitoring/o11y.yaml` (Metadata)**

```yaml
name: o11y
enabled: true
matchers:
  - environment=production
  - service=monitoring
  - team=platform
```

**`monitoring/o11y.alloy` (Contents)**

```alloy
local.file "o11y" {
  filename = "o11y.txt"
}
```

### Directory Structure

```text
.
└── pipelines/
    ├── my-pipeline.yaml
    ├── my-pipeline.alloy
    └── monitoring/
        ├── o11y.yaml
        └── o11y.alloy
```

## Troubleshooting

### Authentication Failures

If you encounter authentication errors, verify that:

- `fm-username` and `fm-token` are correctly set in your workflow
- The token has not expired
- The Fleet Management URL is correct for your environment

### Missing Alloy Files

If you see errors about missing Alloy files, ensure that:

- Every YAML file has a corresponding `.alloy` file with the same base name
- Both files are in the same directory
- The Alloy file has the correct extension (`.alloy`, not `.alloy.txt` or similar)

### Duplicate Pipeline Names

Pipeline names must be unique across all discovered files. If you have duplicate names (either explicit or defaulted from filenames), the action will fail. To resolve:

- Explicitly set different `name` values in the YAML metadata files, or
- Rename the files to have unique base names
