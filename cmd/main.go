package main

import (
	"github.com/jlauser/gocache/internal/config"
	"log"
)

func panicHandler() {
	if r := recover(); r != nil {
		log.Println("recovered from error:", r)
	}
}

func main() {
	defer panicHandler()

	env, err := config.LoadConfig()
	if err != nil {
		panic(err)
	}
	print(env.Mode)
}
