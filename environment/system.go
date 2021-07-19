package environment

import "runtime"

type OsEnum int32

const (
	OS_UNKNOW OsEnum = iota
	OS_Win
	OS_LINUX
	OS_DARWIN
)

func WhichOS() OsEnum {
	switch runtime.GOOS {
	case "darwin":
		return OS_DARWIN
	case "windows":
		return OS_Win
	case "linux":
		return OS_LINUX
	default:
		return OS_UNKNOW
	}
}

func (o OsEnum) ToInt32() int32 {
	return int32(o)
}
