# yj

[![Build Status](https://travis-ci.org/sclevine/yj.svg?branch=master)](https://travis-ci.org/sclevine/yj)
[![GoDoc](https://godoc.org/github.com/sclevine/yj?status.svg)](https://godoc.org/github.com/sclevine/yj)

Convert between YAML, TOML, JSON, and HCL.

Preserves map order.

Format versions:
- YAML: **v1.2** (gopkg.in/yaml.v3)
- TOML: **v1.0.0-rc.1** (github.com/pelletier/go-toml)
- JSON: **RFC 7159** (encoding/json)
- HCL: **v1** (github.com/hashicorp/hcl)

```
opal:yj stephen$ yj -h
Usage: yj [-][ytjcrneikhv]

Convert between YAML, TOML, JSON, and HCL.
Preserves map order.

-x[x]  Convert using stdin. Valid options:
          -yj, -y = YAML to JSON (default)
          -yy     = YAML to YAML
          -yt     = YAML to TOML
          -yc     = YAML to HCL
          -tj, -t = TOML to JSON
          -ty     = TOML to YAML
          -tt     = TOML to TOML
          -tc     = TOML to HCL
          -jj     = JSON to JSON
          -jy, -r = JSON to YAML
          -jt     = JSON to TOML
          -jc     = JSON to HCL
          -cy     = HCL to YAML
          -ct     = HCL to TOML
          -cj, -c = HCL to JSON
          -cc     = HCL to HCL
-n     Do not covert inf, -inf, and NaN to/from strings (YAML or TOML only)
-e     Escape HTML (JSON out only)
-i     Indent output (JSON or TOML out only)
-k     Attempt to parse keys as objects or numbers types (YAML out only)
-h     Show this help message
-v     Show version
```

### Installation

`yj` is available for macOS via [Homebrew](https://brew.sh/):
```bash
$ brew install yj
```

Binaries for macOS, Linux, and Windows are [attached to each release](https://github.com/sclevine/yj/releases).

`yj` is also available as a [Docker image](https://hub.docker.com/r/sclevine/yj).

### Go Package

Packages contained in this repo may also be used to convert all supported data formats to a normalized tree of ordered Go objects.
See [godoc](https://godoc.org/github.com/sclevine/yj) for details.
