package commons

import (
	"strconv"
	"time"
)

func UXSecs() string {
	return strconv.FormatInt(time.Now().Unix(), 10)
}

func UXNano() string {
	return strconv.FormatInt(time.Now().UnixNano(), 10)
}
