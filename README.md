# SDSSH Challenge 2023 - Garten Aufräumen - Testing Tool

A testing tool and simulator for the [SDSSH Challenge 2023](https://github.com/maxwellmatthis/sdssh-challenge-2023).

## Installation

Download the executable from the [releases](https://github.com/juho05/sdssh-gartentester/releases/latest) page.

## Usage

*The following instructions are for macOS and Linux. You might have to tweak the syntax a bit to work on Windows.

### Pipe the output of your programm into the testing tool

```sh
./my-program | ./gartentester example-garden.txt
```

### Read commands from a file

```sh
./gartentester -input commands.txt example-garden.txt
```

### Disable animation (faster)

```sh
./gartentester -no-delay -input commands.txt example-garden.txt
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