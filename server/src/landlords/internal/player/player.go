package player

import p "landlords/internal/poker"

// Player information
type Player struct {
	Name  string
	cards []p.Poker
}

// Players information
type Players struct {
	user1     Player
	user2     Player
	user3     Player
	HoleCards p.Pokers
}
