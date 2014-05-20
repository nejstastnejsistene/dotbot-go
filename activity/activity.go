package activity

import (
	"errors"
	"fmt"
	"os/exec"
	"regexp"
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
