package util

import (
	"fmt"
	"strings"
	"time"
)

func ParseBuildTime(buildName string) (time.Time, error) {
	if !strings.HasPrefix(buildName, "builds_") {
		return time.Time{}, fmt.Errorf("input doesn't include the right prefix")
	}

	timePart := strings.SplitN(buildName, "_", 2)[1]
	form := "2006_01_02"

	return time.Parse(form, timePart)
}
