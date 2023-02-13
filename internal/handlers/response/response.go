package response

type Response struct {
	contentType string
	body        []byte
	status      int
	errors      []error
}

func NewResponse(contentType string) *Response {
	return &Response{contentType: contentType, errors: make([]error, 0)}
}

func (r *Response) AddError(err error) *Response {
	r.errors = append(r.errors, err)
	return r
}

func (r *Response) SetStatus(status int) *Response {
	r.status = status
	return r
}

func (r *Response) SetBody(body []byte) *Response {
	r.body = body
	return r
}

func (r *Response) GetContentType() string {
	return r.contentType
}

func (r *Response) GetStatus() int {
	return r.status
}

func (r *Response) GetBody() []byte {
	return r.body
}
