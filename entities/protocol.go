package entities

type RegisterMessage struct {
	Type      string `json:"type"`
	TargetURL string `json:"target_url"`
}

type RegisteredMessage struct {
	Type      string `json:"type"`
	SessionID string `json:"session_id"`
	PublicURL string `json:"public_url"`
}

type ForwardRequestMessage struct {
	Type       string              `json:"type"`
	RequestID  string              `json:"request_id"`
	Method     string              `json:"method"`
	Path       string              `json:"path"`
	RawQuery   string              `json:"raw_query"`
	Headers    map[string][]string `json:"headers"`
	BodyBase64 string              `json:"body_base64"`
}

type ForwardResponseMessage struct {
	Type       string              `json:"type"`
	RequestID  string              `json:"request_id"`
	StatusCode int                 `json:"status_code"`
	Headers    map[string][]string `json:"headers"`
	BodyBase64 string              `json:"body_base64"`
}
