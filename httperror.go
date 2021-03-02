package gcp

type HttpError struct {
	Code int   `json:"code"`
	Err  error `json:"error"`
}

func (h *HttpError) Error() string {
	if h.Err != nil {
		return h.Err.Error()
	}

	return ""
}

func (h *HttpError) IsStatusCode(err error, code int) bool {
	if e, ok := err.(*HttpError); ok {
		return e.Code == code
	}

	return false
}
