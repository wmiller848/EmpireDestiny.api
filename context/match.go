package context

import (
	"encoding/json"
	"fmt"
	"net/http"
	"runtime"

	"github.com/wmiller848/EmpireDestiny/card"
	"github.com/wmiller848/EmpireDestiny/match"
	"github.com/wmiller848/EmpireDestiny/player"

	rethink "github.com/dancannon/gorethink"
	"github.com/gocraft/web"
)

type EnemyRsp struct {
	Field    card.Deck
	Home     card.Deck
	Destrict card.Deck
	God      card.Card
	Fortress card.Card
}

type PlayerRsp struct {
	Empire   card.HiddenDeck
	Destiny  card.HiddenDeck
	Field    card.Deck
	Home     card.Deck
	Destrict card.Deck
	Hand     card.Deck
	God      card.Card
	Fortress card.Card
	Gold     int32
}

func (c *Context) SetMatchInstance(mtch *match.Match) error {
	resp, err := rethink.DB("ed").Table("matches").Insert(mtch).Run(c.Session)
	if err != nil {
		return err
	}
	fmt.Println("SetMatchInstance", resp)
	return nil
}

func (c *Context) GetMatchInstance(id string) (*match.Match, error) {
	cursor, err := rethink.DB("ed").Table("matches").Get(id).Run(c.Session)
	if err != nil {
		return nil, err
	}
	defer cursor.Close()
	matchInstance := match.Match{}
	err = cursor.One(&matchInstance)
	if err != nil {
		return nil, err
	}
	return &matchInstance, nil
}

func (c *MatchContext) GetMatch(rw web.ResponseWriter, req *web.Request) {
	defer c.Session.Close(rethink.CloseOpts{})
	_, err := c.GetMatchInstance(c.sessionid)
	if err == nil {
		rw.Header().Set("Location", "/player/match/info")
		rw.WriteHeader(http.StatusFound)
		return
	}
	fmt.Println("GetMatch", err.Error())
	// deckid := req.PathParams["deckid"]
	plr, err := player.LoadPlayerAccount(c.playerid, c.Session)
	if err != nil {
		fmt.Println("GetMatch", err.Error())
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	deckid := req.PathParams["deckid"]
	fmt.Println("GetMatch", deckid)
	dck := plr.Decks[deckid]
	if dck == nil {
		rw.WriteHeader(http.StatusNotFound)
		return
	}

	plrInstance, err := plr.LoadPlayer(deckid, c.Session)
	if err != nil {
		fmt.Println("GetMatch - ", deckid, err.Error())
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	c.Queue.AddToQueue(plrInstance)
	q, pchan := c.Queue.Listen(c.playerid)
	fmt.Println(q, pchan)
	if pchan == nil {
		rw.WriteHeader(http.StatusNotFound)
		return
	}
	session, open := <-pchan
	fmt.Println(session)
	// match has started
	if session != nil && open == true {
		session.Lock()
		_, err = c.GetMatchInstance(c.sessionid)
		if err != nil {
			mtch := match.CreateMatch(session.PlayerA.Id, session.PlayerB.Id, session.Sessionid)
			mtch.LoadPlayers(c.Session)
			c.SetMatchInstance(mtch)
		}
		session.Unlock()
		runtime.Gosched()
		cookie := &http.Cookie{
			Name:  "EDSESSION",
			Value: session.Sessionid,
			Path:  "/",
		}
		http.SetCookie(rw, cookie)
	}
	rw.WriteHeader(http.StatusCreated)
}

func (c *MatchContext) GetInfo(rw web.ResponseWriter, req *web.Request) {
	mtch, err := c.GetMatchInstance(c.sessionid)
	fmt.Println(mtch, err)
	if err != nil {
		fmt.Println(err.Error())
		rw.WriteHeader(http.StatusNotFound)
		return
	}

	//err = mtch.LoadPlayers(c.Session)
	//if err != nil {
	//fmt.Println(err.Error())
	//rw.WriteHeader(http.StatusNotFound)
	//return
	//}

	var plr *player.Player
	var enemy *player.Player
	if mtch.PlayerAID == c.playerid {
		plr = mtch.PlayerA
		enemy = mtch.PlayerB
	} else {
		plr = mtch.PlayerB
		enemy = mtch.PlayerA
	}

	fmt.Println(plr, enemy)

	rsp := make(map[string]interface{})
	rsp["Session"] = c.sessionid
	rsp["Enemy"] = EnemyRsp{
		Field:    enemy.Field,
		Home:     enemy.Home,
		Destrict: enemy.District,
		God:      enemy.Deck.GodCard,
		Fortress: enemy.Deck.FortressCard,
	}
	empire, destiny := plr.Cards()
	rsp["Player"] = PlayerRsp{
		Field:    plr.Field,
		Home:     plr.Home,
		Destrict: plr.District,
		God:      plr.Deck.GodCard,
		Fortress: plr.Deck.FortressCard,
		Empire:   empire.ToHiddenDeck(mtch.Key),
		Destiny:  destiny.ToHiddenDeck(mtch.Key),
		Hand:     plr.Hand,
	}
	fmt.Println(rsp)
	d, err := json.Marshal(rsp)
	if err != nil {
		d = []byte{}
	}
	fmt.Fprint(rw, string(d))
	rw.WriteHeader(http.StatusOK)
}

func (c *MatchContext) GetCommit(rw web.ResponseWriter, req *web.Request) {
	// m := c.Matches[c.sessionid]
	// fmt.Println("GetCommit", req.PathParams)
	// fmt.Println(m, req)
}

func (c *MatchContext) GetReset(rw web.ResponseWriter, req *web.Request) {
	// m := c.Matches[c.sessionid]
	// fmt.Println("GetReset", req.PathParams)
	// fmt.Println(m, req)
}

func (c *MatchContext) GetMoveAttempt(rw web.ResponseWriter, req *web.Request) {
	// m := c.Matches[c.sessionid]
	// fmt.Println(m, req)
}

func (c *MatchContext) GetMatchFinish(rw web.ResponseWriter, req *web.Request) {
	// buf := make([]byte, 1024)
	// _, err := req.Body.Read(buf)
	// if err != nil {
	// 	fmt.Println(err.Error())
	// 	return
	// }
	// fmt.Println(buf)
	// m := c.Matches[c.sessionid]
	// fmt.Println(m)
}
