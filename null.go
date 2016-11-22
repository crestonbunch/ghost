package ghost

import (
	"net/http"
)

// A null model doesn't do anything
type NullModel struct {
}

func (m *NullModel) FromRequest(r *http.Request) HttpError {
	return nil
}

// A null validator doesn't do anything
type NullValidator struct {
}

func (v *NullValidator) Validate(interface{}) HttpError {
	return nil
}

// A null processor doesn't do anything
type NullProcessor struct {
}

func (p *NullProcessor) Process(interface{}) (interface{}, HttpError) {
	return nil, nil
}

// A null writer doesn't do anything
type NullWriter struct {
}

func (w *NullWriter) Serialize(interface{}) ([]byte, HttpError) {
	return []byte{}, nil
}
