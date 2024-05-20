package main

import (
	"fmt"
	"log"
	"main/api/rest"
	"main/core"
	"main/core/node"
	"main/core/node/cfgclients"
	"net/http"
	"time"
)

func main() {
	zkc, err := cfgclients.NewZKCfg(
		core.Config.ZKHosts,
		time.Second*time.Duration(core.Config.ZKTimeout),
	)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Connected to db successfully")

	root, err := node.InitNodes(zkc)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("State initialized successfully")

	router := rest.Init(root)
	if err = http.ListenAndServe(
		fmt.Sprint(":", core.Config.RestPort),
		router,
	); err != nil {
		log.Fatal(err)
	}
	log.Println("Exit")
}
