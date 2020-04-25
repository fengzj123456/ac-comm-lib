package httprpc

import (
	"context"
	"net/http/httptest"
	"reflect"
	"testing"
)

var testServerURL string

func TestClient(t *testing.T) {
	c := NewClient(testServerURL, nil)

	s1, s2 := "", "1+2=3"
	tests := []struct {
		classMethod string
		args        interface{}
		reply       interface{}
		result      interface{}
	}{
		{classMethod: "arith/Add", args: Args{A: 1, B: 2}, reply: &Reply{}, result: &Reply{C: 3}},
		{classMethod: "arith/Mul", args: Args{A: 1, B: 2}, reply: &Reply{}, result: &Reply{C: 2}},
		{classMethod: "arith/Div", args: Args{A: 8, B: 2}, reply: &Reply{}, result: &Reply{C: 4}},
		{classMethod: "arith/String", args: Args{A: 1, B: 2}, reply: &s1, result: &s2},
	}
	for _, tt := range tests {
		if err := c.Call(context.Background(), tt.classMethod, tt.args, tt.reply); err != nil {
			t.Fatalf("Call(%s): %v", tt.classMethod, err)
		}
		if got, want := tt.reply, tt.result; !reflect.DeepEqual(got, want) {
			t.Fatalf("Call(%s): reply: got %v, want %v", tt.classMethod, got, want)
		}
		t.Logf("Call(%s): reply: %v", tt.classMethod, reflect.ValueOf(tt.reply).Elem().Interface())
	}
}

func TestMain(m *testing.M) {
	var a Arith
	s := NewServer(nil)
	if err := s.Register("/arith", &a); err != nil {
		panic(err)
	}
	svr := httptest.NewServer(s)
	defer svr.Close()
	testServerURL = svr.URL
	m.Run()
}
