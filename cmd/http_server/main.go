package main

import (
	"gigaAPI/config"
	"gigaAPI/internal/http/http_server"
	"log"
)

func main() {

	conf := config.InitConfig("../../postgres.yaml")

	//TODO: logger

	//TODO: server -> db

	server := http_server.InitServer(conf)

	log.Fatal(server.Run())

	/*server := http_server.NewServerHTTPS(":1234", "../../localhost.crt", "../../localhost.key", "../../.env")

	log.Fatal(server.RunServer())*/
}
