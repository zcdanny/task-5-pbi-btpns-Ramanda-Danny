package main

import (
	"net/http"

	"github.com/zcdanny/task-5-pbi-btpns-Ramanda-Danny/config"
	"github.com/zcdanny/task-5-pbi-btpns-Ramanda-Danny/router"
)

func main() {
	config.Load()
	r := router.Init()
	r.StaticFS("/uploads", http.Dir("uploads"))
	r.Run(":8080")
}
