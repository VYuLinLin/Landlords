package person

import p "landlords/internal/poker"

// Person information
type Person struct {
	Name  string
	cards []p.Poker
}

// Person information
type Persons struct {
	user1     Person
	user2     Person
	user3     Person
	HoleCards p.Pokers
}
