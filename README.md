# Unreal Version Updater

This Go Binary looks through a directory at all `*.ini` files for the `ProjectVersion` Key under the header `/Script/EngineSettings.GeneralProjectSettings` and updates it to whatever version you specify as the arg.

# Getting Started
## Github Action
```yaml
      - uses: Benbentwo/UnrealGameVersionUpdater@master
        with:
          version: ${{ steps.create_release.outputs.tag_name }}
```

### Sample `auto-releaser.yaml` for Unreal Engine Projects
```yaml
name: auto-release

on:
  push:
    branches:
      - main
      - production

jobs:
  publish:
    runs-on: ubuntu-latest
    steps:
      # Get PR from merged commit to master
      - uses: actions-ecosystem/action-get-merged-pull-request@v1
        id: get-merged-pull-request
        with:
          github_token: ${{ secrets.GITHUB_ACTION }}

      - name: Checkout
        uses: actions/checkout@v3

      # prepares and sets versions for future steps
      - uses: release-drafter/release-drafter@v5
        id: create_release
        with:
          publish: false
          prerelease: true
          config-name: auto-release.yaml
        env:
          GITHUB_TOKEN: ${{ secrets.PUBLIC_REPO_ACCESS_TOKEN }}

      - uses: Benbentwo/UnrealGameVersionUpdater@master
        with:
          version: ${{ steps.create_release.outputs.tag_name }}

      - uses: EndBug/add-and-commit@v9 # You can change this to use a specific version.
        with:
          message: 'Updating Unreal Engine Version to `${{ steps.create_release.outputs.tag_name }}`'

      - uses: eregon/publish-release@v1
        env:
          GITHUB_TOKEN: ${{ secrets.PUBLIC_REPO_ACCESS_TOKEN }}
        with:
          release_id: ${{ steps.create_release.outputs.id }}
```


## Run Locally with Docker:
Example running locally to set the version to 1.0.0
```shell
# Linux
docker run --rm -it -v $(pwd):/app ghcr.io/benbentwo/unrealgameversionupdater:latest 1.0.0

# Windows (PowerShell)
docker run --rm -it -v ${PWD}:/app ghcr.io/benbentwo/unrealgameversionupdater:latest 1.0.0

# Windows (CMD)
docker run --rm -it -v %cd%:/app ghcr.io/benbentwo/unrealgameversionupdater:latest 1.0.0
```

# CLI
this runs as a CLI

## Args:
| Arg | Shorthand | Description | Default |
| --- | --- | --- | --- |
| `--project` | `-p` | Is a Project (unused right now) | `true`
| `--plugin` | `-l` | Is a Plugin (unused right now) | `false`
| `--config` | `-c` | Folder to search for INI Files. This can be changed if your version lives in a nested folder. | `Config`
| `--verbose` | `-v` | Verbose Logging (sets log level to debug) | null

