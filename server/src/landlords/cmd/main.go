package main

import (
	"fmt"

	_ "landlords/internal/game/room"
	_ "landlords/internal/mysql"
	_ "landlords/internal/socketio"
	"landlords/internal/version"
)

func main() {
	fmt.Println("starting serve", version.Version)

	fmt.Println("ending serve")
}
