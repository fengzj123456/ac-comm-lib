package main

import (
	"log"
	"net/http"

	"git.ablecloud.cn/ablecloud/ac-comm-lib/httprpc"
	"git.ablecloud.cn/ablecloud/ac-comm-lib/httprpc/examples/arith"
)

func main() {
	//httprpc.StackTrace = true

	var a arith.Arith
	s := httprpc.NewServer(nil)
	if err := s.Register("/arith/v0", &a); err != nil {
		log.Fatalf("register: %v", err)
	}
	http.ListenAndServe(":8000", s)
}
