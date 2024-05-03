package pokerlogic

import (
	p "landlords/internal/game/poker"
	"math/rand"
	"reflect"
	"time"
)

// Card 一副新牌
type Card struct {
	NewCards  p.Pokers
	Cards     [3]p.Pokers
	HoleCards p.Pokers
}

// GetNewPokers 生成一副新牌
func (c *Card) GetNewPokers() {
	// 实例化52张牌
	for _, f := range p.CardShape {
		for _, val := range p.Cards {
			newCard := p.Poker{Shape: f, Value: val}
			c.NewCards = append(c.NewCards, newCard)
		}
	}
	// 实例化大小王
	for _, val := range p.Kings {
		newCard := p.Poker{Shape: p.King, Value: val}
		c.NewCards = append(c.NewCards, newCard)
	}
	// 洗牌
	c.shuffle()
	//	发牌
	c.dealCards()
}

// 洗牌
func (c *Card) shuffle() {
	len := len(c.NewCards)
	swap := reflect.Swapper(c.NewCards)
	rand.Seed(time.Now().Unix())
	for i := len - 1; i >= 0; i-- {
		j := rand.Intn(len)
		swap(i, j)
	}
}

// 发牌
func (c *Card) dealCards() {
	c.Cards = [3]p.Pokers{
		c.NewCards[0:17],
		c.NewCards[17:34],
		c.NewCards[34:51],
	}
	c.HoleCards = c.NewCards[51:]
}
