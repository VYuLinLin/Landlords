package pokerlogic

import (
	"math/rand"
	"reflect"
	"time"

	"landlords/internal/player"
	p "landlords/internal/poker"
)

// Card 一副新牌
type Card struct {
	Cards     p.Pokers
	User1     player.Player
	User2     player.Player
	User3     player.Player
	HoleCards p.Pokers
}

// NewPokers 生成一副新牌
func (c *Card) NewPokers() {
	// 实例化52张牌
	for _, f := range p.CardShape {
		for _, val := range p.Cards {
			newCard := p.Poker{Shape: f, Value: val}
			c.Cards = append(c.Cards, newCard)
		}
	}
	// 实例化大小王
	for _, val := range p.Kings {
		newCard := p.Poker{Shape: p.King, Value: val}
		c.Cards = append(c.Cards, newCard)
	}
	// 洗牌
	c.shuffle()
}

// 洗牌
func (c *Card) shuffle() {
	len := len(c.Cards)
	swap := reflect.Swapper(c.Cards)
	rand.Seed(time.Now().Unix())
	for i := len - 1; i >= 0; i-- {
		j := rand.Intn(len)
		swap(i, j)
	}
}

// 发牌
func dealCards() {
	// pokers := NewPokers()

}
