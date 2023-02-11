package core

import (
	"fmt"

	"github.com/gorilla/websocket"
	"github.com/sharadregoti/devops/model"
	"github.com/sharadregoti/devops/proto"
	"github.com/sharadregoti/devops/utils/logger"
)

type SessionManager struct {
	m map[string]*sessionInfo
}

type sessionInfo struct {
	conn *websocket.Conn
	c    *CurrentPluginContext
}

func NewSM() (*SessionManager, error) {
	return &SessionManager{
		m: make(map[string]*sessionInfo),
	}, nil
}

func (s *SessionManager) SessionCount() int {
	return len(s.m)
}

func (s *SessionManager) DeleteClient(ID string) {
	logger.LogInfo("Deleting client with id (%s)", ID)
	pCtx, err := s.GetClient(ID)
	if err != nil {
		return
	}
	pCtx.Close()
	delete(s.m, ID)
}

func (s *SessionManager) GetClient(ID string) (*CurrentPluginContext, error) {
	info, ok := s.m[ID]
	if !ok {
		return nil, fmt.Errorf("session with this id does not exists")
	}

	return info.c, nil
}

func (s *SessionManager) AddClient(ID, authID, contextID string) error {
	_, ok := s.m[ID]
	if !ok {
		pCtx, err := Start(false, &proto.AuthInfo{IdentifyingName: authID, Name: contextID})
		if err != nil {
			return err
		}

		dataPipe := make(chan model.WebsocketResponse, 10)
		pCtx.SetDataPipe(dataPipe)

		s.m[ID] = &sessionInfo{
			// conn: conn,
			c: pCtx,
		}
		logger.LogInfo("New client with ID (%s) is added", ID)
		return nil
	}
	return fmt.Errorf("session with this id already exists")
}
