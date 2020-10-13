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
	// King 大小王
	King
)

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

// CardShape 扑克牌类型
var CardShape = map[cardShapeType]cardShapeFlag{
	"S": spade,
	"H": heart,
	"C": club,
	"D": diamond,
}

// Kings 大小王
var Kings = map[cardType]cardValue{
	"kx": 14,
	"kd": 15,
}

// Poker information
type Poker struct {
	Shape cardShapeFlag
	Value cardValue
}

// Pokers information
type Pokers []Poker
