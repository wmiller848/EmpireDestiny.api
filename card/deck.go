package card

import (
	"crypto/aes"
	"crypto/cipher"
	"fmt"
	"io/ioutil"

	"github.com/wmiller848/EmpireDestiny/util"

	"gopkg.in/yaml.v2"
)

type HiddenDeck []*HiddenCard

type Deck []*Card

func (d Deck) Shuffle() {

}

func (d Deck) ToHiddenDeck(key string) HiddenDeck {
	cypher, err := aes.NewCipher([]byte(key))
	if err != nil {
		fmt.Println(err.Error())
	}
	hiddenEmpire := make(HiddenDeck, len(d))
	for i, crd := range d {
		nonce := util.Random(aes.BlockSize)
		encrypter := cipher.NewCFBEncrypter(cypher, nonce)
		crdId := []byte(crd.Id)
		encrypted := make([]byte, len(crdId))
		encrypter.XORKeyStream(encrypted, crdId)
		hiddenEmpire[i] = &HiddenCard{
			Nonce: util.Hex(nonce),
			Id:    util.Hex(encrypted),
		}
	}
	return hiddenEmpire
}

func (d Deck) Next() *Card {
	l := len(d)
	if l > 0 {
		return d[l-1]
	} else {
		return nil
	}
}

type TraitExp struct {
	Exp         string `Exp`
	Targets     string `Targets`
	Name        string `Name`
	Description string `Description`
}

type EmpireCardYML struct {
	Name        string     `Name`
	AttackPower int32      `AttackPower`
	Armor       int32      `Armor`
	Lifeforce   int32      `Lifeforce`
	Tags        []string   `Tags`
	TraitExps   []TraitExp `TraitExps`
	Props       CardProp   `Props`
}

func LoadEmpireDeckFromYML(path, name string) (Deck, error) {
	deck := []EmpireCardYML{}
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	err = yaml.Unmarshal([]byte(data), &deck)
	if err != nil {
		return nil, err
	}

	empireDeck := Deck{}

	for _, c := range deck {
		// cProp := []CardProp{}
		// for _, prop := range c.Props {
		// 	cProp = append(cProp, CardProp{
		// 		Name:  prop.Name,
		// 		Value: prop.Value,
		// 	})
		// }
		empireCard := CreateCard("empire", c.Props, c.Tags, c.TraitExps, c.Name)
		empireDeck = append(empireDeck, empireCard)
	}

	return empireDeck, nil
}

type DestinyCardYML struct {
	Name      string     `Name`
	Tags      []string   `Tags`
	TraitExps []TraitExp `TraitExps`
	Props     CardProp   `Props`
}

func LoadDestinyDeckFromYML(path, name string) (Deck, error) {
	deck := []DestinyCardYML{}
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	err = yaml.Unmarshal([]byte(data), &deck)
	if err != nil {
		return nil, err
	}

	destinyDeck := Deck{}

	for _, c := range deck {
		destinyCard := CreateCard("destiny", c.Props, c.Tags, c.TraitExps, c.Name)
		destinyDeck = append(destinyDeck, destinyCard)
	}

	return destinyDeck, nil
}

func LoadFortressCardsFromYML(path string) (Deck, error) {
	deck := []EmpireCardYML{}
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	err = yaml.Unmarshal([]byte(data), &deck)
	if err != nil {
		return nil, err
	}

	fortressCards := Deck{}

	for _, c := range deck {
		fortressCard := CreateCard("fortress", c.Props, c.Tags, c.TraitExps, c.Name)
		fortressCards = append(fortressCards, fortressCard)
	}

	return fortressCards, nil
}

func LoadGodCardsFromYML(path string) (Deck, error) {
	deck := []EmpireCardYML{}
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	err = yaml.Unmarshal([]byte(data), &deck)
	if err != nil {
		return nil, err
	}

	godCards := Deck{}

	for _, c := range deck {
		godCard := CreateCard("god", c.Props, c.Tags, c.TraitExps, c.Name)
		godCards = append(godCards, godCard)
	}

	return godCards, nil
}
