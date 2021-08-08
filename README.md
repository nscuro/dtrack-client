# dtrack-client

[![CI](https://github.com/nscuro/dtrack-client/actions/workflows/ci.yml/badge.svg)](https://github.com/nscuro/dtrack-client/actions/workflows/ci.yml)
[![PkgGoDev](https://pkg.go.dev/badge/github.com/nscuro/dtrack-client)](https://pkg.go.dev/github.com/nscuro/dtrack-client)
[![License](https://img.shields.io/badge/license-Apache%202.0-brightgreen.svg)](LICENSE)

*Go client library for [OWASP Dependency-Track](https://dependencytrack.org/)*

## Installation

```
GO111MODULE=on go get github.com/nscuro/dtrack-client
```

## Compatibility

|  Go   | Dependency-Track |
| :---: | :--------------: |
| 1.14+ |      4.0.0+      |

## Usage

See [`examples`](./examples).

## API Coverage

*dtrack-client* primarily covers those parts of the Dependency-Track API that I personally need.
If you'd like to use this library, and your desired functionality is not yet available, please consider creating a PR.
