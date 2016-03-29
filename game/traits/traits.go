package traits

import (
	"github.com/wmiller848/EmpireDestiny/player"
)

type Trait interface {
	Handle(string)
}

// gold:
//   keys:
//     - increase
//     - decrease

type GoldTrait struct {
}

func (g *GoldTrait) Increase(p player.Player, amount int32) {
	// p.SetGold(p.Gold() + amount)
}

func (g *GoldTrait) Decrease(p player.Player, amount int32) {
	// p.SetGold(p.Gold() + -amount)
}

// spirit:
//   keys:
//     - increase
//     - decrease

// change:
//   keys:
//     - properties
//     - attackpower
//     - armor
//     - lifeforce

// locate_card:
//   keys:
//     - all
//     - decks
//     - empire_deck
//     - destiny_deck

// play:
//   keys:
//     - inplay
//     - field
//     - districts

// hand:
//   keys:
//     - enemy
//     - ally

// create:
//   keys:
//     - enemy
//     - ally
