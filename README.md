# freefall

Landing humanitarian airdrops in hostile territory, control the parachute to avoid fire while you fall.

**Download or play in browser at: https://sinisterstuf.itch.io/freefall**

This is my entry to this year's [Nokia 3310 Jam](https://itch.io/jam/nokiajam5).

## For game testers

<!-- TODO: add a link to the latest downloads page -->

Game controls:
- F: toggle full-screen
- Q: quit the game
- Space / Numpad 5 / Tap screen: toggle parachute

## For programmers

Make sure you have [Go 1.19 or later](https://go.dev/) to contribute to the game

To build the game yourself, run: `go build .` it will produce an freefall file and on Windows freefall.exe.

To run the tests, run: `go test ./...` assuming there even are any.

The project has a very simple, flat structure, the first place to start looking is the main.go file.

## Attribution

This game was written using the [Ebitengine](https://ebitengine.org/) library. The graphics and animations were drawn in [Aseprite](https://www.aseprite.org/). The music was written in [LMMS](https://lmms.io/) using the nokia_3310_soundfont2.sf2 made by Krasno using samples imitating Nokia 3310 sounds made by Eamonn Watt.
