# yj
CLI - YAML &lt;-> TOML &lt;-> JSON &lt;-> HCL

```
opal:yj stephen$ yj -h
Usage: yj [-][ytjcrnekh]

Convert YAML, TOML, JSON, or HCL to YAML, TOML, JSON, or HCL.

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
-n     Do not covert Infinity, -Infinity, and NaN to/from strings
-e     Escape HTML (JSON output only)
-k     Attempt to parse keys as objects or numbers types (YAML output only)
-h     Show this help message
```
