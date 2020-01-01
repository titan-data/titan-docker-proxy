# Project Development

For general information about contributing changes, see the
[Contributor Guidelines](https://github.com/titan-data/.github/blob/master/CONTRIBUTING.md).

## How it Works

The 

## Building

To build the project, run `go build ./...`. This is equivalent to building `cmd/docker-volume-proxy/main.go`. This
will create a binary named `docker-volume-proxy` in the root of the directory.

## Testing

To test the project, run `go test ./...`. This will run all tests.

## Releasing

To release, create a tag and push it. This will build the resulting go binary for Linux (the runtime for the
titan-server container) and upload it as an artifact to the draft release. Release notes are maintained on each
push through the release drafter action.
