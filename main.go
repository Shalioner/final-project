package main

import (
	"final-project/pkg/db"
	"final-project/pkg/server"
)

func main() {
	if err := db.Init("scheduler.db"); err != nil {
		panic(err)
	}
	defer db.DB.Close()

	server.Run()
}
