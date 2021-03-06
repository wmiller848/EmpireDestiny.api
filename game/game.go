package game

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/wmiller848/EmpireDestiny/card"
	"github.com/wmiller848/EmpireDestiny/game/effects"
	"github.com/wmiller848/EmpireDestiny/game/traits"
	"github.com/wmiller848/EmpireDestiny/util"

	"gopkg.in/yaml.v2"
)

type GameYML struct {
	Fortresses struct {
		Path string `Path`
	} `Fortresses`
	Gods struct {
		Path string `Path`
	} `Gods`
	Decks struct {
		Path string `Path`
	} `Decks`
	Effects struct {
		Path string `Path`
	} `Effects`
	Traits struct {
		Path string `Path`
	} `Traits`
}

type TwinDeck struct {
	Name    string    `Name`
	Empire  card.Deck `Empire`
	Destiny card.Deck `Destiny`
}

type Game struct {
	Traits     []traits.Trait
	Effects    []effects.Effect
	TwinDecks  []TwinDeck
	Fortresses card.Deck
	Gods       card.Deck
	Index      map[string]*card.Card
}

//func (g *Game) MarshalJSON() ([]byte, error) {
//jsonData := make(map[string]interface{})
//gameData := make(map[string]interface{})
//gameData["TwinDecks"] = g.TwinDecks()
//gameData["Fortresses"] = g.Fortresses()
//gameData["Gods"] = g.Gods()
//jsonData["Data"] = gameData
//jsonData["Index"] = g.Index()
//return json.Marshal(jsonData)
//}

func (g *Game) sync(crd *card.Card) card.CardRevision {
	jsn, _ := json.Marshal(crd)
	return card.CardRevision{
		Id:     util.Hash(string(jsn)),
		CardId: crd.Id,
		Card:   crd,
	}
}

func (g *Game) Sync() chan card.CardRevision {
	chn := make(chan card.CardRevision)
	go func() {
		for _, crd := range g.Fortresses {
			chn <- g.sync(crd)
		}
		for _, crd := range g.Gods {
			chn <- g.sync(crd)
		}
		for _, dck := range g.TwinDecks {
			for _, crd := range dck.Empire {
				chn <- g.sync(crd)
			}
			for _, crd := range dck.Destiny {
				chn <- g.sync(crd)
			}
		}
		close(chn)
	}()
	return chn
}

func LoadGameFromYML(path string) (*Game, error) {
	game := GameYML{}
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	err = yaml.Unmarshal([]byte(data), &game)
	if err != nil {
		return nil, err
	}

	file, err := os.Open(game.Decks.Path)
	if err != nil {
		return nil, err
	}

	folders, err := file.Readdirnames(0)
	if err != nil {
		return nil, err
	}

	g := &Game{
		TwinDecks: []TwinDeck{},
		Index:     make(map[string]*card.Card),
	}
	for _, folder := range folders {
		root := game.Decks.Path + "/" + folder
		empire, err := card.LoadEmpireDeckFromYML(root+"/empire.yml", folder)
		if err != nil {
			fmt.Println(err.Error())
		}
		destiny, err := card.LoadDestinyDeckFromYML(root+"/destiny.yml", folder)
		if err != nil {
			fmt.Println(err.Error())
		}
		fmt.Println(empire, destiny)
		g.TwinDecks = append(g.TwinDecks, TwinDeck{Name: folder, Empire: empire, Destiny: destiny})
	}

	fortresses, err := card.LoadFortressCardsFromYML(game.Fortresses.Path)
	if err != nil {
		fmt.Println(err.Error())
	} else {
		g.Fortresses = fortresses
	}

	gods, err := card.LoadGodCardsFromYML(game.Gods.Path)
	if err != nil {
		fmt.Println(err.Error())
	} else {
		g.Gods = gods
	}

	for _, crd := range g.Fortresses {
		g.Index[crd.Id] = crd
	}
	for _, crd := range g.Gods {
		g.Index[crd.Id] = crd
	}
	for _, dck := range g.TwinDecks {
		for _, crd := range dck.Empire {
			g.Index[crd.Id] = crd
		}
		for _, crd := range dck.Destiny {
			g.Index[crd.Id] = crd
		}
	}

	return g, nil
}
