package httprpc

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"reflect"
	"sync"

	"git.ablecloud.cn/ablecloud/ac-comm-lib/httprpc/codes"
)

type Server struct {
	codec   Codec
	classes sync.Map
	next    NextMiddleware
}

func NewServer(codec Codec) *Server {
	if codec == nil {
		codec = DefaultCodec
	}
	return &Server{codec: codec}
}

func (s *Server) Register(prefix string, rcvr interface{}) error {
	return s.register(rcvr, prefix)
}

func (s *Server) register(rcvr interface{}, name string) error {
	val := reflect.ValueOf(rcvr)
	tname := reflect.Indirect(val).Type().Name()
	if !isExported(tname) {
		return fmt.Errorf("register: type %s is not exported", tname)
	}

	name = normalizePath(name)
	c, err := parseClass(name, val)
	if err != nil {
		return fmt.Errorf("register: parse class: %v", err)
	}

	if _, loaded := s.classes.LoadOrStore(name, c); loaded {
		return fmt.Errorf("register: class already defined: %s", name)
	}
	return nil
}

func (s *Server) AddMiddleware(middlewares ...Middleware) {
	if s.next == nil {
		s.next = s.serveHTTP
	}
	for _, m := range middlewares {
		h := m
		n := s.next
		s.next = func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
			return h.ServeHTTP(ctx, w, r, n)
		}
	}
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method == "OPTIONS" {
		setCors(w.Header(), r.Header.Get("Origin"))
		return
	}

	next := s.next
	if next == nil {
		next = s.serveHTTP
	}
	ctx := &Context{
		Context:  context.Background(),
		TraceID:  getHeaderTraceID(r.Header),
		Request:  r,
		Response: w,
	}
	if err := next(ctx, w, r); err != nil {
		s.setError(w, err, r)
	}
}

func (s *Server) serveHTTP(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	// check
	if err := s.checkContentType(r.Header); err != nil {
		return NewError(codes.InvalidHeader, err)
	}

	// lookup method
	rcvr, meth, err := s.lookupByPath(r.URL.Path)
	if err != nil {
		return NewError(codes.InvalidPath, err)
	}

	// decode args
	args, err := s.decodeArgs(r.Body, meth.args)
	if err != nil {
		return NewError(codes.DecodeBodyFail, err)
	}

	// call method
	reply := newReply(meth.reply)
	if err = call(ctx, meth.method, rcvr, args, reply); err != nil {
		return err
	}

	// set response header
	s.setResponseHeader(w, ctx, r)

	// encode reply
	if !isNilInterface(meth.reply) {
		if err = s.codec.Encode(w, reply.Interface()); err != nil {
			return NewError(codes.EncodeBodyFail, err)
		}
	}

	return nil
}

func (s *Server) checkContentType(h http.Header) error {
	if v := getHeaderContentType(h); v != "" && v != s.codec.ContentType() {
		return fmt.Errorf("not support %s Content-Type header, the Content-Type header must be %s or empty", v, s.codec.ContentType())
	}
	return nil
}

func (s *Server) lookupByPath(path string) (rcvr reflect.Value, meth *method, err error) {
	className, methodName := splitPath(path)
	rcvr, meth, err = s.lookupMethod(className, methodName)
	if err != nil {
		return reflect.Value{}, nil, fmt.Errorf("lookup method: %v", err)
	}
	return rcvr, meth, err
}

func (s *Server) lookupMethod(className, methodName string) (reflect.Value, *method, error) {
	v, ok := s.classes.Load(className)
	if !ok {
		return reflect.Value{}, nil, fmt.Errorf("can not find class %s/%s", className, methodName)
	}
	c := v.(*class)
	meth, ok := c.methods[methodName]
	if !ok {
		return reflect.Value{}, nil, fmt.Errorf("can not find method %s/%s", className, methodName)
	}
	return c.rcvr, meth, nil
}

func (s *Server) decodeArgs(r io.Reader, argsType reflect.Type) (args reflect.Value, err error) {
	isValue := false
	if argsType.Kind() == reflect.Ptr {
		args = reflect.New(argsType.Elem())
	} else {
		args = reflect.New(argsType)
		isValue = true
	}
	if !isNilInterface(argsType) {
		if err = s.codec.Decode(r, args.Interface()); err != nil {
			return args, err
		}
	}
	if isValue {
		args = args.Elem()
	}
	return args, nil
}

func newReply(replyType reflect.Type) (reply reflect.Value) {
	if isNilInterface(replyType) {
		reply = reflect.New(replyType)
	} else {
		reply = reflect.New(replyType.Elem())
		switch replyType.Elem().Kind() {
		case reflect.Map:
			reply.Elem().Set(reflect.MakeMap(replyType.Elem()))
		case reflect.Slice:
			reply.Elem().Set(reflect.MakeSlice(replyType.Elem(), 0, 0))
		}
	}
	return reply
}

func call(ctx context.Context, method reflect.Method, rcvr, args, reply reflect.Value) (err error) {
	defer func() {
		if r := recover(); r != nil {
			buf := runtimeStack()
			log.Printf("panic: %v\n%s", r, buf)

			if e, ok := r.(error); ok {
				err = e
			} else {
				err = fmt.Errorf("%v", r)
			}
			err = NewError(codes.Panic, err)
		}
	}()

	rets := method.Func.Call([]reflect.Value{rcvr, reflect.ValueOf(ctx), args, reply})
	erri := rets[0].Interface()
	if erri != nil {
		err = erri.(error)
	}
	return err
}

func (s *Server) setResponseHeader(w http.ResponseWriter, ctx context.Context, r *http.Request) {
	setHeaderContentType(w.Header(), s.codec.ContentType())
	setCors(w.Header(), r.Header.Get("Origin"))
	if rctx, ok := ctx.(*Context); ok {
		setHeaderTraceID(w.Header(), rctx.TraceID)
		for key, values := range rctx.ResponseHeader {
			for _, value := range values {
				w.Header().Add(key, value)
			}
		}
	}
}

func (s *Server) setError(w http.ResponseWriter, err error, r *http.Request) {
	code, cause, stack := GetErrorCode(err), GetErrorCause(err), GetErrorStack(err)
	er := errReply{
		Code:  int(code),
		Error: code.String(),
		Cause: cause.Error(),
		Stack: string(stack),
	}
	setHeaderContentType(w.Header(), s.codec.ContentType())
	setCors(w.Header(), r.Header.Get("Origin"))
	w.WriteHeader(code.Status())
	s.codec.Encode(w, er)
}
