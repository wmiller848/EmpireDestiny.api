package card

import "github.com/wmiller848/EmpireDestiny/util"

type Event struct {
}

type CardRevision struct {
	Id     string `Id`
	CardId string `CardId`
	Card   *Card  `Card`
}

type CardProp map[string]interface{}

type Card struct {
	Id        string     `Id`
	Name      string     `Name`
	Type      string     `Type`
	Props     CardProp   `Props`
	Tags      []string   `Tags`
	TraitExps []TraitExp `TraitExps`
}

func CreateCard(ctype string, props CardProp, tags []string, traitExps []TraitExp, name string) *Card {
	card := &Card{
		Id:        util.Hash(name),
		Name:      name,
		Type:      ctype,
		Tags:      tags,
		TraitExps: traitExps,
		Props:     props,
	}
	return card
}

// type CardInstance interface {
// 	Card
// 	Engage()
// 	Disengage()
// 	Engaged() bool //
//
// 	Bow()
// 	Unbow()
// 	Bowed() bool
//
// 	Kill()
// 	Revive()
// 	Dead() bool
//
// 	Lock()
// 	Unlock()
// 	Locked() bool
// }

// func (b *BaseCard) Engage() {
// 	if b.locked == false {
// 		b.engaged = true
// 		b.Bow()
// 	}
// }
// func (b *BaseCard) Disengage() {
// 	if b.locked == false {
// 		b.engaged = false
// 		b.Unbow()
// 	}
// }
// func (b *BaseCard) Engaged() bool {
// 	return b.engaged
// }
//
// func (b *BaseCard) Bow() {
// 	if b.locked == false {
// 		b.bowed = true
// 	}
// }
// func (b *BaseCard) Unbow() {
// 	if b.locked == false {
// 		b.bowed = false
// 	}
// }
// func (b *BaseCard) Bowed() bool {
// 	return b.bowed
// }
//
// func (b *BaseCard) Kill() {
// 	if b.locked == false {
// 		b.dead = true
// 	}
// }
// func (b *BaseCard) Revive() {
// 	if b.locked == false {
// 		b.dead = false
// 	}
// }
// func (b *BaseCard) Dead() bool {
// 	return b.dead
// }
//
// func (b *BaseCard) Lock() {
// 	b.locked = true
// }
// func (b *BaseCard) Unlock() {
// 	b.locked = false
// }
// func (b *BaseCard) Locked() bool {
// 	return b.locked
// }
