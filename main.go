package main

import (
	"fmt"
	"log"
	"net/http"
	"regexp"

	"github.com/wmiller848/EmpireDestiny/context"
	"github.com/wmiller848/EmpireDestiny/game"
	"github.com/wmiller848/EmpireDestiny/queue"

	rethink "github.com/dancannon/gorethink"
	"github.com/gocraft/web"
)

var global context.Global

func ensure_users() {

}

func ensure_auth_sessions() {

}

func ensure_cards() {

}

func ensure_matches() {

}

func ensure_tables(session *rethink.Session) {
	cursor, err := rethink.TableList().Run(session)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	tables := []interface{}{}
	err = cursor.All(&tables)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	tablesMap := make(map[string]bool)
	for _, db := range tables {
		tablesMap[db.(string)] = true
	}
	if !tablesMap["ed"] {
		_, err := rethink.TableCreate("users").RunWrite(session)
		if err != nil {
			fmt.Print(err.Error())
			return
		}
	}
}

func ensure_db(session *rethink.Session) {
	cursor, err := rethink.DBList().Run(session)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	dbs := []interface{}{}
	err = cursor.All(&dbs)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	dbsMap := make(map[string]bool)
	for _, db := range dbs {
		dbsMap[db.(string)] = true
	}
	if !dbsMap["ed"] {
		_, err := rethink.DBCreate("ed").RunWrite(session)
		if err != nil {
			fmt.Print(err.Error())
			return
		}
	}
}

func setup() {
	rethinkSession, err := rethink.Connect(rethink.ConnectOpts{
		Address: "localhost:28015",
		// Addresses: []string{"localhost:28015", "localhost:28016"},
		// Database: "test",
		// AuthKey:       "14daak1cad13dj",
		DiscoverHosts: true,
		// Timeout:       1 * time.Second,
		// MaxIdle:       100,
		// MaxOpen:       100,
	})
	if err != nil {
		log.Fatalln(err.Error())
		return
	}
	ensure_db(rethinkSession)
	ensure_tables(rethinkSession)
}

func start() {
	rethinkSession, err := rethink.Connect(rethink.ConnectOpts{
		Address: "localhost:28015",
		// Addresses: []string{"localhost:28015", "localhost:28016"},
		// Database: "test",
		// AuthKey:       "14daak1cad13dj",
		DiscoverHosts: true,
		// Timeout:       1 * time.Second,
		// MaxIdle:       100,
		// MaxOpen:       100,
	})
	if err != nil {
		log.Fatalln(err.Error())
		return
	}

	gameInstance, err := game.LoadGameFromYML("game/data/game.yml")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	chn := gameInstance.Sync()
	for {
		crdRevison, closed := <-chn
		if closed == false {
			break
		}
		fmt.Println("Card Revision", crdRevison)
		resp, err := rethink.DB("ed").Table("cards").Insert(crdRevison, rethink.InsertOpts{
			Conflict: "replace",
		}).RunWrite(rethinkSession)
		if err != nil {
			fmt.Println(err.Error())
		} else {
			fmt.Println(resp)
		}
	}

	exp, err := regexp.Compile(`^/player/match/queue/[0-9a-zA-z]*`)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	requests := context.RequestCounter{}
	global := context.Global{
		Game:       gameInstance,
		Queue:      queue.CreateQueue(),
		MatchRegex: exp,
	}
	context.GlobalPtr = &global

	fmt.Println(requests, global)

	global.Lock()
	global.Queue = queue.CreateQueue()
	global.Unlock()
}

func main() {

	//setup()
	go start()
	rootRouter := web.New(context.Context{}). // Create your router
							Middleware(web.LoggerMiddleware). // included logging middleware
							Middleware(web.ShowErrorsMiddleware).
							Middleware((*context.Context).Count).     // Set the global indexes
							Get("/", (*context.Context).Root).        //
							Get("/auth", (*context.Context).GetAuth). //
							Head("/auth", (*context.Context).GetAuth) //

	playerRouter := rootRouter.Subrouter(context.PlayerContext{}, "/player").
		Middleware((*context.PlayerContext).SetGlobal).        // Set the global indexes
		Middleware((*context.PlayerContext).PlayerMiddleware). // User Session Validation
		Get("/", (*context.PlayerContext).GetPlayer).          //
		Post("/decks", (*context.PlayerContext).PostDecks).    //
		Put("/decks/:id", (*context.PlayerContext).PutDecks).  //
		Get("/game", (*context.PlayerContext).GetGame)         //

	matchRouter := playerRouter.Subrouter(context.MatchContext{}, "/match").
		// Middleware((*context.MatchContext).PlayerMiddleware).       // User Session Validation
		Middleware((*context.MatchContext).MatchSessionMiddleware). // Match Session Validation
		Get("/queue/:deckid", (*context.MatchContext).GetMatch).    //
		Get("/finish", (*context.MatchContext).GetMatchFinish).     //
		Get("/info", (*context.MatchContext).GetInfo).              //
		// Get("/move/engage/:id/taget/:targetid", (*context.MatchContext).GetEngage). //
		// Get("/move/disengage/:id", (*context.MatchContext).GetDisengage).           //
		// Get("/move/bow/:id/target/:targetid", (*context.MatchContext).GetBow).      //
		// Get("/move/unbow/:id", (*context.MatchContext).GetUnbow).                   //
		Get("/move/commit", (*context.MatchContext).GetCommit).
		Get("/move/reset", (*context.MatchContext).GetReset)

	fmt.Println(rootRouter, playerRouter, matchRouter)
	http.ListenAndServe("localhost:3000", rootRouter)
}
