package main

import (
	"crypto/rand"
	"fmt"
	"io"
	"log"
	"net/http"
)

type SessionManager struct {
	sessions map[string]session
}

type session struct {
	s_id string
	user string
}

func (s *SessionManager) generateID() string {
	buf := make([]byte, 32)
	_, err := io.ReadFull(rand.Reader, buf)
	if err != nil {
		log.Fatal(err)
	}
	return string(buf)
}

func (s *SessionManager) Get(id string) (session, error) {
	sess, pre := s.sessions[id]
	if pre {
		return sess, nil
	} else {
		return session{}, nil
	}
}

func (s *SessionManager) AddSession(user string) string {
	id := s.generateID()
	s.sessions[id] = session{id, user}
	return id
}

var sessionManager SessionManager

func login(w http.ResponseWriter, r *http.Request) {

	id, pre := r.Cookie("user")
	if pre == nil {
		sess, err := sessionManager.Get(id.Value)
		if err != nil {
			w.Write([]byte("invalid session id"))
			return
		} else {
			msg := fmt.Sprintf("Your name %v , session id%v", sess.user, sess.s_id)
			w.Write([]byte(msg))
		}
	} else {

	}
	cookie := &http.Cookie{
		Name:  "s_id",
		Value: sessionManager.generateID(),
	}
	http.SetCookie(w, cookie)
	w.Write([]byte("Hello,world"))
}

func browse(w http.ResponseWriter, r *http.Request) {
	id, _ := r.Cookie("user")
	msg := fmt.Sprintf("Your session id is [%v]", id.Value)
	w.Write([]byte(msg))
}

func init() {
	sessionManager = make(map[string]session)
}

func main() {

	http.HandleFunc("/login", login)
	http.HandleFunc("/browse", browse)
	http.ListenAndServe(":10000", nil)

	// var w http.ResponseWriter
}
