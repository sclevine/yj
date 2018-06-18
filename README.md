# yj
CLI - YAML &lt;-> TOML &lt;-> JSON

```
opal:yj stephen$ yj -h
Usage: ./yj [-][ytjrnekh]

Convert YAML, TOML, or JSON to YAML, TOML, or JSON.

-x[x]  Convert using stdin. Valid options:
          -yj, -y = YAML to JSON (default)
          -yy     = YAML to YAML
          -yt     = YAML to TOML
          -tj, -t = TOML to JSON
          -ty     = TOML to YAML
          -tt     = TOML to TOML
          -jj     = JSON to JSON
          -jy, -r = JSON to YAML
          -jt     = JSON to TOML
-n     Do not covert Infinity, -Infinity, and NaN to/from strings
-e     Escape HTML (JSON output only)
-k     Attempt to parse keys as objects or numbers types (YAML output only)
-h     Show this help message
```
