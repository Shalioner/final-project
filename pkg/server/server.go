package server

import (
	"fmt"
	"net/http"
	"os"

	"final-project/pkg/api"
)

func Run() {
	api.Init()

	port := os.Getenv("TODO_PORT")
	if port == "" {
		port = "7540"
	}

	fmt.Println("Server started at http://localhost:" + port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		panic(err)
	}
}
