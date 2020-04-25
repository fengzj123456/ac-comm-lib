package codes

import (
	"fmt"
	"net/http"
)

type Code int

type entry struct {
	desc   string
	status int
}

var codes = map[Code]entry{}

func Register(code Code, desc string, status int) {
	if _, ok := codes[code]; ok {
		panic(fmt.Sprintf("code(%d) is registered", code))
	}
	codes[code] = entry{desc: desc, status: status}
}

func RegisterDesc(code Code, desc string) {
	Register(code, desc, http.StatusInternalServerError)
}

func (c Code) String() string {
	if e, ok := codes[c]; ok {
		return e.desc
	}
	return fmt.Sprintf("code(%d)", c)
}

func (c Code) Status() int {
	if e, ok := codes[c]; ok {
		return e.status
	}
	return http.StatusInternalServerError
}

const (
	OK      Code = 0
	Unknown Code = -1
	Panic   Code = -2

	InvalidPath   Code = -101
	InvalidHeader Code = -102

	EncodeBodyFail Code = -201
	DecodeBodyFail Code = -202
)

func init() {
	Register(OK, "ok", http.StatusOK)
	Register(Unknown, "unknown error", http.StatusInternalServerError)
	Register(Panic, "panic error", http.StatusInternalServerError)

	Register(InvalidPath, "invalid url path", http.StatusBadRequest)
	Register(InvalidHeader, "invalid http header", http.StatusBadRequest)

	Register(EncodeBodyFail, "encode http body fail", http.StatusInternalServerError)
	Register(DecodeBodyFail, "decode http body fail", http.StatusBadRequest)
}
