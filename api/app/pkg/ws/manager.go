package ws

import (
	log "api/app/pkg/logger"
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"net/http"
	"sync"
	"sync/atomic"
)

type Manager struct {
	upgrader        websocket.Upgrader
	processorFabric PipeProcessorFabric
	clients         map[string]Client

	isClosed   atomic.Bool
	deadSignal chan string
	mu         sync.Mutex
	logger     Logger
}

func NewManager(
	opts ...OptionFunc,
) *Manager {
	manager := &Manager{
		upgrader:        websocket.Upgrader{},
		processorFabric: &PipeProcessorFabricImpl{},
		clients:         make(map[string]Client, defaultConnsLimit),
		deadSignal:      make(chan string, defaultConnsLimit),
		logger:          log.Default(),
		mu:              sync.Mutex{},
		isClosed:        atomic.Bool{},
	}
	return manager.With(opts...)
}

func (m *Manager) Process(uniqueID int, w http.ResponseWriter, r *http.Request, header http.Header) error {
	conn, err := m.upgrader.Upgrade(w, r, header)
	if err != nil {
		return err
	}

	processor, err := m.processorFabric.NewPipeProcessor(r.Context(), uniqueID)
	if err != nil {
		return err
	}
	if m.isClosed.Load() {
		return ErrManagerClosed
	}
	client := NewDefaultClient(conn, uuid.New().String(), m.deadSignal, processor, m.logger)
	m.addClient(client.GetClientID(), client)

	go func() {
		err := client.Run(context.WithoutCancel(r.Context()))
		if err != nil {
			m.logger.Errorw("DefaultClient run error", "error", err)
		}
	}()
	return nil
}
func (m *Manager) addClient(id string, client Client) {
	defer m.mu.Unlock()
	m.mu.Lock()
	m.clients[id] = client
}
func (m *Manager) Run(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case id, ok := <-m.deadSignal:
			if !ok {
				return
			}
			func() {
				defer m.mu.Unlock()
				m.mu.Lock()
				_, ok := m.clients[id]
				if !ok {
					m.logger.Errorw("DefaultClient not found", "id", id)
					return
				}
				delete(m.clients, id)
			}()
		}
	}
}

func (m *Manager) Close() error {
	var err error
	m.isClosed.Store(true)
	for id, client := range m.clients {
		err := client.Close()
		if err != nil {
			m.logger.Errorw("error closing connection", "error", err)
			err = errors.Join(err, err)
		}
		delete(m.clients, id)
	}
	close(m.deadSignal)
	return err
}
