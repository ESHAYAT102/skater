#!/usr/bin/env sh
set -eu

binary_name="skater"
install_dir="${HOME}/.local/bin"

need() {
	if ! command -v "$1" >/dev/null 2>&1; then
		echo "error: $1 is required" >&2
		exit 1
	fi
}

need go

mkdir -p "$install_dir"

if ! command -v skate >/dev/null 2>&1 && [ ! -x "$install_dir/skate" ]; then
	echo "skate not found; installing skate"
	GOBIN="$install_dir" go install github.com/charmbracelet/skate@latest
	chmod +x "$install_dir/skate"
fi

GOBIN="$install_dir" go install github.com/ESHAYAT102/skater@latest
chmod +x "$install_dir/$binary_name"

echo "installed $binary_name to $install_dir/$binary_name"

case ":$PATH:" in
	*":$install_dir:"*) ;;
	*)
		echo "warning: $install_dir is not in PATH"
		echo "add this to your shell profile:"
		echo "  export PATH=\"\$HOME/.local/bin:\$PATH\""
		;;
esac
