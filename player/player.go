package player

import (
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
	EmpireDeck   []int32
	DestinyDeck  []int32
	FortressCard int32
	GodCard      int32
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
	// cursor, err := rethink.DB("ed").Table("players").Get(p.Id).Run(session)
	// if err != nil {
	// 	return nil, err
	// }
	// defer cursor.Close()
	// suc := cursor.One(&pdb)
	// fmt.Println("LoadPlayer", suc, pdb)
	// dck := card.GetDeckFromIds(plr.Deck(deckid))
	// dck := card.GetDeckFromIds(plr.Deck(deckid))
	// p.deckName = dck.Name()
	// p.deck = &Deck{}
	dck := Deck{
		EmpireDeck:  card.Deck{},
		DestinyDeck: card.Deck{},
	}
	for _, crdIndex := range p.Decks[deck].EmpireDeck {
		// crdInstance := card.CreateEmpireCard(crd, armor int32, life int32, tags []string, traitExps []card.TraitExp, name string)
		// dck.EmpireDeck = append(dck.EmpireDeck, crdInstance)
		crd := p.Cards[crdIndex]
		// fmt.Println("LoadPlayer", crd)
		// index := make(map[string]interface{})
		// index["index"] = "CardId"
		cursor, err := rethink.DB("ed").Table("cards").GetAllByIndex("CardId", crd).Run(session)
		if err != nil {
			return nil, err
		}
		if cursor.IsNil() {
			return nil, errors.New("CardId not found")
		}
		crdRev := card.CardRevision{}
		// crdInstance := make(map[string]interface{})
		err = cursor.One(&crdRev)
		if err != nil {
			return nil, err
		}
		cursor.Close()
		// 	suc := res.Next(&crdInstance)
		// fmt.Println("LoadPlayer", crdInstance)
		dck.EmpireDeck = append(dck.EmpireDeck, crdRev.Card)
	}
	// for _, crdIndex := range pdb.Decks[deck].DestinyDeck {
	// 	crd := pdb.Cards[crdIndex]
	// 	fmt.Println(crd)
	// }
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
	District card.Deck `District`
	Gold     int32     `Gold`
}

func CreatePlayer(id string) *Player {
	return &Player{
		Id:       id,
		Hand:     card.Deck{},
		Field:    card.Deck{},
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
