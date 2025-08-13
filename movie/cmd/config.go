package main

type config struct {
	API              apiConfig              `yaml:"api"`
	ServiceDiscovery serviceDiscoveryConfig `yaml:"serviceDiscovery"`
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
