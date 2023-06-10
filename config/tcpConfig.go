package config

// TcpServerConfiguration define global tcp server config
type TcpServerConfiguration struct {
	Host       string
	Port       int
	MaxClients int
}

// Configuration is global tcp server config
var Configuration *TcpServerConfiguration

// Init init global tcp server config
func Init() *TcpServerConfiguration {
	Configuration = &TcpServerConfiguration{
		Host:       "127.0.0.1",
		Port:       6379,
		MaxClients: 10000,
	}
	return Configuration
}
