package activity

import (
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"regexp"

	"github.com/mkrautz/plist/xmlplist"
)

const PackageName = "com.nerdyoctopus.gamedots"

var onTopPattern = regexp.MustCompile(
	fmt.Sprintf("mCurrentFocus=.*%s.*\n", PackageName))

func IsDotsOnTop() (onTop bool, err error) {
	cmd := exec.Command("dumpsys", "window", "windows")
	output, err := cmd.Output()
	if err != nil {
		return
	}
	onTop = onTopPattern.Find(output) != nil
	return
}

func IsDotsRunning() (running bool, err error) {
	err = errors.New("not implemented")
	return
}

type DotsData struct {
	file  io.ReadWriteCloser
	plist map[string]interface{}
}

func OpenDotsData(filename string) (*DotsData, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	var plist map[string]interface{}
	err = xmlplist.NewDecoder(f).Decode(&plist)
	if err != nil {
		return nil, err
	}
	return &DotsData{f, plist}, err
}

func (d DotsData) Save() error {
	//return xmlplist.NewEncoder(d.file).Encode(d.plist)
	return errors.New("not implemented")
}

func (d DotsData) Close() error {
	return d.file.Close()
}

const MaxInt32 int32 = int32(^uint32(0) >> 1)

func (d DotsData) Powerups() (timeFreezes, shrinkers, expanders int32) {
	timeFreezes = convertToInt32(d.plist["number_of_time_freezes"])
	shrinkers = convertToInt32(d.plist["number_of_shrinkers"])
	expanders = convertToInt32(d.plist["number_of_expanders"])
	return
}

func convertToInt32(v interface{}) int32 {
	switch v := v.(type) {
	case int32:
		return v
	case uint8:
	case uint16:
	case uint32:
	case uint64:
	case int8:
	case int16:
	case int64:
		return int32(v)
	}
	panic(fmt.Sprintf("Unable to convert to int32: %#+v (type %T)\n", v, v))
	return 0
}

func (d DotsData) MaximizePowerups() {
	d.SetPowerups(MaxInt32, MaxInt32, MaxInt32)
}

func (d DotsData) SetPowerups(timeFreezes, shrinkers, expanders int32) {
	if timeFreezes > 0 {
		d.plist["number_of_time_freezes"] = timeFreezes
	}
	if shrinkers > 0 {
		d.plist["number_of_shrinkers"] = shrinkers
	}
	if expanders > 0 {
		d.plist["number_of_expanders"] = expanders
	}
}
