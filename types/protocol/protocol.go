package protocol

type inner struct {
	Code    code   `json:"code"`
	Message string `json:"message"`
	Status  status `json:"status"`
}

func newInner(c code, msg string, status status) *inner {
	return &inner{c, msg, status}
}

var (
	SUCCESS            = newInner(success, "", successStatus)
	FailedQueryParsing = newInner(queryParsingFailed, "Failed to parsing uri query", failedStatus)
	FailedBodyParsing  = newInner(bodyParsingFailed, "Failed to parsing in body", failedStatus)
)
