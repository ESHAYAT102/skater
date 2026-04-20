
# skater

A minimal Bubble Tea TUI for [skate](https://github.com/charmbracelet/skate).

## Install

The installer clones this repo, builds `skater` with Go, copies the binary into your user binary directory, and removes the cloned repo folder when it is done.

If `go` or `git` is missing, the installer will ask before installing them with the platform package manager:

- Arch: `pacman`
- Ubuntu/Debian: `apt`
- Fedora/RHEL: `dnf` or `yum`
- openSUSE: `zypper`
- Alpine: `apk`
- Void: `xbps-install`
- macOS: `brew`
- Windows: `winget`, Chocolatey, or Scoop

It also asks before installing `skate` if it is missing.

Manual install with Go:

```sh
go install github.com/charmbracelet/skate@latest #skate is required
git clone https://github.com/ESHAYAT102/skater.git
cd skater
go build -o "$HOME/.local/bin/skater" .
```

MacOS and Linux:

```sh
curl -fsSL https://raw.githubusercontent.com/ESHAYAT102/skater/main/scripts/install.sh | sh
```

Windows PowerShell:

```powershell
irm https://raw.githubusercontent.com/ESHAYAT102/skater/main/scripts/install.ps1 | iex
```

On macOS and Linux, the binary is installed to:

```sh
~/.local/bin/skater
~/.local/bin/skate
```

On Windows, the binary is installed to:

```powershell
$HOME\.local\bin\skater.exe
$HOME\.local\bin\skate.exe
```

## Run

```sh
skater
```

## Controls

- `tab`: move focus
- `enter`: save when an input is focused, edit selected row when the table is focused
- `d`: delete selected row
- `r`: refresh
- `q` or `ctrl+c`: quit
