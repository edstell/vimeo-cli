#!/usr/bin/env bash

# Check enough arguments passed.
if [ "$#" -ne 1 ]; then
	echo "Usage: $0 JSON_FILE" >&2; exit 1
fi

# Check argument is a file.
if ! [ -f $1 ]; then
	echo "JSON_FILE: Must be a file" >&2; exit 1
fi

# Check second argument is valid JSON.
if ! python -m json.tool "$1" > /dev/null; then
	echo "$1 must contain valid JSON" >&2; exit 1
fi

# Pipe file content to json script to escape content.
printf '%s' `cat $1` | python -c 'import json,sys; print(json.dumps(sys.stdin.read()))'
