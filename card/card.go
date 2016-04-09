package card

import "github.com/wmiller848/EmpireDestiny/util"

type Event struct {
}

type HiddenCard struct {
	Nonce string `BlockKey`
	Id    string `Id`
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
