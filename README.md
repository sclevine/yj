# yj
CLI - YAML &lt;-> JSON

```
opal:yj stephen$ .yj -h
Usage: yj [-][rcjenkh]

Converts stdin from JSON/YAML to YAML/JSON.

-r     Convert JSON to YAML instead of YAML to JSON
-c     Use CandiedYAML parser instead of GoYAML parser
-n     Do not covert Infinity, -Infinity, and NaN to/from strings
-h     Show this help message

YAML to JSON options:

-e     Escape HTML in JSON output (ignored for JSON to YAML)

JSON to YAML (-r) options:

-y     Use a YAML decoder instead of a JSON decoder to parse JSON
-k     Attempt to parse keys as JSON objects/numbers
```
