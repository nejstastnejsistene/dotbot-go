package touchscreen

import (
	"bytes"
	"errors"
	"image"
	"os/exec"
	"regexp"
	"strconv"
)

type TouchScreenInfo struct {
	Type        TouchScreenType
	VirtualSize image.Rectangle
}

type TouchScreenType string

const (
	SingleTouch TouchScreenType = "SingleTouch"
	MultiTouch                  = "MultiTouch"
)

var (
	geteventPattern = regexp.MustCompile(`^add device \d+: (.*)$`)
	absPattern      = regexp.MustCompile(`^.{16}([0-9a-f]{4}).*min (\d+), max (\d+)`)
	eventsSection   = []byte("  events:")
)

// Determine a touch screen's type and capabilities from the output of
// getevent -p.
func DeviceInfo(filename string) (*TouchScreenInfo, error) {
	cmd := exec.Command("getevent", "-p")
	buf, err := cmd.Output()
	if err != nil {
		return nil, err
	}
	return deviceInfo(filename, buf)
}

// Determine information about a touch screen given the filename and the output
// getevent -p.
func deviceInfo(filename string, buf []byte) (*TouchScreenInfo, error) {
	lines := make(chan []byte)
	go grepSection(filename, buf, lines)
	for line := range lines {
		if bytes.Equal(line, eventsSection) {
			return readCapabilities(lines)
		}
	}
	return nil, errors.New("touchscreen: unable to determine screen type")
}

// TODO This is a little sloppy. It works on my Nexus, but could probably be
// rewritten to more gracefully handle other devices.
func readCapabilities(lines chan []byte) (*TouchScreenInfo, error) {
	info := new(TouchScreenInfo)
	for line := range lines {
		m := absPattern.FindSubmatch(line)
		if len(m) != 4 {
			break
		}
		code, err := strconv.ParseUint(string(m[1]), 16, 16)
		min, err := strconv.ParseInt(string(m[2]), 10, 32)
		max, err := strconv.ParseInt(string(m[3]), 10, 32)
		if err != nil {
			return nil, err
		}
		if info.Type == "" {
			switch uint16(code) {
			case ABS_X, ABS_Y:
				info.Type = SingleTouch
			case ABS_MT_POSITION_X, ABS_MT_POSITION_Y, ABS_MT_TRACKING_ID:
				info.Type = MultiTouch
			}
		}
		switch uint16(code) {
		case ABS_X, ABS_MT_POSITION_X:
			info.VirtualSize.Min.X = int(min)
			info.VirtualSize.Max.X = int(max)
		case ABS_Y, ABS_MT_POSITION_Y:
			info.VirtualSize.Min.Y = int(min)
			info.VirtualSize.Max.Y = int(max)
		}
	}
	return info, nil
}

// Filter out the lines of getevent -p that are relevant to our file.
func grepSection(filename string, buf []byte, lines chan []byte) {
	correctSection := false
	for _, line := range bytes.Split(buf, []byte{'\n'}) {
		m := geteventPattern.FindSubmatch(line)
		if len(m) == 2 {
			correctSection = bytes.Equal(m[1], []byte(filename))
		}
		if correctSection {
			lines <- line
		}
	}
	close(lines)
}
