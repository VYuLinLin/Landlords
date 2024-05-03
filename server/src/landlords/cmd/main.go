package main

import (
	"fmt"
	"landlords/internal/conf"
	_ "landlords/internal/db"
	_ "landlords/internal/game/room"
	_ "landlords/internal/router"
)

func main() {
	fmt.Println("starting serve", conf.GameConf.Version)

	fmt.Println("ending serve")
}
