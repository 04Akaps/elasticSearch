package protocol

type status string

const (
	successStatus = status("success")
	failedStatus  = status("failed")
	errorStatus   = status("error")
)
