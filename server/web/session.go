package web

import (
	"base"
	"net/http"

	"github.com/gorilla/sessions"
)

const (
	session_store    = "team-ide-session"
	session_key      = "team-ide"
	session_bean_key = "team-ide-bean"
)

var sessionStore = sessions.NewCookieStore([]byte(session_store))

func init() {
	sessionStore.MaxAge(60 * 60 * 2)
}

func SetSessionUser(w http.ResponseWriter, r *http.Request, user *base.LoginUserBean) error {
	param := GetSessionBean(r)
	param.User = user
	return SetSessionBean(w, r, param)
}

func SetSessionBean(w http.ResponseWriter, r *http.Request, sessionBean *base.SessionBean) error {
	var err error
	session, _ := sessionStore.Get(r, session_key)
	value := ""
	if sessionBean != nil {
		var by []byte
		by, err = base.JSON.Marshal(sessionBean)
		if err != nil {
			return err
		}
		value = string(by)
	}
	session.Values[session_bean_key] = value
	err = session.Save(r, w)
	return err
}

func GetSessionBean(r *http.Request) *base.SessionBean {
	session, _ := sessionStore.Get(r, session_key)
	value, ok := session.Values[session_bean_key]
	res := &base.SessionBean{}
	if ok && value != "" {
		base.JSON.Unmarshal([]byte(value.(string)), res)
	} else {
		res = &base.SessionBean{}
	}
	return res
}
