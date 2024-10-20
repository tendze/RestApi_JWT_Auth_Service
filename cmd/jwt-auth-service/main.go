package main

import (
	"fmt"
	"jwt-auth-service/internal/config"
)

func main() {
	cfg := config.MustLoad()
	fmt.Println(cfg)
	// TODO: LOAD CONFIG
	// TODO: INIT LOGGER
	// TODO: INIT STORAGE: POSTGRESQL
	// TODO: INIT ROUTER CHI:"chi render"
	// TODO: RUN SERVER
}
