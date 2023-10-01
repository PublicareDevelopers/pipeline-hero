[![Maintenance](https://img.shields.io/badge/Maintained%3F-yes-green.svg)](https://github.com/PublicareDevelopers/pipeline-hero/graphs/commit-activity)
![Maintainer](https://img.shields.io/badge/maintainer-DionTech-blue)
[![made-with-Go](https://img.shields.io/badge/Made%20with-Go-1f425f.svg)](https://go.dev/)
[![GitHub license](https://img.shields.io/github/license/PublicareDevelopers/pipeline-hero.svg)](https://github.com/PublicareDevelopers/pipeline-hero/blob/main/LICENSE)
[![GitHub branches](https://badgen.net/github/branches/PublicareDevelopers/pipeline-hero)](https://github.com/PublicareDevelopers/pipeline-hero/)
[![Github tag](https://badgen.net/github/tag/PublicareDevelopers/pipeline-hero)](https://github.com/PublicareDevelopers/pipeline-hero/tags/)
[![GitHub latest commit](https://badgen.net/github/last-commit/PublicareDevelopers/pipeline-hero)](https://GitHub.com/PublicareDevelopers/pipeline-hero/commit/)
[![GitHub pull-requests](https://img.shields.io/github/issues-pr/PublicareDevelopers/pipeline-hero.svg)](https://GitHub.com/PublicareDevelopers/pipeline-hero/pull/)
[![Github all releases](https://img.shields.io/github/downloads/PublicareDevelopers/pipeline-hero/total.svg)](https://GitHub.com/PublicareDevelopers/pipeline-hero/releases/)

# pipeline-hero

Warning: This is a work in progress. The API is not stable and may change at any time. Use at your own risk.

## install
```shell
go install github.com/PublicareDevelopers/pipeline-hero@latest
```

## get help
```shell
pipeline-hero --help
```

## use in a pipeline
```shell
pipeline-hero pipe analyse
```

```
Usage:
  pipeline-hero pipe analyse [flags]

Flags:
  -c, --coverage-threshold float   Coverage threshold to use (default 75)
  -e, --env stringToString         Environment variables to set (default [])
  -h, --help                       help for analyse
  -s, --slack                      Send results to slack
  -t, --test-setup string          Test setup to use (default "./...")
```
