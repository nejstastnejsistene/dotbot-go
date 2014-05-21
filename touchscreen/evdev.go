package touchscreen

import (
	"syscall"
)

type InputEvent struct {
	Time  syscall.Timeval
	Type  uint16
	Code  uint16
	Value int32
}

const (
	EV_SYN             uint16 = 0x00
	EV_KEY                    = 0x01
	EV_ABS                    = 0x03
	SYN_REPORT                = 0x00
	BTN_TOUCH                 = 0x14a
	ABS_X                     = 0x00
	ABS_Y                     = 0x01
	ABS_MT_POSITION_X         = 0x35
	ABS_MT_POSITION_Y         = 0x36
	ABS_MT_TRACKING_ID        = 0x39

	KeyUp   int32 = 0
	KeyDown       = 1
)
