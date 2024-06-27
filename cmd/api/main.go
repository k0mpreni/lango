package main

import (
	"fmt"
	"lango/internal/database"
	"lango/internal/server"
)

func main() {
	if err := database.Init(); err != nil {
		fmt.Println("error init database", err)
	}
	defer database.Close()

	server := server.NewServer()

	err := server.ListenAndServe()
	if err != nil {
		panic(fmt.Sprintf("cannot start server: %s", err))
	}
}
