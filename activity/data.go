package activity

import (
	"fmt"
	"io"

	plist "github.com/nejstastnejsistene/plist/xmlplist"
)

type DotsData map[string]interface{}

func ReadDotsData(r io.Reader) (data DotsData, err error) {
	err = plist.NewDecoder(r).Decode(&data)
	return
}

func (data DotsData) WriteTo(w io.Writer) error {
	v := map[string]interface{}(data)
	return plist.NewEncoder(w).Encode(v)
}

const MaxInt32 = int32(^uint32(0) >> 1)

func (data DotsData) Powerups() (timeFreezes, shrinkers, expanders int32) {
	timeFreezes = powerupToInt32(data["number_of_time_freezes"])
	shrinkers = powerupToInt32(data["number_of_shrinkers"])
	expanders = powerupToInt32(data["number_of_expanders"])
	return
}

func powerupToInt32(v interface{}) int32 {
	if v == nil {
		return 0
	}
	switch v := v.(type) {
	case uint8:
		return int32(v)
	case uint16:
		return int32(v)
	case uint32:
		return int32(v)
	case uint64:
		return int32(v)
	case int8:
		return int32(v)
	case int16:
		return int32(v)
	case int32:
		return v
	case int64:
		return int32(v)
	}
	panic(fmt.Sprintf("Unable to convert to int32: %#+v (type %T)\n", v, v))
	return 0
}

func (data DotsData) MaximizePowerups() {
	data.SetPowerups(MaxInt32, MaxInt32, MaxInt32)
}

func (data DotsData) SetPowerups(timeFreezes, shrinkers, expanders int32) {
	if timeFreezes > 0 {
		data["number_of_time_freezes"] = timeFreezes
	}
	if shrinkers > 0 {
		data["number_of_shrinkers"] = shrinkers
	}
	if expanders > 0 {
		data["number_of_expanders"] = expanders
	}
}
