package main

type config struct {
	API              apiConfig              `yaml:"api"`
	ServiceDiscovery serviceDiscoveryConfig `yaml:"serviceDiscovery"`
	MessageQueue     messageQueueConfig     `yaml:"messageQueue"`
	Database         databaseConfig         `yaml:"database"`
}

type apiConfig struct {
	Port int `yaml:"port"`
}

type serviceDiscoveryConfig struct {
	Name   string       `yaml:"name"`
	Consul consulConfig `yaml:"consul"`
}

type consulConfig struct {
	Address string `yaml:"address"`
}

type messageQueueConfig struct {
	Address string `yaml:"address"`
	GroupID string `yaml:"groupID"`
	Topic   string `yaml:"topic"`
}

type databaseConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	DBName   string `yaml:"dbname"`
}
