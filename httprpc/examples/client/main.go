package main

import (
	"context"
	"fmt"
	"log"

	"git.ablecloud.cn/ablecloud/ac-comm-lib/httprpc"
	"git.ablecloud.cn/ablecloud/ac-comm-lib/httprpc/examples/arith"
)

func main() {
	c := httprpc.NewClient("http://localhost:8000", nil)

	args := arith.Args{A: 1, B: 4}
	reply := arith.Reply{}
	if err := c.Call(context.Background(), "/arith/v0/Add", args, &reply); err != nil {
		log.Fatalf("call: %v", err)
	}
	fmt.Printf("reply: %v\n", reply)
}
