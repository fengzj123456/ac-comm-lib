package httprpc

import (
	"encoding/json"
	"io"
)

type Codec interface {
	ContentType() string
	Encode(io.Writer, interface{}) error
	Decode(io.Reader, interface{}) error
}

type JSONCodec struct{}

func (JSONCodec) ContentType() string {
	return "application/json"
}

func (JSONCodec) Encode(w io.Writer, v interface{}) error {
	return json.NewEncoder(w).Encode(v)
}

func (JSONCodec) Decode(r io.Reader, v interface{}) error {
	return json.NewDecoder(r).Decode(v)
}

var DefaultCodec = JSONCodec{}
