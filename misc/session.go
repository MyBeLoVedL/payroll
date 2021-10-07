package misc

import (
	"payroll/db/util"
)

type SessionManager struct {
	sessions map[string]Session
}

var GSS SessionManager

func init() {
	GSS = SessionManager{}
	GSS.sessions = make(map[string]Session)
}

type Session struct {
	S_id string
	User int64
}

func (s *SessionManager) generateID() string {
	return util.RandStr(32)
}

func (s *SessionManager) Get(id string) (Session, error) {
	sess, pre := s.sessions[id]
	if pre {
		return sess, nil
	} else {
		return Session{}, nil
	}
}

func (s *SessionManager) AddSession(user int64) string {
	id := s.generateID()
	s.sessions[id] = Session{id, user}
	return id
}
