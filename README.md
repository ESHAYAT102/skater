
# skater

A minimal Bubble Tea TUI for [skate](https://github.com/charmbracelet/skate).

## Install

Requirements:

- `go`
- `skate`

Install with Go:

```sh
go install github.com/charmbracelet/skate@latest #skate is required
go install github.com/ESHAYAT102/skater@latest
```

MacOS and Linux:

```sh
curl -fsSL https://raw.githubusercontent.com/ESHAYAT102/skater/main/scripts/install.sh | sh
```

Windows PowerShell:

```powershell
irm https://raw.githubusercontent.com/ESHAYAT102/skater/main/scripts/install.ps1 | iex
```

The installer installs `skate` first if it is missing, then installs `skater` with `go install`.

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

