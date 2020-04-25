package httprpc

import (
	"encoding/json"
	"testing"
)

func TestErrReply(t *testing.T) {
	er := errReply{
		Code:  1,
		Error: "test error",
		Cause: "test cause",
	}
	data, err := json.Marshal(er)
	if err != nil {
		t.Fatalf("json marshal: %v", err)
	}
	t.Logf("data: %s", data)
}
