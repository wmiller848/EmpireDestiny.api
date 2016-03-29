package context

import (
	"crypto/rand"
	"fmt"
	"net/http"
	"time"

	"github.com/wmiller848/EmpireDestiny/util"

	rethink "github.com/dancannon/gorethink"
	"github.com/gocraft/web"
)

type AuthSession struct {
	Id          string `AuthID`
	AccessToken string `AccessToken`
	Nonce       string `Nonce`
	RemoteIP    string `RemoteIP`
	PlayerID    string `PlayerID`
	Created     string `Created`
	Valid       string `Valid`
}

func (c *Context) SetAuthSession(authSession *AuthSession) error {
	resp, err := rethink.DB("ed").Table("auth_sessions").Insert(authSession).Run(c.Session)
	if err != nil {
		return err
	}
	fmt.Println("SetAuthSession", resp.Type())
	return nil
}

func (c *Context) GetAuthSession(id string) (*AuthSession, error) {
	cursor, err := rethink.DB("ed").Table("auth_sessions").Get(id).Run(c.Session)
	if err != nil {
		return nil, err
	}
	defer cursor.Close()
	authSession := AuthSession{}
	err = cursor.One(&authSession)
	if err != nil {
		return nil, err
	}
	return &authSession, nil
}

func (c *Context) GetAuth(rw web.ResponseWriter, req *web.Request) {
	c.Connect()
	defer c.Session.Close(rethink.CloseOpts{})
	// MOCK Get request_token FROM
	// TWITTER / FACEBOOK / GITHUB / ETC..
	unauth_request_token := make([]byte, 64)
	_, err := rand.Read(unauth_request_token)
	if err != nil {
		fmt.Println("GetAuth", err.Error())
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}
	// MOCK get access_token
	auth_request_token := util.Hash(string(unauth_request_token))
	access_token := util.Hash(auth_request_token)
	// Get secret
	nonce := make([]byte, 128)
	_, err = rand.Read(nonce)
	if err != nil {
		fmt.Println("GetAuth", err.Error())
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}
	playerID := make([]byte, 48)
	_, err = rand.Read(playerID)
	if err != nil {
		fmt.Println("GetAuth", err.Error())
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}
	_, err = c.GetPlayerInstance(util.Hex(playerID))
	if err != nil {
		_, err = c.CreatePlayer(util.Hex(playerID))
		if err != nil {
			fmt.Println("GetAuth", err.Error())
			rw.WriteHeader(http.StatusInternalServerError)
		}
	}
	access_token_client := util.Hash(util.Hex(nonce) + util.Hex(playerID))
	authSession := AuthSession{
		Nonce:       util.Hex(nonce),
		RemoteIP:    req.RemoteAddr,
		Id:          access_token_client,
		AccessToken: access_token,
		PlayerID:    util.Hex(playerID),
		Created:     time.Now().String(),
		Valid:       time.Now().Add(14 * 24 * time.Hour).String(),
	}
	err = c.SetAuthSession(&authSession)
	if err != nil {
		fmt.Println("GetAuth", err.Error())
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	cookie := &http.Cookie{
		Name:  "EDAUTH",
		Value: access_token_client,
		Path:  "/",
	}
	http.SetCookie(rw, cookie)
	rw.WriteHeader(http.StatusNoContent)
	fmt.Println("GetAuth", authSession.Nonce, authSession.PlayerID, authSession.Id)
}
