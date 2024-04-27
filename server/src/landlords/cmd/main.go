package main

import (
	"fmt"
	_ "landlords/internal/conf"
	_ "landlords/internal/db"
	_ "landlords/internal/game/room"
	_ "landlords/internal/router"
	"landlords/internal/version"
)

func main() {
	fmt.Println("starting serve", version.Version)

	fmt.Println("ending serve")
}
