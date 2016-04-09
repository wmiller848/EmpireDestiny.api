package match

import (
	"errors"
	"fmt"

	"github.com/wmiller848/EmpireDestiny/card"
	"github.com/wmiller848/EmpireDestiny/player"
	"github.com/wmiller848/EmpireDestiny/util"

	rethink "github.com/dancannon/gorethink"
)

const stageMatchPhase uint32 = 0
const harvestPhase uint32 = 1
const eventPhase uint32 = 2
const eventResponsePhase uint32 = 3
const conquestPhase uint32 = 4
const buildingPhase uint32 = 5

type Match struct {
	Round     uint32                           `Round`
	Phase     uint32                           `Phase`
	PlayerA   *player.Player                   `PlayerA`
	PlayerB   *player.Player                   `PlayerB`
	PlayerAID string                           `PlayerAID`
	PlayerBID string                           `PlayerBID`
	Events    map[string]map[string]*card.Card `Events`
	Id        string                           `Id`
	Key       string                           `Key`
}

func CreateMatch(PlayerAID, PlayerBID, id string) *Match {
	key := util.Random(16)
	match := &Match{
		Round:     0,
		Phase:     stageMatchPhase,
		PlayerAID: PlayerAID,
		PlayerBID: PlayerBID,
		Id:        id,
		Key:       util.Hex(key),
	}
	return match
}

func (m *Match) LoadPlayers(session *rethink.Session) error {
	pdbA, err := player.LoadPlayerAccount(m.PlayerAID, session)
	m.PlayerA, err = pdbA.LoadPlayer("Default", session)
	if err != nil {
		return err
	}

	pdbB, err := player.LoadPlayerAccount(m.PlayerBID, session)
	m.PlayerB, err = pdbB.LoadPlayer("Default", session)
	if err != nil {
		return err
	}
	return nil
}

func (m *Match) PlayRound(playerAMove, playerBMove *player.Move) error {
	var attackerMove *player.Move
	var attackingPlayer *player.Player

	// var defenderMove *player.Move
	var defendingPlayer *player.Player
	if m.Round == stageMatchPhase {
		m.prepareMatchPhase()
	}

	if m.Round%2 == 0 {
		attackerMove = playerAMove
		// defenderMove = playerBMove
		attackingPlayer = m.PlayerA
		defendingPlayer = m.PlayerB
	} else {
		attackerMove = playerBMove
		// defenderMove = playerAMove
		attackingPlayer = m.PlayerB
		defendingPlayer = m.PlayerA
	}

	if attackingPlayer == nil || attackerMove == nil || defendingPlayer == nil {
		return errors.New("Players Deck and/or Attacker Move nil")
	}

	// Phases
	if m.Phase == harvestPhase {
		m.harvestPhase(attackingPlayer)
	} else if m.Phase == eventPhase {
		m.eventPhase(attackerMove, attackingPlayer, defendingPlayer)
	} else if m.Phase == eventResponsePhase {
		m.eventPhase(attackerMove, attackingPlayer, defendingPlayer)
	} else if m.Phase == conquestPhase {
		m.conquestPhase(attackerMove, attackingPlayer, defendingPlayer)
	} else if m.Phase == buildingPhase {
		m.buildingPhase(attackerMove, attackingPlayer, defendingPlayer)
		m.Round++
	}

	m.Phase++
	if m.Phase > buildingPhase {
		m.Phase = harvestPhase
	}
	return nil
}

func (m *Match) prepareMatchPhase() {
	//
	m.PlayerA.Shuffle()
	m.PlayerB.Shuffle()
}

func (m *Match) endMatchPhase() {
	// reward points
}

//
//
func (m *Match) harvestPhase(p *player.Player) {
	// empireDeck, destinyDeck := p.Cards()
	field := p.Field
	for _, c := range field {
		// c.Unbow()
		fmt.Println(c)
	}
}

func (m *Match) eventPhase(move *player.Move, attackingPlayer, defendingPlayer *player.Player) {
	// Reset the event map
	m.Events = make(map[string]map[string]*card.Card)

	// Populate the event map
	defenderField := defendingPlayer.Field
	for _, c := range defenderField {
		for _, tag := range c.Tags {
			m.Events[tag][c.Id] = c
		}
	}

	// attackerHand := attackingPlayer.Hand()
	// for _, c := range move.CardsToPlay() {
	// 	cardToPlay := attackerHand[c]
	// 	if cardToPlay != nil {
	// 		traits := cardToPlay.TraitExps()
	// 		for _, trait := range traits {
	// 			for _, tag := range trait.Tags() {
	// 				events := m.events[tag]
	// 				for _, cardToEffect := range events {
	// 					fmt.Println(cardToPlay, cardToEffect)
	// 				}
	// 			}
	// 		}
	// 	} else {
	//
	// 	}
	// }
}

func (m *Match) conquestPhase(move *player.Move, p *player.Player, defendingPlayer *player.Player) {

}

func (m *Match) buildingPhase(move *player.Move, p *player.Player, defendingPlayer *player.Player) {
}
