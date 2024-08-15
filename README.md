# SRM Manifest Generator

This is a [Steam ROM Manager](https://github.com/SteamGridDB/steam-rom-manager) manifest generator
made in _golang_. It's over engineered with concurrency but it gets the job done 😅

The main reason I made this, is because I want Steam ROM Manager to be able to find my Pokemon Fan
Games that aren't on any major platform.

Powered by [huh](https://github.com/charmbracelet/huh) and [log](https://github.com/charmbracelet/log).

## Usage

Build the program using `make`, run:

```sh
make
```

The program will be compiled to /bin, so to run it you can run

```sh
./bin/srm-manifest-gen
```

After that proceed with the program instructions, it will find every executable (.exe) and make you
choose it.
