package queue

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"sync"

	"github.com/wmiller848/EmpireDestiny/player"
)

type Session struct {
	sync.Mutex
	PlayerA   *player.Player
	PlayerB   *player.Player
	Sessionid string
}

type PlayerQueue struct {
	player  *player.Player
	channel chan *Session
}

type Queue struct {
	next    *Queue
	players [2]*PlayerQueue
}

func CreateQueue() *Queue {
	return &Queue{
		players: [2]*PlayerQueue{
			&PlayerQueue{
				channel: make(chan *Session),
			},
			&PlayerQueue{
				channel: make(chan *Session),
			},
		},
	}
}

func (q *Queue) AddToQueue(player *player.Player) {
	if q.next != nil {
		q.next.AddToQueue(player)
	} else {
		playerA := q.players[0]
		playerB := q.players[1]
		if playerA.player != nil {
			if playerA.player.Id == player.Id {
				return
			}
		} else if playerB.player != nil {
			if playerB.player.Id == player.Id {
				return
			}
		}

		if playerA.player == nil {
			playerA.player = player
		} else if playerB.player == nil {
			playerB.player = player
		} else {
			q.next = CreateQueue()
			q.next.AddToQueue(player)
		}
		// q.players[0] = playerA
		// q.players[1] = playerB
	}
}

func (q *Queue) Listen(playerid string) (*Queue, chan *Session) {
	playerA := q.players[0]
	playerB := q.players[1]

	if playerA.player != nil {
		if playerA.player.Id == playerid {
			go q.Consume()
			return q, playerA.channel
		}
	}

	if playerB.player != nil {
		if playerB.player.Id == playerid {
			go q.Consume()
			return q, playerB.channel
		}
	}

	if q.next != nil {
		return q.next.Listen(playerid)
	}

	return q, nil
}

func (q *Queue) Consume() {
	playerA := q.players[0]
	playerB := q.players[1]

	if playerA.player == nil || playerB.player == nil {
		return
	}

	buf := make([]byte, 32)
	_, err := rand.Read(buf)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	sessionid := hex.EncodeToString(buf)
	session := Session{
		PlayerA:   playerA.player,
		PlayerB:   playerB.player,
		Sessionid: sessionid,
	}
	playerA.channel <- &session
	playerB.channel <- &session
	// //
	// close(playerA.channel)
	// close(playerB.channel)
}

// func (q *Queue) Consume(playerQueue *map[string]chan *Session) *Queue {
// 	pq := *playerQueue
// 	for {
// 		if q == nil {
// 			return nil
// 		}
// 		playerA := q.players[0]
// 		playerB := q.players[1]
// 		if playerA != nil && playerB != nil {
// 			go func() {
// 				for {
// 					time.Sleep(1000 * time.Millisecond)
// 					if pq == nil || playerA == nil || playerB == nil {
// 						return
// 					}
// 					if pq[playerA.Id] != nil && pq[playerB.Id] != nil {
// 						playerChanA := pq[playerA.Id]
// 						playerChanB := pq[playerB.Id]
// 						buf := make([]byte, 32)
// 						_, err := rand.Read(buf)
// 						if err != nil {
// 							fmt.Println(err.Error())
// 							return
// 						}
// 						sessionid := hex.EncodeToString(buf)
// 						session := Session{
// 							PlayerA:   playerA,
// 							PlayerB:   playerB,
// 							Sessionid: sessionid,
// 							Mutex:     &sync.Mutex{},
// 						}
// 						playerChanB <- &session
// 						playerChanA <- &session
// 						q.next.Consume(playerQueue)
// 						delete(pq, playerA.Id)
// 						delete(pq, playerB.Id)
// 					}
// 				}
// 			}()
// 			return q.next
// 		}
// 	}
// }
