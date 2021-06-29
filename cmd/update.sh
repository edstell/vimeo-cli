#!/usr/bin/env bash
export VIMEO_ACCESS_TOKEN=TODO

csv2json=./csv2json.jq
vimeo=./vimeo
template=./template
ops_tmpl=./ops.tmpl

config_pipe=/tmp/config

# Check whether the access token has been set.
if [ $VIMEO_ACCESS_TOKEN = "TODO" ]; then
	echo "Please set your access token" >&2; exit 1
fi

# Check enough arguments passed.
if [ "$#" -ne 1 ]; then
	echo "Usage: $0 CONFIG_CSV" >&2; exit 1
fi

# Check config file exists.
if ! [ -f $1 ]; then
	echo "CONFIG_CSV: Doesn't exist" >&2; exit 1
fi

# Check our csv2json converter is present.
if ! [ -f $csv2json ]; then
	echo "$csv2json csv converter missing" >&2; exit 1
fi

# Check the operations template exists.
if ! [ -f $ops_tmpl ]; then
	echo "$ops_tmpl operations template missing" >&2; exit 1
fi

# Check 'vimeo' executable present.
if ! command -v $vimeo &> /dev/null; then
	echo "'$vimeo' executable missing" >&2; exit 1
fi

# Check 'template' executable present.
if ! command -v $template &> /dev/null; then
	echo "'$template' executable missing" >&2; exit 1
fi

# Create our formatting pipe.
if [[ ! -p $config_pipe ]]; then
	mkfifo $config_pipe
fi
trap "rm -rf $config_pipe" EXIT

# Convert our config csv to json for templating, then iterate over each file
# in the config.
jq -c '.[]' <(jq -R -s -f $csv2json $1) | while read i; do
	jq -c '.[]' <(cat $ops_tmpl | $template <(echo $i)) | while read i; do
		service=$(echo $i | jq -r '.service')
		operation=$(echo $i | jq -r '.operation')
		$(echo $i | jq -r '.arguments') | $vimeo $service $operation 
		if [ $? -ne 0 ]; then
			echo "failed to execute $service $operation"
		fi
	done
done
