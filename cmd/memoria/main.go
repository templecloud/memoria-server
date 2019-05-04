package main

import (
	"github.com/templecloud/memoria-server/internal/memoria/boot"
)

func main() {
	// Start the memoria api server.
	boot.Start()
}