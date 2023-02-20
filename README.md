# freefall

A WIP 2D game for the 2023 Nokia 3310 Jam using the [Ebitengine](https://ebitengine.org/) library.

## For game testers

<!-- TODO: add a link to the latest downloads page -->

Game controls:
- F: toggle full-screen
- Q: quit the game
- Space: move up

## For programmers

Make sure you have [Go 1.19 or later](https://go.dev/) to contribute to the game

To build the game yourself, run: `go build .` it will produce an freefall file and on Windows freefall.exe.

To run the tests, run: `go test ./...` assuming there even are any.

The project has a very simple, flat structure, the first place to start looking is the main.go file.
