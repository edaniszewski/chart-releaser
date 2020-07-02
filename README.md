# chart-releaser

A CI tool to automate Helm Chart updates for new application releases.

This tool was developed to help make maintaining Helm Chart repos more sane. It works
for a single Chart, or for repos holding multiple charts. By integrating the tool into
CI pipelines, a new application release can trigger an automatic update to its corresponding
Chart to update the appVersion and Chart version.

> **Note**: This tool is still in its infancy so there may be bugs, missing features,
> non-ergonomic defaults, etc. If you have any questions, notice any issues, or have
> suggestions for improvement, [open an issue](https://github.com/edaniszewski/chart-releaser/issues)
> or open a PR - all contributions are welcome.

## Getting Started

### Getting

#### Latest

The latest release is available from the [releases](https://github.com/edaniszewski/chart-releaser/releases)
page. All tags generate a release with pre-compiled binaries attached as assets.

#### Docker

A lightweight Docker image is also available:

```
docker pull chartreleaser/chart-releaser
``` 

#### Homebrew

macOS users can also install via [Homebrew](https://brew.sh/) by first adding the tap,
then installing the package.

```
brew tap edaniszewski/tap
brew install chart-releaser
```

### Running

`chart-releaser` can be run from the command line with no argument or with the `--help` flag
to print usage info. The main entrypoint for kicking of a Chart update is through the `chart-releaese update`
command.

It looks for a `.chartreleaser.yml` configuration file in the directory which it is run out of
and executes an update based on those config options.

```
chart-releaser update
```


## Configuring

`chart-releaser` is configured on a per-project basis, which usually translates to one config file
for a repository. The config file should live in the repository for the application, not for the
repository for its Helm Chart (assuming they are different).

The configuration is defined in a YAML-formatted file named`.chartreleaser.yml`.

Below is a description and examples of the different config options that are currently supported.

### v1

These configuration options apply to config files following the v1 schema, as denoted by the
`version: v1` field in the configuration.

Note that in the tables below, a default value of `-` indicates no default and that the configuration
is required.

#### Top Level Keys

Top-level keys for a `chart-releaser` configuration.

| Key | Description | Value |
| --- | ----------- | ----- |
| `version` | The version of the configuration scheme. This must be `v1` for v1 configs. | `-` |
| `chart`   | Defines where the Helm Chart exists for the configured project. | [Chart](#chart) |
| `publish` | Defines how Chart/file updates should be published to the chart repo. | [Publish](#publish) |
| `commit`  | Defines the author of the commit(s) made to the chart repo. | [Commit](#commit) |
| `release` | Defines behavior for how application releases are targeted and they affect the Chart version. | [Release](#release) |
| `extras`  | Defines any non-Chart.yaml files that should also be updated. | [Extras](#extras) |

#### Chart

> Defines where the Helm Chart exists for the configured project.

| Key | Description | Default |
| --- | ----------- | ------- |
| `name` | The name of the Chart. | `-` |
| `repo` | The name of the repository holding the Chart for the project. This is required. It should follow the format `{{ RepoType }}/{{ Owner }}/{{ Name }}`. The currently supported repo types are: `github.com`. | `-` |
| `path` | The sub-path to the chart in the repository. If this is empty, it assumes the Chart.yaml is in the root of the specified repository. | `""` |

#### Publish

> Defines how Chart/file updates should be published to the chart repo.

This defines a `commit` and `pr` key. If neither are specified, it defaults to "pr". Only one of
the two keys may be specified at once.

`commit` commits directly to the configured branch, but will not open a PR. `pr` will commit to the
configured branch and open a PR for the changes.

| Key | Description | Default |
| --- | ----------- | ------- |
| `commit.branch` | The name of the branch to commit to. | `master` |
| `commit.base` | The base ref to create the new branch from. | `master` |
| `pr.branch_template` | The name of the branch to commit to. This may be a template populated by update context. | `chartreleaser/{{ .Chart.Name }}/{{ .Chart.NewVersion }}` |
| `pr.base` | The base ref to create the new branch from. | `master` |
| `pr.title_template` | The title to use for the pull request. | see: [templates.go](./pkg/templates/templates.go) |
| `pr.body_template` | The pull request body comment. | see: [templates.go](./pkg/templates/templates.go) |

#### Commit

> Defines the author of the commit(s) made to the chart repo.

| Key | Description | Default |
| --- | ----------- | ------- |
| `author.name` | Name of the commit author for published updates. | The git user name from gitconfig |
| `author.email` | Email of the commit author for published updates. | The git user email from gitconfig |
| `templates.update` | The template for the commit message to use on update. | see: [templates.go](./pkg/templates/templates.go) |

#### Release

> Defines behavior for how application releases are targeted and they affect the Chart version.

| Key | Description | Default |
| --- | ----------- | ------- |
| `matches` | A list of regex-compilable strings defining constraints that an application tag needs to meet to be eligible for chart-releaser update. | `[]` |
| `ignores` | A list of regex-compilable strings defining constraints which prevent an application tag from being eligible for chart-releaser update. | `[]` |
| `strategy` | The release strategy to use. See below for supported strategies. | `default` |

##### Strategies

The supported release strategies are:

* `default`: Bump the patch version of the Chart for any non-prerelease update to the app version. If the app has a prerelease bump, the chart version either gets a new prerelease suffix, or gets an existing prerelease suffix bumped.
* `major`: Bump the major version of the Chart for any update to the app version.
* `minor`: Bump the minor version of the Chart for any update to the app version.
* `patch`: Bump the patch version of the Chart for any update to the app version.

#### Extras

> Defines any non-Chart.yaml files that should also be updated.

Extras as specified as a list, e.g.

```yaml
extras:
  - path: some/path/to/README.md
    updates:
      - search: foo
        replace: bar
        limit: 1
```

| Key | Description | Default |
| --- | ----------- | ------- |
| `path` | The path to the file to update, starting from the root of the charts repository. | `""` |
| `updates[*].search` | A regex-compilable string defining the pattern that should be matched up for replacement within the file. | `-` |
| `updates[*].replace` | The value to replace the found matches. This may be a template, which gets rendered with the update context. | `-` |
| `updates[*].limit` | Limit the number of replaces to be performed in the file. A value of `0` indicates that there is no limit and all found matches should be replaced. | `0` |

#### Context

The `v1` config schema update context is used to render various templates in the configuration.
See [context.go](./pkg/v1/ctx/context.go) for details on the fields provided by the context.

## License

`chart-releaser` is released under the [MIT License](LICENSE).
