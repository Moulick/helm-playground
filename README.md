# Helm Playground

## Forked from https://github.com/shipmight/helm-playground

This is a fork of the original [Helm Playground](https://github.com/shipmight/helm-playground), with updated dependencies.
Main changes:
- Updated dependencies to the latest versions.
- Auto publish website when updated.
- Remove the start button on website load and instead directly load the wasm module.

This repository contains the source code for https://helm-playground.moulick.xyz

[![Screenshot of Helm Playground – Click to open](screenshot.png)](https://helm-playground.moulick.xyz)

## How it works

A piece of Go code is compiled to a [Wasm](https://en.wikipedia.org/wiki/WebAssembly) module which can be ran in the browser. This code implements a simple function which takes two inputs:

- YAML template
- YAML values

Then it simply renders the given template with the given values using [Sprig](https://github.com/Masterminds/sprig), which is also what Helm uses.

The Wasm module is compiled in a GitHub action. You can find the workflow in [`.github/workflows/compile.yaml`](.github/workflows/compile.yaml). When a commit is pushed to `master`, the workflow is triggered, the code is compiled and uploaded to Github Pages. The `master` branch is hosted live via GitHub Pages at https://helm-playground.moulick.xyz.

## Development

### Pull the repository

```bash
git clone https://github.com/Moulick/helm-playground
```

### Run tests for the golang code

```bash
make test
```

### Build Wasm from golang

```bash
make build
```

### Test the built Wasm code in headless browser

```bash
make browser-test
```

### Locally develop in browser

You need a HTTP server to run the site locally, because fetch doesn't work under `file://` protocol.

```bash
make serve
```

### Updating dependencies

Dependency versions in `go.mod` should mirror a recent Helm release.

See the top comment in [`go.mod`](go.mod) for the current mirrored Helm version and upgrade instructions.

## License

Some files in this repository contain embedded license notes.

Other files in this repository are licensed under GNU AGPLv3 (see [LICENSE](./LICENSE)).
