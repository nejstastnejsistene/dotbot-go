package activity

import (
	"fmt"
	"os/exec"
	"regexp"
)

const PackageName = "com.nerdyoctopus.gamedots"

var onTopPattern = regexp.MustCompile(
	fmt.Sprintf("mCurrentFocus=.*%s.*\n", PackageName))

func IsOnTop() (onTop bool, err error) {
	cmd := exec.Command("dumpsys", "window", "windows")
	output, err := cmd.Output()
	if err != nil {
		return
	}
	onTop = onTopPattern.Find(output) != nil
	return
}
