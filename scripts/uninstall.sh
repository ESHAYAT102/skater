#!/usr/bin/env sh
set -eu

binary_name="skater"
skate_binary_name="skate"
install_dir="${XDG_BIN_HOME:-${HOME}/.local/bin}"
remove_skate=0

usage() {
	echo "usage: uninstall.sh [--with-skate]"
	echo
	echo "Removes $install_dir/$binary_name."
	echo "Pass --with-skate to also remove $install_dir/$skate_binary_name."
}

while [ "$#" -gt 0 ]; do
	case "$1" in
		--with-skate)
			remove_skate=1
			;;
		-h | --help)
			usage
			exit 0
			;;
		*)
			echo "error: unknown option: $1" >&2
			usage >&2
			exit 1
			;;
	esac
	shift
done

remove_binary() {
	path="$1"
	name="$2"

	if [ -e "$path" ]; then
		rm -f "$path"
		echo "removed $name from $path"
	else
		echo "$name was not found at $path"
	fi
}

remove_binary "$install_dir/$binary_name" "$binary_name"

if [ "$remove_skate" -eq 1 ]; then
	remove_binary "$install_dir/$skate_binary_name" "$skate_binary_name"
else
	echo "left $skate_binary_name installed; rerun with --with-skate to remove $install_dir/$skate_binary_name"
fi
