dotbot-go [![Build Status](https://travis-ci.org/nejstastnejsistene/dotbot-go.svg?branch=master)](https://travis-ci.org/nejstastnejsistene/dotbot-go)
=========

I originally wrote DotBot for PennApps Fall 2013, and then rewrote it in C so that it could compute possible moves in a reasonable amount of time. It was much faster, but it ended up being a mish-mash two C programs communicating over adb via python scripts. That code is available [here](https://github.com/nejstastnejsistene/DotBot).

This is DotBot ported to Go. It was for fun, to teach myself Go, to explore continuous integration, and to see how how quickly it runs in comparison to the original C code. It turns out it is slower, but only marginally so. At worst, the moves are calculated within a few seconds on my Nexus 7. This trade-off in speed is justified by the ease and joy of development that comes with writing in Go.

### Run on Infinite Mode

First make sure Dots is open to Infinite mode, then run the following commands. Tested on a 2013 Nexus 7 with the dark color scheme.

```sh
./build.sh
adb push dotbot-go /data/local
adb shell /data/local/dotbot-go
```
