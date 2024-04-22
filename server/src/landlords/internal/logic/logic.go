package logic

import (
	"fmt"
	"landlords/internal/logic/pokerlogic"
	"landlords/internal/player"
)

// InitRoom 初始化房间
func InitRoom() {
	card := &pokerlogic.Card{}
	card.NewPokers()
	fmt.Println(card.Cards)
	fmt.Println(len(card.Cards))

	players := &player.Players{}
}
