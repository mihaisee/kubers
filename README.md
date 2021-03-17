# kubers

[![Build Go](https://github.com/mihaisee/kubers/actions/workflows/build.yml/badge.svg)](https://github.com/mihaisee/kubers/actions/workflows/build.yml)

CLI for providing an easier way to inspect cluster resource usage.

**FYI:** It uses your active `kubectl` context for connecting to cluster.


**SECOND FYI:** It expects you have at least `go 1.16` installed.
```shell
brew install go
```

## Installation

`make install`

## Usage
```shell
kubers -h
```

## Example usage
```shell
kubers get po # Get a list of pods with resources (usage/request/limit)
kubers get ns # Get a list of namespaces with resources (usage/request/limit)
```

### Filter by ns
```shell
kubers get po -n staging
```

### Filter by labels
```shell
kubers get ns -l selector=team-ns
```

### Order
```shell
kubers get po -o asc -b cpu
kubers get ns -o desc -b mem
```

### Get pods with containers detailed
```shell
kubers get po -c
```
