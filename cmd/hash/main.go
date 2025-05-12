package main

import (
	"fmt"
	"os"

	"github.com/budhilaw/personal-website-backend/pkg/util"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run cmd/hash/main.go <password>")
		os.Exit(1)
	}

	password := os.Args[1]
	hash, err := util.HashPassword(password)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Argon2id hash for password '%s':\n%s\n", password, hash)
}
