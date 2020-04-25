package httprpc

type errReply struct {
	Code  int
	Error string
	Cause string
	Stack string `json:",omitempty"`
}
