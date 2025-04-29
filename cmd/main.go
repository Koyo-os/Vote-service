package main

import (
	"github.com/Koyo-os/Vote-service/pkg/config"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		return
	}
}
