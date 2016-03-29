package effects

import (
	"github.com/wmiller848/EmpireDestiny/card"
	"github.com/wmiller848/EmpireDestiny/player"
)

type Effect interface {
	Handle(p player.Player, c card.Card, key string, value interface{})
}

// change:
//   keys:
//     - properties
//     - attackpower
//     - armor
//     - lifeforce
//     - gold
//     - spirit

type ChangeEffect struct {
}

func (e *ChangeEffect) Handle(p player.Player, c card.Card, key string, value int32) {
	// if key == "attackpower" {
	// 	attack, armor, life := c.Prop()
	// 	c.SetProp(attack+value, armor, life)
	// } else if key == "armor" {
	// 	attack, armor, life := c.Prop()
	// 	c.SetProp(attack, armor+value, life)
	// } else if key == "lifeforce" {
	// 	attack, armor, life := c.Prop()
	// 	c.SetProp(attack, armor, life+value)
	// } else if key == "gold" {
	// 	p.SetGold(p.Gold() + value)
	// }
}

// set:
//   keys:
//     - properties
//     - attackpower
//     - armor
//     - lifeforce
//     - gold
//     - spirit

type SetEffect struct {
}

func (e *SetEffect) Handle(p player.Player, c card.Card, key string, value int32) {
	// if key == "attackpower" {
	// 	_, armor, life := c.Prop()
	// 	c.SetProp(value, armor, life)
	// } else if key == "armor" {
	// 	attack, _, life := c.Prop()
	// 	c.SetProp(attack, value, life)
	// } else if key == "lifeforce" {
	// 	attack, armor, _ := c.Prop()
	// 	c.SetProp(attack, armor, value)
	// } else if key == "gold" {
	// 	p.SetGold(value)
	// }
}

// locate_card:
//   keys:
//     - all
//     - decks
//     - empire_deck
//     - destiny_deck

type LocateCardEffect struct {
}

// func (e *LocateCardEffect) Handle(p player.Player, key string, value int32) card.Card {
// 	if key == "all" {
//
// 	} else if key == "armor" {
//
// 	} else if key == "decks" {
//
// 	} else if key == "empire_deck" {
//
// 	} else if key == "destiny_deck" {
//
// 	}
//
// 	return card.EmpireCard{}
// }

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
