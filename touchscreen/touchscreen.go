package touchscreen

import (
	"encoding/binary"
	"errors"
	"fmt"
	"image"
	"os"
	"time"
)

const (
	ShortDelay time.Duration = 10 * time.Millisecond // A short delay between groups of touch events.
	LongDelay                = 10 * ShortDelay       // A stylish delay before a gesture is ended.
	NumPoints  int           = 5                     // How many points are interpolated for gestures.
)

type TouchScreen struct {
	Device           *os.File         // The device node to write to.
	PixelSize        *image.Rectangle // The size in pixels of the display.
	*TouchScreenInfo                  // Additional info about the display.
}

func New(filename string, size *image.Rectangle) (v *TouchScreen, err error) {
	dev, err := os.OpenFile(filename, os.O_WRONLY, 0666)
	if err != nil {
		return
	}
	info, err := DeviceInfo(filename)
	if err != nil {
		return
	}
	v = &TouchScreen{dev, size, info}
	return
}

func (v *TouchScreen) SetPixelSize(size image.Rectangle) {
	v.PixelSize = &size
}

// Write an input event to the device.
func (v TouchScreen) event(evType, code uint16, value int) {
	event := InputEvent{
		Type:  evType,
		Code:  code,
		Value: int32(value),
	}
	err := binary.Write(v.Device, binary.LittleEndian, event)
	if err != nil {
		panic(err)
	}
	if evType == EV_SYN && code == SYN_REPORT {
		time.Sleep(ShortDelay)
	}
}

// Write a sync event.
func (v TouchScreen) sync() {
	v.event(EV_SYN, SYN_REPORT, 0)
}

// Set the position of the cursor.
func (v TouchScreen) setPos(p image.Point) {
	switch v.Type {
	case SingleTouch:
		v.event(EV_ABS, ABS_X, p.X)
		v.event(EV_ABS, ABS_Y, p.Y)
		return
	case MultiTouch:
		v.event(EV_ABS, ABS_MT_POSITION_X, p.X)
		v.event(EV_ABS, ABS_MT_POSITION_Y, p.Y)
		return
	}
	panic(fmt.Sprint("touchscreen: not implemented:", v.Type))
}

var trackingId int = 0

// Indicate that a touch event is beginning.
func (v TouchScreen) fingerDown(p image.Point) {
	switch v.Type {
	case SingleTouch:
		v.setPos(p)
		v.event(EV_KEY, BTN_TOUCH, 1)
		v.sync()
		return
	case MultiTouch:
		v.event(EV_ABS, ABS_MT_TRACKING_ID, trackingId)
		v.setPos(p)
		v.sync()
		trackingId++
		return
	}
	panic(fmt.Sprint("touchscreen: not implemented:", v.Type))
}

// Indicate that a touch event is finished.
func (v TouchScreen) fingerUp() {
	switch v.Type {
	case SingleTouch:
		v.event(EV_KEY, BTN_TOUCH, 0)
		v.sync()
		return
	case MultiTouch:
		v.event(EV_ABS, ABS_MT_TRACKING_ID, -1)
		v.sync()
		return
	}
	panic(fmt.Sprint("touchscreen: not implemented:", v.Type))
}

// Tap the screen.
func (v TouchScreen) Tap(p image.Point) error {
	ps := make([]image.Point, 1)
	ps[0] = p
	return v.Gesture(ps)
}

// Perform a touch gesture, indicated by the waypoints ps.
// NumPoints points are interpolated linearly between the given points.
// It is an error to pass a nil or empty slice as the argument.
func (v TouchScreen) Gesture(ps []image.Point) (err error) {
	if ps == nil || len(ps) == 0 {
		return errors.New("touchscreen: expecting at least one point.")
	}
	if v.PixelSize == nil {
		return errors.New("touchscreen: must call SetPixelSize() first")
	}
	for i, p := range ps {
		ps[i] = v.coord(p)
	}
	func() {
		// For ease of writing, the touch event code uses panics to indicate
		// errors. We recover from that here, and return it as an error.
		defer func() {
			r := recover()
			if r != nil {
				err = errors.New(fmt.Sprint(r))
			}
		}()
		v.fingerDown(ps[0])
		prev := ps[0]
		for _, p := range ps[1:] {
			// Interpolate points up to and including p.
			for j := 1; j <= NumPoints; j++ {
				v.setPos(prev.Add(p.Sub(prev).Mul(j).Div(NumPoints)))
				v.sync()
			}
			prev = p
		}
		time.Sleep(LongDelay)
		v.fingerUp()
	}()
	return
}

// Scale coordinates from pixels to the touchscreen.
// Size is the size of the screen in pixels.
func (v TouchScreen) coord(p image.Point) image.Point {
	size := v.PixelSize
	vSize := v.VirtualSize
	return image.Point{
		vSize.Min.X + (p.X-size.Min.X)*vSize.Max.X/size.Max.X,
		vSize.Min.Y + (p.Y-size.Min.Y)*vSize.Max.Y/size.Max.Y,
	}
}
