package pluginmanager

import (
	"fmt"

	"github.com/sharadregoti/devops-plugin-sdk/proto"
	"github.com/sharadregoti/devops/model"
	"github.com/sharadregoti/devops/utils/logger"
)

type SessionManager struct {
	conf *model.Config
	m    map[string]*sessionInfo
}

type sessionInfo struct {
	c *CurrentPluginContext
}

func NewSM(conf *model.Config) (*SessionManager, error) {
	return &SessionManager{
		conf: conf,
		m:    make(map[string]*sessionInfo),
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

func (s *SessionManager) AddClient(ID, pluginName, authID, contextID string) error {
	_, ok := s.m[ID]
	if !ok {
		pCtx, err := Start(pluginName, s.conf, &proto.AuthInfo{IdentifyingName: authID, Name: contextID})
		if err != nil {
			return err
		}

		dataPipe := make(chan model.WebsocketResponse, 10)
		pCtx.SetDataPipe(dataPipe)

		s.m[ID] = &sessionInfo{
			c: pCtx,
		}
		logger.LogInfo("New client with ID (%s) is added", ID)
		return nil
	}
	return fmt.Errorf("session with this id already exists")
}
