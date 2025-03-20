#!/usr/bin/env bash

set -e

dir="$(dirname "$(realpath "$0")")"
root="$(realpath "${dir}/..")"

echo "watching ${root} for changes..."
while true; do
	changed="$(inotifywait -e modify,create,delete,move --quiet --recursive --include "\.go$" --format "%w%f" "${root}")"
  changedRelativePath="$(realpath -m --relative-to="${root}" "${changed}")"

	echo "change in ./${changedRelativePath}, rebuilding..."
	devbox run render
  echo "rebuilt."
done
