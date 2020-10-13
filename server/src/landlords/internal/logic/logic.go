package logic

import (
	"fmt"
	"landlords/internal/logic/pokerlogic"
	"landlords/internal/person"
)

// InitRoom 初始化房间
func InitRoom() {
	card := &pokerlogic.Card{}
	card.NewPokers()
	fmt.Println(card.Cards)
	fmt.Println(len(card.Cards))

	persons := &person.Persons{}
}
