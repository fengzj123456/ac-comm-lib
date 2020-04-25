# README

## Overview

`httpRPC` is a rpc framework, it's based on http.

## Example

service interface code
```
type Args struct {
	A, B int
}

type Reply struct {
	C int
}

type Arith int

func (t *Arith) Add(ctx context.Context, args Args, reply *Reply) error {
	log.Infof("[%s] Add Args: %v", httprpc.GetContextTraceID(ctx), args)
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
```

server code
```
func main() {
	var a arith.Arith
	s := httprpc.NewServer(nil)
	if err := s.Register(&a); err != nil {
		log.Fatalf("register: %v", err)
	}
	http.ListenAndServe(":8000", s)
}
```

client code
```
func main() {
	c := httprpc.NewClient("http://localhost:8000", nil)

	args := arith.Args{A: 1, B: 4}
	reply := arith.Reply{}
	if err := c.Call(context.Background(), "Arith.Add", args, &reply); err != nil {
		log.Fatalf("call: %v", err)
	}
	fmt.Printf("reply: %v\n", reply)
}
```

