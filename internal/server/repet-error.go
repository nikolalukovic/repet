package server

type RepetErrorCode int16

const (
	ConfigPortNotSet RepetErrorCode = iota
	ConfigAddrNotSet RepetErrorCode = iota
)

type RepetError struct {
	Code RepetErrorCode
}

func (re *RepetError) Error() string {
	switch re.Code {
	case ConfigAddrNotSet:
		return "REPET_ADDR not set as an environment variable"
	case ConfigPortNotSet:
		return "REPET_PORT not set as an environemnt variable"
	default:
		panic("Repet error code not handled")
	}
}
