package httprpc

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"git.ablecloud.cn/ablecloud/ac-comm-lib/httprpc/codes"
)

type Args struct {
	A, B int
}

type Reply struct {
	C int
}

type Arith int

// Some of Arith's methods have value args, some have pointer args. That's deliberate.

func (t *Arith) Add(ctx context.Context, args Args, reply *Reply) error {
	reply.C = args.A + args.B
	return nil
}

func (t *Arith) Mul(ctx context.Context, args *Args, reply *Reply) error {
	reply.C = args.A * args.B
	return nil
}

func (t *Arith) Div(ctx context.Context, args Args, reply *Reply) error {
	if args.B == 0 {
		return errors.New("divide by zero")
	}
	reply.C = args.A / args.B
	return nil
}

func (t *Arith) String(ctx context.Context, args *Args, reply *string) error {
	*reply = fmt.Sprintf("%d+%d=%d", args.A, args.B, args.A+args.B)
	return nil
}

func (t *Arith) Scan(ctx context.Context, args string, reply *Reply) (err error) {
	_, err = fmt.Sscan(args, &reply.C)
	return
}

func (t *Arith) Error(ctx context.Context, args *Args, reply *Reply) error {
	panic("ERROR")
}

type BuiltinTypes struct{}

func (BuiltinTypes) Map(ctx context.Context, args *Args, reply *map[int]int) error {
	(*reply)[args.A] = args.B
	return nil
}

func (BuiltinTypes) Slice(ctx context.Context, args *Args, reply *[]int) error {
	*reply = append(*reply, args.A, args.B)
	return nil
}

func (BuiltinTypes) Array(ctx context.Context, args *Args, reply *[2]int) error {
	(*reply)[0] = args.A
	(*reply)[1] = args.B
	return nil
}

type hidden int

func (t *hidden) Exported(ctx context.Context, args Args, reply *Reply) error {
	reply.C = args.A + args.B
	return nil
}

type Embed struct {
	hidden
}

type Printer int

func (Printer) Println(ctx context.Context, s string, reply interface{}) error {
	fmt.Println(s)
	return nil
}

func TestServerRegisterCorrect(t *testing.T) {
	var arith Arith
	var embed Embed
	var builtinTypes BuiltinTypes
	tests := []struct {
		rcvr  interface{}
		name  string
		class string
	}{
		{rcvr: &arith, name: "/arith", class: "/arith"},
		{rcvr: &embed, name: "", class: "/"},
		{rcvr: builtinTypes, name: "/Builtin/", class: "/Builtin"},
	}

	var s Server
	for i, tt := range tests {
		if err := s.register(tt.rcvr, tt.name); err != nil {
			t.Fatalf("case%d: register: %v", i, err)
		}
		if _, ok := s.classes.Load(tt.class); !ok {
			t.Fatalf("case%d: %q not found", i, tt.class)
		}
	}
}

func TestServerRegisterError(t *testing.T) {
	var s Server

	// unexport
	var h hidden
	if err := s.register(&h, ""); err == nil {
		t.Errorf("register return error is nil")
	} else {
		t.Logf("register: %v", err)
	}

	// no method
	type A struct{}
	if err := s.register(A{}, ""); err == nil {
		t.Errorf("register return error is nil")
	} else {
		t.Logf("register: %v", err)
	}

	// repeat register
	var a Arith
	if err := s.register(&a, ""); err != nil {
		t.Fatalf("register: %v", err)
	}
	if err := s.register(&a, ""); err == nil {
		t.Errorf("register return error is nil")
	} else {
		t.Logf("register: %v", err)
	}
}

func getMethod(rcvr interface{}, methodName string) *method {
	m, ok := reflect.TypeOf(rcvr).MethodByName(methodName)
	if !ok {
		panic(fmt.Errorf("%s method not found", methodName))
	}
	return &method{
		method: m,
		args:   m.Type.In(2),
		reply:  m.Type.In(3),
	}
}

func TestLookupMethod(t *testing.T) {
	var a Arith
	var b BuiltinTypes
	var s Server
	if err := s.Register("/arith", &a); err != nil {
		t.Fatalf("Register: %v", err)
	}
	if err := s.Register("/builtin", &b); err != nil {
		t.Fatalf("Register: %v", err)
	}

	tests := []struct {
		class  string
		method string
		rcvr   reflect.Value
		meth   *method
	}{
		{class: "/arith", method: "Add", rcvr: reflect.ValueOf(&a), meth: getMethod(&a, "Add")},
		{class: "/arith", method: "Div", rcvr: reflect.ValueOf(&a), meth: getMethod(&a, "Div")},
		{class: "/builtin", method: "Map", rcvr: reflect.ValueOf(&b), meth: getMethod(&b, "Map")},
		{class: "/builtin", method: "Slice", rcvr: reflect.ValueOf(&b), meth: getMethod(&b, "Slice")},
	}
	for _, tt := range tests {
		rcvr, meth, err := s.lookupMethod(tt.class, tt.method)
		if err != nil {
			t.Fatalf("lookupMethod(%s, %s): %v", tt.class, tt.method, err)
		}
		if got, want := rcvr, tt.rcvr; got != want {
			t.Fatalf("lookupMethod(%s, %s): rcvr: got %v, want %v", tt.class, tt.method, got, want)
		}
		if got, want := meth.method.Name, tt.meth.method.Name; got != want {
			t.Fatalf("lookupMethod(%s, %s): method name: got %v, want %v", tt.class, tt.method, got, want)
		}
		if got, want := meth.args, tt.meth.args; got != want {
			t.Fatalf("lookupMethod(%s, %s): method args: got %v, want %v", tt.class, tt.method, got, want)
		}
		if got, want := meth.reply, tt.meth.reply; got != want {
			t.Fatalf("lookupMethod(%s, %s): method reply: got %v, want %v", tt.class, tt.method, got, want)
		}
	}
}

