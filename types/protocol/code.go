package protocol

type code int32

const (
	success code = iota
	queryParsingFailed
	bodyParsingFailed
)
