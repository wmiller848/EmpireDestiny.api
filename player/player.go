package player

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/wmiller848/EmpireDestiny/card"

	rethink "github.com/dancannon/gorethink"
)

type Move struct {
}

func (m *Move) CardsToPlay() []string {
	return []string{}
}

type DeckRef struct {
	EmpireDeck   []int32 `EmpireDeck`
	DestinyDeck  []int32 `DestinyDeck`
	FortressCard int32   `FortressCard`
	GodCard      int32   `GodCard`
}

type PlayerAccount struct {
	Id    string              `Id`
	Decks map[string]*DeckRef `Decks`
	Cards []string            `Cards`
}

func CreatePlayerAccount(id string) *PlayerAccount {
	return &PlayerAccount{
		Id: id,
	}
}

func LoadPlayerAccount(id string, session *rethink.Session) (*PlayerAccount, error) {
	cursor, err := rethink.DB("ed").Table("players").Get(id).Run(session)
	if err != nil {
		return nil, err
	}
	defer cursor.Close()
	pdb := PlayerAccount{}

	if cursor.IsNil() == false {
		err = cursor.One(&pdb)
		fmt.Println(cursor, pdb, err)
		if err != nil {
			return nil, err
		} else {
			return &pdb, nil
		}
	} else {
		return nil, errors.New("No user found")
	}
}

func (p *PlayerAccount) Save(session *rethink.Session) error {
	resp, err := rethink.DB("ed").Table("players").Insert(p, rethink.InsertOpts{
		Conflict: "replace",
	}).RunWrite(session)
	if err != nil {
		return err
	}
	fmt.Println("Save", resp)
	return nil
}

func (p *PlayerAccount) LoadPlayer(deck string, session *rethink.Session) (*Player, error) {
	dck := Deck{
		EmpireDeck:  card.Deck{},
		DestinyDeck: card.Deck{},
	}

	jsn, _ := json.Marshal(p)
	fmt.Println(deck, string(jsn))
	for _, crdIndex := range p.Decks[deck].EmpireDeck {
		crd := p.Cards[crdIndex]
		cursor, err := rethink.DB("ed").Table("cards").GetAllByIndex("CardId", crd).Run(session)
		if err != nil {
			return nil, err
		}
		if cursor.IsNil() {
			return nil, errors.New("CardId not found")
		}
		crdRev := card.CardRevision{}
		err = cursor.One(&crdRev)
		if err != nil {
			return nil, err
		}
		cursor.Close()
		dck.EmpireDeck = append(dck.EmpireDeck, crdRev.Card)
	}

	for _, crdIndex := range p.Decks[deck].DestinyDeck {
		crd := p.Cards[crdIndex]
		cursor, err := rethink.DB("ed").Table("cards").GetAllByIndex("CardId", crd).Run(session)
		if err != nil {
			return nil, err
		}
		if cursor.IsNil() {
			return nil, errors.New("CardId not found")
		}
		crdRev := card.CardRevision{}
		err = cursor.One(&crdRev)
		if err != nil {
			return nil, err
		}
		cursor.Close()
		dck.DestinyDeck = append(dck.DestinyDeck, crdRev.Card)
	}

	plr := CreatePlayer(p.Id)
	plr.Deck = &dck
	return plr, nil
}

func (p *PlayerAccount) AddDeck(session *rethink.Session, deckid string, deck *DeckRef) error {
	if p.Decks == nil {
		p.Decks = make(map[string]*DeckRef)
	}
	p.Decks[deckid] = deck
	return p.Save(session)
}

type Deck struct {
	EmpireDeck   card.Deck `EmpireDeck`
	DestinyDeck  card.Deck `DestinyDeck`
	FortressCard card.Card `FortressCard`
	GodCard      card.Card `GodCard`
}

type Player struct {
	Id       string    `Id`
	DeckName string    `DeckName`
	Deck     *Deck     `Deck`
	Hand     card.Deck `Hand`
	Field    card.Deck `Field`
	Home     card.Deck `Home`
	District card.Deck `District`
	Gold     int32     `Gold`
}

func CreatePlayer(id string) *Player {
	return &Player{
		Id:       id,
		Hand:     card.Deck{},
		Field:    card.Deck{},
		Home:     card.Deck{},
		District: card.Deck{},
		Deck:     &Deck{},
	}
}

func (p *Player) Cards() (card.Deck, card.Deck) {
	return p.Deck.EmpireDeck, p.Deck.DestinyDeck
}

func (p *Player) Shuffle() {
	p.Deck.EmpireDeck.Shuffle()
	p.Deck.DestinyDeck.Shuffle()
}
