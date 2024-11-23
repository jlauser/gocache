package main

import "github.com/jlauser/gocache/internal/config"

func main() {
	env, err := config.LoadConfig()
	if err != nil {
		panic(err)
	}
	print(env.Mode)
}
