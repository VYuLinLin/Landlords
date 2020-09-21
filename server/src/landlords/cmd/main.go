package main

import (
	"fmt"

	_ "landlords/internal/socketio"
	"landlords/internal/version"
)

func main() {
	fmt.Println("starting serve", version.Version)

	fmt.Println("ending serve")
}
