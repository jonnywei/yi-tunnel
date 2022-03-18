package common

type Config struct {
	Local_address            string
	Local_port               int
	Remote_address           string
	Remote_port              int
	Tunnel_config            string
	Secret                   string
	Stream_count_per_channel int
}
