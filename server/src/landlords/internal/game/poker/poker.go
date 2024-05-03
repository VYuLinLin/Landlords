package poker

type cardType string
type cardValue int
type cardShapeType string
type cardShapeFlag int

const (
	spade   = iota + 1 // 黑桃♠
	heart              // 红桃♥
	club               // 梅花♣
	diamond            // 方片♦
	King               // 大小王
)

// CardShape 扑克牌类型
var CardShape = map[cardShapeType]cardShapeFlag{
	"S": spade,   // 黑桃♠
	"H": heart,   // 红桃♥
	"C": club,    // 梅花♣
	"D": diamond, // 方片♦
}

// Cards 扑克牌
var Cards = map[cardType]cardValue{
	"A":  12,
	"2":  13,
	"3":  1,
	"4":  2,
	"5":  3,
	"6":  4,
	"7":  5,
	"8":  6,
	"9":  7,
	"10": 8,
	"J":  9,
	"Q":  10,
	"K":  11,
}

// Kings 大小王
var Kings = map[cardType]cardValue{
	"kx": 15, // 大王
	"kd": 14, // 小王
}

// Poker information
type Poker struct {
	Shape cardShapeFlag `json:"shape"`
	Value cardValue     `json:"value"`
}

// Pokers information
type Pokers []Poker
