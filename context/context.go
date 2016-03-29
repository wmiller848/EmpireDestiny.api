package context

import (
	"fmt"
	"net/http"
	"regexp"
	"sync"

	"github.com/wmiller848/EmpireDestiny/queue"

	"github.com/wmiller848/EmpireDestiny/game"

	rethink "github.com/dancannon/gorethink"
	"github.com/gocraft/web"
)

var GlobalPtr *Global
var Requests *RequestCounter = &RequestCounter{}

type RequestCounter struct {
	sync.Mutex
	Value uint32
}

// TODO
// Make all this shit a db or cache
type Global struct {
	sync.Mutex
	Game       *game.Game
	Queue      *queue.Queue
	MatchRegex *regexp.Regexp
}

type Context struct {
	*Global
	Session *rethink.Session
}

type PlayerContext struct {
	*Context
	playerid string
}

type MatchContext struct {
	*PlayerContext
	sessionid string
}

func (c *Context) Count(rw web.ResponseWriter, req *web.Request, next web.NextMiddlewareFunc) {
	Requests.Lock()
	Requests.Value++
	Requests.Unlock()
	fmt.Println("Requests -", Requests.Value)
	next(rw, req)
}

func (c *Context) Connect() {
	sess, err := rethink.Connect(rethink.ConnectOpts{
		Address: "localhost:28015",
		// Addresses: []string{"localhost:28015", "localhost:28016"},
		// Database: "test",
		// AuthKey:       "14daak1cad13dj",
		// DiscoverHosts: true,
		// Timeout:       1 * time.Second,
		// MaxIdle:       100,
		// MaxOpen:       100,
	})
	fmt.Println("RETHINK SESSION")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	if sess != nil {
		fmt.Println("Connecting to rethink")
		c.Session = sess
	}
	c.Global = GlobalPtr
}

func (c *PlayerContext) SetGlobal(rw web.ResponseWriter, req *web.Request, next web.NextMiddlewareFunc) {
	c.Connect()
	fmt.Println("SetGlobal - Forwarding")
	next(rw, req)
}

func (c *PlayerContext) PlayerMiddleware(rw web.ResponseWriter, req *web.Request, next web.NextMiddlewareFunc) {
	authid, err := req.Cookie("EDAUTH")
	_authid := ""
	if err != nil {
		_authid = ""
	} else {
		_authid = authid.Value
	}
	fmt.Println("PlayerMiddleware", req.RequestURI, _authid)
	authSession, err := c.GetAuthSession(_authid)
	if (err != nil) && req.RequestURI != "/auth" {
		fmt.Println("Redirect to /auth")
		rw.Header().Set("Location", "/auth")
		rw.WriteHeader(http.StatusUnauthorized)
		return
	} else if err == nil {
		c.playerid = authSession.PlayerID
	}
	fmt.Println("PlayerMiddleware - Forwarding")
	rw.Header().Set("Content-Type", "application/json")
	next(rw, req)
}

func (c *MatchContext) MatchSessionMiddleware(rw web.ResponseWriter, req *web.Request, next web.NextMiddlewareFunc) {
	sessionid, err := req.Cookie("EDSESSION")
	ssid := ""
	if err != nil {
		fmt.Println("MatchSessionMiddleware", err.Error())
		ssid = ""
	} else {
		ssid = sessionid.Value
	}

	fmt.Println("MatchSessionMiddleware", req.RequestURI, ssid)
	_, err = c.GetMatchInstance(ssid)
	if (ssid == "" || err != nil) &&
		c.MatchRegex.Match([]byte(req.RequestURI)) == false {
		rw.WriteHeader(http.StatusUnauthorized)
		return
	} else {
		c.sessionid = ssid
	}
	fmt.Println("MatchSessionMiddleware - Forwarding")
	next(rw, req)
}

func (c *Context) Root(rw web.ResponseWriter, req *web.Request) {
	fmt.Fprintf(rw, "<html>Empire Destiny</html>")
}
