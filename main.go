package main

import (
	"log"
	"net/http"
	"raft-example/raft"
	"raft-example/service"
)

func main() {
	node := raft.NewNode("192.168.56.1:6060")
	node.Start(true)
	store := service.NewStore(node.GetRaftInstince())
	handel := service.NewWeb(store)
	if err := http.ListenAndServe(":8080", handel); err != nil {
		log.Fatal("HTTP服务启动失败:", err)
	}
}
