package main

import (
	"github.com/templecloud/memoria-server/internal/memoria/boot"
)

func main() {
	// Start the Memoria API server.
	boot.Start()
}