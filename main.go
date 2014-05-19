package main

import (
	"fmt"

	"github.com/nejstastnejsistene/dotbot-go/activity"
)

func main() {
	fmt.Println(activity.IsDotsOnTop())
	d, err := activity.OpenDotsData("com.nerdyoctopus.gamedots.plist")
	fmt.Println(err)
	fmt.Println(d.Powerups())
	d.MaximizePowerups()
	fmt.Println(d.Powerups())
	err = d.Save()
	fmt.Println(err)
	err = d.Close()
	fmt.Println(err)
}
