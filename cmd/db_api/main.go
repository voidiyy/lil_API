package main

import (
	"fmt"
	"gigaAPI/internal/db_psql"
	"log"
)

func main() {
	psql, err := db_psql.NewPSQL(".env")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("succ")
	psql.Run()
}
