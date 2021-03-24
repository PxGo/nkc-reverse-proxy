package config

type Server struct {
	Name           string   `json:"name"`
	Type           string   `json:"type"`
	Hostname       []string `json:"hostname"`
	Https          bool     `json:"https"`
	HttpTarget     []string `json:"httpTarget"`
	WsTarget       []string `json:"wsTarget"`
	NoResponsePage string   `json:"noResponsePage"`
	SSL            SSL      `json:"SSL"`
}

type Ports struct {
	HttpPort  int64 `json:"httpPort"`
	HttpsPort int64 `json:"httpsPort"`
}

type Profile struct {
	Ports           Ports    `json:"ports"`
	Servers         []Server `json:"servers"`
	NoResponsePage  string   `json:"noResponsePage"`
	SSL             SSL      `json:"SSL"`
	ConcurrentLimit int64    `json:"concurrentLimit"`
}

type SSL struct {
	Cert string `json:"cert"`
	Key  string `json:"key"`
}
