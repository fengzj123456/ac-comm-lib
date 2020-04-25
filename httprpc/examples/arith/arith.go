package arith

import (
	"context"
	"errors"
	"fmt"
	"log"

	"git.ablecloud.cn/ablecloud/ac-comm-lib/httprpc"
)

type Args struct {
	A, B int
}

type Reply struct {
	C int
}

type Arith int

// Some of Arith's methods have value args, some have pointer args. That's deliberate.

func (t *Arith) Add(ctx *httprpc.Context, args Args, reply *Reply) error {
	log.Printf("[%s] Add Args: %v", ctx.TraceID, args)
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

func (t *Arith) Error(ctx context.Context, args interface{}, reply *Reply) error {
	panic("ERROR")
}
