package protocol

type response struct {
	*inner
	Data interface{} `json:"data"`
}

func Response(data interface{}, inner *inner) *response {
	r := &response{
		inner: inner,
		Data:  data,
	}

	if v, ok := r.Data.(error); ok {
		r.Data = v.Error()
	}

	if r.Data == nil || r.Data == "" {
		r.Data = inner.Message
	}

	return r
}

func (r response) Error() string {
	if str, ok := r.Data.(string); ok {
		return str
	}

	if err, ok := r.Data.(error); ok {
		return err.Error()
	}

	return ""
}
