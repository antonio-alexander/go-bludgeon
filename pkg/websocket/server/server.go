package server

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/antonio-alexander/go-bludgeon/common"
	"github.com/antonio-alexander/go-bludgeon/pkg/config"
	"github.com/antonio-alexander/go-bludgeon/pkg/logger"

	"github.com/gorilla/websocket"
	"github.com/pkg/errors"
)

type server struct {
	sync.RWMutex
	sync.WaitGroup
	*websocket.Conn
	logger.Logger
	lastPing     time.Time
	lastPong     time.Time
	connected    bool
	disconnected chan struct{}
	configured   bool
	config       *Configuration
}

func New(parameters ...interface{}) interface {
	Server
	common.Configurer
	common.Closer
} {
	var upgrader websocket.Upgrader
	var writer http.ResponseWriter
	var request *http.Request
	var header http.Header

	s := &server{Logger: logger.NewNullLogger()}
	for _, parameter := range parameters {
		switch p := parameter.(type) {
		case *Configuration:
			s.config = p
			s.configured = true
		case Configuration:
			s.config = &p
			s.configured = true
		case logger.Logger:
			s.Logger = p
		case http.ResponseWriter:
			writer = p
		case *http.Request:
			request = p
		case http.Header:
			header = p
		}
	}
	if !s.configured {
		s.Error(logAlias + "not configured")
		return nil
	}
	conn, err := upgrader.Upgrade(writer, request, header)
	if err != nil {
		s.Error(logAlias+"unable to upgrade: %s", err)
		return nil
	}
	s.Conn = conn
	s.connected = true
	s.disconnected = make(chan struct{})
	s.launchPinger()
	s.SetPongHandler(s.pongHandler)
	s.SetPingHandler(s.pingHandler)
	s.SetCloseHandler(s.closeHandler)
	return s
}

func (s *server) disconnect() bool {
	select {
	default:
		s.connected = false
		close(s.disconnected)
		return true
	case <-s.disconnected:
		return false
	}
}

func (s *server) pingHandler(ping string) error {
	if !s.connected {
		return ErrNotConnected
	}
	s.Trace(logAlias+"ping received: %s", ping)
	deadline := time.Now().Add(s.config.WriteTimeout)
	message := []byte(fmt.Sprint(time.Now()))
	if err := s.WriteControl(websocket.PongMessage, message, deadline); err != nil {
		s.Error(logAlias+"error while writing control pong: %s", err)
		return err
	}
	s.Trace(logAlias + "pong sent")
	return nil
}

func (s *server) pongHandler(pong string) error {
	if !s.connected {
		return ErrNotConnected
	}
	s.Trace(logAlias+"pong received: %s", pong)
	s.lastPong = time.Now()
	return nil
}

func (s *server) launchPinger() {
	started := make(chan struct{})
	s.Add(1)
	go func() {
		defer s.Done()

		pingFx := func() bool {
			switch {
			default:
				s.lastPing = time.Now()
				deadline := time.Now().Add(s.config.WriteTimeout)
				message := []byte(fmt.Sprint(s.lastPing))
				if err := s.WriteControl(websocket.PingMessage, message, deadline); err != nil {
					s.Error(logAlias+"error while writing control ping: %s", err)
					break
				}
				s.Trace(logAlias + "ping sent")
			case s.lastPong.Before(s.lastPing):
				if s.lastPong.Sub(s.lastPing) >= s.config.PingTimeout {
					if s.disconnect() {
						s.Debug(logAlias + "disconnected because ping timeout")
					}
					return true
				}
			}
			return false
		}
		tPing := time.NewTicker(s.config.PingInterval)
		defer tPing.Stop()
		close(started)
		s.Debug(logAlias + "pinger started")
		if pingFx() {
			return
		}
		for {
			select {
			case <-s.disconnected:
				return
			case <-tPing.C:
				if pingFx() {
					return
				}
			}
		}
	}()
	<-started
}

func (s *server) closeHandler(code int, text string) error {
	if !s.connected {
		return ErrNotConnected
	}
	s.Trace(logAlias+"close received code: %d, text: %s", code, text)
	if s.disconnect() {
		s.Debug(logAlias + "disconnected because close received")
	}
	return nil
}

func (s *server) Configure(items ...interface{}) error {
	s.Lock()
	defer s.Unlock()

	var c *Configuration

	for _, item := range items {
		switch v := item.(type) {
		case *Configuration:
			c = v
		}
	}
	if c == nil {
		return errors.New(config.ErrConfigurationNotFound)
	}
	s.config = c
	s.configured = true
	return nil
}

func (s *server) Write(item interface{}) error {
	s.Lock()
	defer s.Unlock()
	if !s.connected {
		return ErrNotConnected
	}
	if s.config.WriteTimeout > 0 {
		deadline := time.Now().Add(s.config.WriteTimeout)
		if err := s.SetWriteDeadline(deadline); err != nil {
			return err
		}
	}
	if err := s.Conn.WriteJSON(item); err != nil {
		if s.disconnect() {
			s.Debug(logAlias + "disconnected because write failed")
		}
		return err
	}
	return nil
}

func (s *server) Read(item interface{}) error {
	//KIM: you can safely read concurrently
	if !s.connected {
		return ErrNotConnected
	}
	if s.config.ReadTimeout > 0 {
		deadline := time.Now().Add(s.config.ReadTimeout)
		if err := s.SetReadDeadline(deadline); err != nil {
			return err
		}
	}
	if err := s.Conn.ReadJSON(item); err != nil {
		if s.disconnect() {
			s.Debug(logAlias + "disconnected because read failed")
		}
		return err
	}
	return nil
}

func (s *server) Close() {
	if !s.connected {
		return
	}
	deadline := time.Now().Add(s.config.WriteTimeout)
	bytes := []byte(fmt.Sprint(time.Now()))
	if err := s.WriteControl(websocket.CloseMessage, bytes, deadline); err != nil {
		s.Error("error while writing control: %s", err)
	}
	if s.disconnect() {
		s.Debug(logAlias + "disconnected because close executed")
	}
	s.Trace(logAlias+"close: %s", bytes)
}
