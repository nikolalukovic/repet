package server

import "fmt"

type RepetErrorCode int16

const (
	ConfigPortNotSet RepetErrorCode = iota
	ConfigAddrNotSet

	CommandNotFound

	UnsupportedMessageVersion
)

type RepetError struct {
	Code    RepetErrorCode
	Details string
}

func (re *RepetError) Error() string {
	switch re.Code {
	case ConfigAddrNotSet:
		return "REPET_ADDR not set as an environment variable"
	case ConfigPortNotSet:
		return "REPET_PORT not set as an environemnt variable"
	case CommandNotFound:
		return fmt.Sprintf("Specified %s command not found", re.Details)
	case UnsupportedMessageVersion:
		return fmt.Sprintf("Passed %s version is not supported", re.Details)
	default:
		panic("Repet error code not handled")
	}
}