func serveTestHTTP(h http.Handler, method, path string, b []byte) (*httptest.ResponseRecorder, error) {
	r := httptest.NewRequest(method, path, bytes.NewReader(b))
	w := httptest.NewRecorder()
	h.ServeHTTP(w, r)
	return w, nil
}

func callTestServer(s *Server, path string, args interface{}, reply interface{}) error {
	var err error
	var buf bytes.Buffer
	if args != nil {
		if err = s.codec.Encode(&buf, args); err != nil {
			return err
		}
	}

	r, err := serveTestHTTP(s, "POST", path, buf.Bytes())
	if err != nil {
		return err
	}
	if r.Code != http.StatusOK {
		var er errReply
		if err = s.codec.Decode(r.Body, &er); err != nil {
			return err
		}
		return Errorf(codes.Code(er.Code), er.Cause)
	}

	if reply != nil {
		if err = s.codec.Decode(r.Body, reply); err != nil {
			return err
		}
	}
	return nil
}

func TestServeHTTP(t *testing.T) {
	var a Arith
	var b BuiltinTypes
	var e Embed
	s := NewServer(nil)
	if err := s.Register("/arith", &a); err != nil {
		t.Fatalf("Register: %v", err)
	}
	if err := s.Register("/builtin", b); err != nil {
		t.Fatalf("Register: %v", err)
	}
	if err := s.Register("/embed", &e); err != nil {
		t.Fatalf("Register: %v", err)
	}

	s1, s2 := "1+2=3", ""
	tests := []struct {
		path   string
		args   interface{}
		result interface{}
		reply  interface{}
	}{
		{path: "/arith/Add", args: Args{A: 1, B: 2}, result: &Reply{C: 3}, reply: &Reply{}},
		{path: "/arith/Div", args: Args{A: 1, B: 2}, result: &Reply{C: 0}, reply: &Reply{}},
		{path: "/arith/String", args: Args{A: 1, B: 2}, result: &s1, reply: &s2},
		{path: "/builtin/Map", args: Args{A: 1, B: 2}, result: &map[int]int{1: 2}, reply: &map[int]int{}},
		{path: "/embed/Exported", args: Args{A: 1, B: 2}, result: &Reply{C: 3}, reply: &Reply{}},
	}
	for _, tt := range tests {
		if err := callTestServer(s, tt.path, tt.args, tt.reply); err != nil {
			t.Fatalf("callTestServer(%s): %v", tt.path, err)
		}
		if got, want := tt.reply, tt.result; !reflect.DeepEqual(got, want) {
			t.Fatalf("callTestServer(%s): reply: got %v, want %v", tt.path, got, want)
		}
		t.Logf("%s: args: %v, reply: %v", tt.path, tt.args, tt.reply)
	}
}

func TestServerMiddleware(t *testing.T) {
	var p Printer
	s := NewServer(nil)
	if err := s.Register("/", p); err != nil {
		t.Fatalf("Register: %v", err)
	}

	m1 := func(ctx context.Context, w http.ResponseWriter, r *http.Request, next NextMiddleware) error {
		fmt.Println("m1 before")
		err := next(ctx, w, r)
		fmt.Println("m1 after")
		return err
	}
	m2 := func(ctx context.Context, w http.ResponseWriter, r *http.Request, next NextMiddleware) error {
		fmt.Println("m2 before")
		err := next(ctx, w, r)
		fmt.Println("m2 after")
		return err
	}
	m3 := func(ctx context.Context, w http.ResponseWriter, r *http.Request, next NextMiddleware) error {
		fmt.Println("m3 before")
		err := next(ctx, w, r)
		fmt.Println("m3 after")
		return err
	}
	s.AddMiddleware(MiddlewareFunc(m1), MiddlewareFunc(m2))
	s.AddMiddleware(MiddlewareFunc(m2), MiddlewareFunc(m1), MiddlewareFunc(m3))

	if err := callTestServer(s, "/Println", "Hello, world!", nil); err != nil {
		t.Fatalf("callTestServer: %v", err)
	}
}
