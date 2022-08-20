package restclient

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	client "github.com/antonio-alexander/go-bludgeon/changes/client"
	data "github.com/antonio-alexander/go-bludgeon/changes/data"
	"github.com/antonio-alexander/go-bludgeon/internal"

	internal_logger "github.com/antonio-alexander/go-bludgeon/internal/logger"
	internal_websocketclient "github.com/antonio-alexander/go-bludgeon/internal/websocket/client"
)

type handler struct {
	sync.RWMutex
	sync.WaitGroup
	internal_logger.Logger
	client interface {
		internal_websocketclient.Client
		internal.Configurer
		internal.Closer
	}
	logAlias     string
	config       *Configuration
	ctx          context.Context
	cancel       context.CancelFunc
	connected    bool
	disconnected chan struct{}
	handlerFx    client.HandlerFx
}

func newHandler(ctx context.Context, logger internal_logger.Logger, handlerId string, config *Configuration, handlerFx client.HandlerFx) *handler {
	ctx, cancel := context.WithCancel(ctx)
	h := &handler{
		Logger:       logger,
		client:       internal_websocketclient.New(),
		handlerFx:    handlerFx,
		disconnected: make(chan struct{}, 1),
		config:       config,
		logAlias:     logAlias + "[" + handlerId + "] ",
		ctx:          ctx,
		cancel:       cancel,
	}
	h.client.Configure(config.Websocket)
	h.launchConnect()
	h.launchChangeReader()
	return h
}

func (h *handler) Close() {
	h.Lock()
	defer h.Unlock()

	h.cancel()
	h.Wait()
	h.client.Close()
}

func (h *handler) launchConnect() {
	started := make(chan struct{})
	h.Add(1)
	go func() {
		defer h.Done()

		connectFx := func() bool {
			uri := fmt.Sprintf("ws://%s:%s"+data.RouteChangesWebsocket,
				h.config.Rest.Address, h.config.Rest.Port)
			response, err := h.client.Connect(h.ctx, uri, http.Header{})
			if err != nil {
				h.Error(logAlias+"error while connecting to websocket", err)
				return false
			}
			defer response.Body.Close()
			return true
		}
		tConnect := time.NewTicker(10 * time.Second)
		defer tConnect.Stop()
		close(started)
		if connectFx() {
			h.connected = true
			tConnect.Stop()
		}
		for {
			select {
			case <-h.ctx.Done():
				return
			case <-h.disconnected:
				if connectFx() {
					h.connected = true
					break
				}
				tConnect = time.NewTicker(10 * time.Second)
			case <-tConnect.C:
				if h.connected {
					break
				}
				if connectFx() {
					h.connected = true
					h.Trace(h.logAlias + "connected")
					tConnect.Stop()
				}
			}
		}
	}()
	<-started
}

func (h *handler) launchChangeReader() error {
	started := make(chan struct{})
	h.Add(1)
	go func() {
		defer h.Done()

		close(started)
		for {
			select {
			case <-h.ctx.Done():
				return
			default:
				if !h.connected {
					time.Sleep(h.config.Websocket.ReadTimeout)
					break
				}
				wrapper := &data.Wrapper{}
				if err := h.client.Read(wrapper); err != nil {
					h.Error(h.logAlias+"error while reading websocket: %s", err)
					break
				}
				switch data.MessageType(wrapper.Type) {
				case data.MessageTypeChange:
					change := &data.Change{}
					if err := json.Unmarshal(wrapper.Bytes, change); err != nil {
						h.Error(h.logAlias+"error while unmarshalling change: %s", err)
						break
					}
					if err := h.handlerFx(change); err != nil {
						h.Error(h.logAlias+"error while handling change: %s", err)
						break
					}
				case data.MessageTypeChangeDigest:
					changeDigest := &data.ChangeDigest{}
					if err := json.Unmarshal(wrapper.Bytes, changeDigest); err != nil {
						h.Error(logAlias+"error while unmarshalling changes: %s", err)
						break
					}
					if err := h.handlerFx(changeDigest.Changes...); err != nil {
						h.Error(logAlias+"error while handling changes: %s", err)
						break
					}
				}
			}
		}
	}()
	<-started
	return nil
}

func (h *handler) isConnected() bool {
	return h.connected
}
