package tunnel

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/squeakycheese75/paytunnel/entities"
)

type Client struct {
	relayURL  string
	targetURL string
}

func NewClient(relayURL, targetURL string) *Client {
	return &Client{
		relayURL:  relayURL,
		targetURL: targetURL,
	}
}

func (c *Client) Run() error {
	conn, _, err := websocket.DefaultDialer.Dial(c.relayURL, nil)
	if err != nil {
		return fmt.Errorf("connect to relay: %w", err)
	}
	defer conn.Close()

	register := entities.RegisterMessage{
		Type:      "register",
		TargetURL: c.targetURL,
	}
	if err := conn.WriteJSON(register); err != nil {
		return fmt.Errorf("send register: %w", err)
	}

	var registered entities.RegisteredMessage
	if err := conn.ReadJSON(&registered); err != nil {
		return fmt.Errorf("read registered response: %w", err)
	}

	log.Println("tunnel: connected to relay")
	log.Printf("tunnel: public_url=%s", registered.PublicURL)
	log.Printf("tunnel: forwarding_to=%s", c.targetURL)

	for {
		var raw map[string]any
		if err := conn.ReadJSON(&raw); err != nil {
			return err
		}

		msgType, _ := raw["type"].(string)
		if msgType != "forward_request" {
			continue
		}

		data, err := json.Marshal(raw)
		if err != nil {
			return err
		}

		var reqMsg entities.ForwardRequestMessage
		if err := json.Unmarshal(data, &reqMsg); err != nil {
			return err
		}

		respMsg, err := c.forwardRequest(reqMsg)
		if err != nil {
			respMsg = entities.ForwardResponseMessage{
				Type:       "forward_response",
				RequestID:  reqMsg.RequestID,
				StatusCode: http.StatusBadGateway,
				Headers: map[string][]string{
					"Content-Type": {"text/plain; charset=utf-8"},
				},
				BodyBase64: base64.StdEncoding.EncodeToString([]byte(err.Error())),
			}
		}

		if err := conn.WriteJSON(respMsg); err != nil {
			return err
		}
	}
}

func (c *Client) forwardRequest(reqMsg entities.ForwardRequestMessage) (entities.ForwardResponseMessage, error) {
	body, err := base64.StdEncoding.DecodeString(reqMsg.BodyBase64)
	if err != nil {
		return entities.ForwardResponseMessage{}, fmt.Errorf("decode request body: %w", err)
	}

	url := c.targetURL + reqMsg.Path
	if reqMsg.RawQuery != "" {
		url += "?" + reqMsg.RawQuery
	}

	req, err := http.NewRequest(reqMsg.Method, url, bytes.NewBuffer(body))
	if err != nil {
		return entities.ForwardResponseMessage{}, fmt.Errorf("build local request: %w", err)
	}

	for key, values := range reqMsg.Headers {
		for _, value := range values {
			req.Header.Add(key, value)
		}
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return entities.ForwardResponseMessage{}, fmt.Errorf("send local request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return entities.ForwardResponseMessage{}, fmt.Errorf("read local response: %w", err)
	}

	return entities.ForwardResponseMessage{
		Type:       "forward_response",
		RequestID:  reqMsg.RequestID,
		StatusCode: resp.StatusCode,
		Headers:    resp.Header,
		BodyBase64: base64.StdEncoding.EncodeToString(respBody),
	}, nil
}
