package utils

// Settings API config struct
type Settings struct {
	ApiParams      Params       `json:"api"`
	PostgresParams PostgresSets `json:"postgresql"`
	Business       Business     `json:"business"`
}

type Params struct {
	ServerName string `json:"server_name"`
	PortRun    int    `json:"port_run"`
	LogFile    string `json:"log_file"`
	Secret     string `json:"secret"`
}

type PostgresSets struct {
	User     string `json:"user"`
	Password string `json:"password"`
	Server   string `json:"server"`
	Port     int    `json:"port"`
	Schema   string `json:"schema"`
}

type Business struct {
	MaxBalance int `json:"max_balance"`
	MinBalance int `json:"min_balance"`
}
