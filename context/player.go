package context

import (
	"crypto/rand"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"net/http"

	"github.com/wmiller848/EmpireDestiny/player"

	rethink "github.com/dancannon/gorethink"
	"github.com/gocraft/web"
)

func (c *Context) CreatePlayer(id string) (*player.Player, error) {
	pdb := player.CreatePlayerAccount(id)
	// Give random cards
	defDeck := player.DeckRef{
		EmpireDeck:   []int32{},
		DestinyDeck:  []int32{},
		FortressCard: -1,
		GodCard:      -1,
	}
	defCards := make([]string, 104)

	cardsGiven := 0
	for cardsGiven < 2 {
		max := big.NewInt(10000)
		num, err := rand.Int(rand.Reader, max)
		if err != nil {
			return nil, err
		}
		fortresses := c.Game.Fortresses()
		forLen := uint64(len(fortresses))
		pick := num.Uint64() % forLen
		crd := fortresses[pick]
		fmt.Println(forLen, pick, crd)
		defCards[cardsGiven] = crd.Id
		if defDeck.FortressCard == -1 {
			defDeck.FortressCard = int32(cardsGiven)
		}
		cardsGiven++
	}

	next := cardsGiven + 2
	for cardsGiven < next {
		max := big.NewInt(10000)
		num, err := rand.Int(rand.Reader, max)
		if err != nil {
			return nil, err
		}
		gods := c.Game.Gods()
		godLen := uint64(len(gods))
		pick := num.Uint64() % godLen
		crd := gods[pick]
		defCards[cardsGiven] = crd.Id
		if defDeck.GodCard == -1 {
			defDeck.GodCard = int32(cardsGiven)
		}
		cardsGiven++
	}
	next = cardsGiven + 50
	for cardsGiven < next {
		for _, dck := range c.Game.TwinDecks() {
			max := big.NewInt(10000)
			num, err := rand.Int(rand.Reader, max)
			if err != nil {
				return nil, err
			}
			empLen := uint64(len(dck.Empire))
			pick := num.Uint64() % empLen
			crd := dck.Empire[pick]
			// TODO GET REVISION ID INSTEAD
			defCards[cardsGiven] = crd.Id
			defDeck.EmpireDeck = append(defDeck.EmpireDeck, int32(cardsGiven))
			cardsGiven++
		}
	}

	next = cardsGiven + 50
	for cardsGiven < next {
		for _, dck := range c.Game.TwinDecks() {
			max := big.NewInt(10000)
			num, err := rand.Int(rand.Reader, max)
			if err != nil {
				return nil, err
			}
			desLen := uint64(len(dck.Destiny))
			pick := num.Uint64() % desLen
			crd := dck.Destiny[pick]
			defCards[cardsGiven] = crd.Id
			defDeck.DestinyDeck = append(defDeck.DestinyDeck, int32(cardsGiven))
			cardsGiven++
		}
	}

	pdb.Cards = append(pdb.Cards, defCards...)
	pdb.Decks = make(map[string]*player.DeckRef)
	pdb.Decks["Default"] = &defDeck

	pdb.Save(c.Session)
	return c.GetPlayerInstance(id)
}

func (c *Context) GetPlayerInstance(id string) (*player.Player, error) {
	if id == "" {
		return nil, errors.New("Invalid PlayerID")
	}
	pdb, err := player.LoadPlayerAccount(id, c.Session)
	if err != nil {
		return nil, err
	}
	fmt.Println("LoadPlayerDB", pdb)
	plrInstance, err := pdb.LoadPlayer("Default", c.Session)
	if err != nil {
		return nil, err
	}
	return plrInstance, nil
}

func (c *PlayerContext) GetPlayer(rw web.ResponseWriter, req *web.Request) {
	defer c.Session.Close(rethink.CloseOpts{})
	// fmt.Fprint(rw, "/")
	fmt.Println("GetPlayer")
	pdb, err := player.LoadPlayerAccount(c.playerid, c.Session)
	if err != nil {
		fmt.Println(err.Error())
		rw.WriteHeader(http.StatusNotFound)
		return
	}

	fmt.Println("LoadPlayerDB", pdb)
	jsn, err := json.Marshal(pdb)
	fmt.Println(string(jsn), err)
	if err != nil {
		fmt.Println(err.Error())
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}
	fmt.Fprint(rw, string(jsn))
	rw.WriteHeader(http.StatusOK)
}

func (c *PlayerContext) PostDecks(rw web.ResponseWriter, req *web.Request) {
	// fmt.Fprint(rw, "/")
}

func (c *PlayerContext) PutDecks(rw web.ResponseWriter, req *web.Request) {
	// fmt.Fprint(rw, "/")
}

func (c *PlayerContext) GetGame(rw web.ResponseWriter, req *web.Request) {
	defer c.Session.Close(rethink.CloseOpts{})
	fmt.Println("GetGame")
	jsn, err := json.Marshal(c.Game)
	if err != nil {
		fmt.Println(err.Error())
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}
	fmt.Fprint(rw, string(jsn))
	// rw.WriteHeader(http.StatusOK)
	fmt.Println("END GetGame")
}
