package relay

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/squeakycheese75/paytunnel/entities"
)

type Server struct {
	listenAddr string

	mu        sync.Mutex
	conn      *websocket.Conn
	sessionID string

	pending map[string]chan entities.ForwardResponseMessage
}

func NewServer(listenAddr string) *Server {
	return &Server{
		listenAddr: listenAddr,
		pending:    make(map[string]chan entities.ForwardResponseMessage),
	}
}

func (s *Server) Run() error {
	mux := http.NewServeMux()
	mux.HandleFunc("/ws", s.handleWS)
	mux.HandleFunc("/t/", s.handleTunnel)

	log.Printf("relay: relay listening on http://localhost%s\n", s.listenAddr)
	log.Printf("relay: websocket endpoint: ws://localhost%s/ws\n", s.listenAddr)

	return http.ListenAndServe(s.listenAddr, mux)
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

func (s *Server) handleWS(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, "failed to upgrade websocket", http.StatusBadRequest)
		return
	}

	defer func() {
		s.mu.Lock()
		if s.conn == conn {
			s.conn = nil
			s.sessionID = ""
		}
		s.mu.Unlock()
		_ = conn.Close()
	}()

	var register entities.RegisterMessage
	if err := conn.ReadJSON(&register); err != nil {
		return
	}

	sessionID := newSessionID()

	s.mu.Lock()
	s.conn = conn
	s.sessionID = sessionID
	s.mu.Unlock()

	registered := entities.RegisteredMessage{
		Type:      "registered",
		SessionID: sessionID,
		PublicURL: fmt.Sprintf("http://localhost%s/t/%s", s.listenAddr, sessionID),
	}

	if err := conn.WriteJSON(registered); err != nil {
		return
	}

	log.Printf("relay: tunnel client registered: session_id=%s target=%s\n", sessionID, register.TargetURL)

	for {
		var raw map[string]any
		if err := conn.ReadJSON(&raw); err != nil {
			return
		}

		msgType, _ := raw["type"].(string)
		if msgType != "forward_response" {
			continue
		}

		data, err := json.Marshal(raw)
		if err != nil {
			continue
		}

		var msg entities.ForwardResponseMessage
		if err := json.Unmarshal(data, &msg); err != nil {
			continue
		}

		s.mu.Lock()
		ch := s.pending[msg.RequestID]
		s.mu.Unlock()

		if ch != nil {
			ch <- msg
		}
	}
}

func (s *Server) handleTunnel(w http.ResponseWriter, r *http.Request) {
	s.mu.Lock()
	conn := s.conn
	sessionID := s.sessionID
	s.mu.Unlock()

	if conn == nil || sessionID == "" {
		http.Error(w, "no tunnel client connected", http.StatusServiceUnavailable)
		return
	}

	prefix := "/t/" + sessionID
	if !strings.HasPrefix(r.URL.Path, prefix) {
		http.NotFound(w, r)
		return
	}

	path := strings.TrimPrefix(r.URL.Path, prefix)
	if path == "" {
		path = "/"
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "failed to read request body", http.StatusBadRequest)
		return
	}

	requestID := fmt.Sprintf("req_%d", time.Now().UnixNano())

	msg := entities.ForwardRequestMessage{
		Type:       "forward_request",
		RequestID:  requestID,
		Method:     r.Method,
		Path:       path,
		RawQuery:   r.URL.RawQuery,
		Headers:    r.Header,
		BodyBase64: base64.StdEncoding.EncodeToString(body),
	}

	respCh := make(chan entities.ForwardResponseMessage, 1)

	s.mu.Lock()
	s.pending[requestID] = respCh
	s.mu.Unlock()

	defer func() {
		s.mu.Lock()
		delete(s.pending, requestID)
		s.mu.Unlock()
	}()

	if err := conn.WriteJSON(msg); err != nil {
		http.Error(w, "failed to forward request to tunnel client", http.StatusBadGateway)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 30*time.Second)
	defer cancel()

	select {
	case resp := <-respCh:
		for key, values := range resp.Headers {
			for _, value := range values {
				w.Header().Add(key, value)
			}
		}

		w.WriteHeader(resp.StatusCode)

		body, err := base64.StdEncoding.DecodeString(resp.BodyBase64)
		if err == nil {
			_, _ = w.Write(body)
		}

	case <-ctx.Done():
		http.Error(w, "timed out waiting for tunnel response", http.StatusGatewayTimeout)
	}
}

func newSessionID() string {
	return fmt.Sprintf("%x", time.Now().UnixNano())
}
