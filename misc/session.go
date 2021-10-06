package misc

import (
	"crypto/rand"
	"io"
	"log"
)

type SessionManager struct {
	sessions map[string]Session
}

var GSS SessionManager

type Session struct {
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

func (s *SessionManager) Get(id string) (Session, error) {
	sess, pre := s.sessions[id]
	if pre {
		return sess, nil
	} else {
		return Session{}, nil
	}
}

func (s *SessionManager) AddSession(user string) string {
	id := s.generateID()
	s.sessions[id] = Session{id, user}
	return id
}
