# yj

[![Build Status](https://travis-ci.org/sclevine/yj.svg?branch=master)](https://travis-ci.org/sclevine/yj)
[![GoDoc](https://godoc.org/github.com/sclevine/yj?status.svg)](https://godoc.org/github.com/sclevine/yj)

Convert between YAML, TOML, JSON, and HCL.

```
opal:yj stephen$ yj -h
Usage: yj [-][ytjcrneikh]

Convert between YAML, TOML, JSON, and HCL.

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
-n     Do not covert inf, -inf, and NaN to/from strings (YAML in/out only)
-e     Escape HTML (JSON out only)
-i     Indent output (JSON or TOML out only)
-k     Attempt to parse keys as objects or numbers types (YAML out only)
-h     Show this help message
```
