#!/usr/bin/env bash
export VIMEO_ACCESS_TOKEN=TODO

# Check whether the access token has been set.
if [ $VIMEO_ACCESS_TOKEN = "TODO" ]; then
	echo "Please set your access token" >&2; exit 1
fi

# Check enough arguments passed.
if [ "$#" -ne 2 ]; then
	echo "Usage: $0 VIDEO_ID CONFIG_FILE" >&2; exit 1
fi

# Check first argument is a number.
re='^[0-9]+$'
if ! [[ $1 =~ $re ]]; then
	echo "VIDEO_ID: Must be a number" >&2; exit 1
fi

# Check second argument is a file.
if ! [ -f $2 ]; then
	echo "CONFIG_FILE: Must be a file" >&2; exit 1
fi

# Check second argument is valid JSON.
if ! python -m json.tool "$2" > /dev/null; then
	echo "$2 must contain valid JSON" >&2; exit 1
fi

# Check 'vimeo' executable present.
if ! command -v ./vimeo &> /dev/null; then
	echo "'vimeo' executable missing" >&2; exit 1
fi

# Execute edit video command, passing in arguments provided.
printf '[%s,%s]' $1 `cat $2` | ./vimeo Videos Edit | python -m json.tool
