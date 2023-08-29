# SDSSH Challenge 2023 - Garten Aufräumen - Testing Tool

A testing tool and simulator for the [SDSSH Challenge 2023](https://github.com/maxwellmatthis/sdssh-challenge-2023).

## Installation

1. Download the executable from the [releases](https://github.com/juho05/sdssh-gartentester/releases/latest) page.
2. Rename the file to `gartentester.exe` (Windows) or `gartentester` (macOS/Linux)
3. On macOS/Linux: make the file executable:
```
chmod +x gartentester
```

## Usage

*The following instructions are for macOS and Linux. You might have to tweak the syntax a bit to work on Windows (replace `./` with `.\`, append `.exe` to executables, …).*

### Generate a random garden

```sh
./gartentester -generate output.txt
```

#### Set size

```sh
# generate a garden with width 64 and height 32
./gartentester -generate -size 64x32 output.txt

# generate a garden with a random size
./gartentester -generate -size random output.txt
```

#### Set area count

```sh
# generate a garden with 8 areas
./gartentester -generate -area-count 8 output.txt

# generate a garden with a random area count
./gartentester -generate -area-count random output.txt
```

### Test a command sequence

#### Pipe the output of your program into the testing tool

```sh
./my-program | ./gartentester example-garden.txt
```

#### Read commands from a file

```sh
./gartentester -input commands.txt example-garden.txt
```

#### Disable animation (faster)

```sh
./gartentester -no-delay -input commands.txt example-garden.txt
```

#### Enable manual stepping

```sh
./gartentester -step -input commands.txt example-garden.txt
```

#### Generate with deterministic seed

```sh
./gartentester -generate -seed 42 garden.txt
```

## Building

### Prerequisites
- [Git](https://git-scm.com)
- [Go](https://go.dev)

```sh
git clone https://github.com/juho05/sdssh-gartentester
cd sdssh-gartentester
```

### For your current computer

```sh
go build
```

### For other computers

On a UNIX-like system (e.g. macOS/Linux):
```sh
./build.sh
```

## License

Copyright (c) 2023 Julian Hofmann

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program.  If not, see <http://www.gnu.org/licenses/>.
