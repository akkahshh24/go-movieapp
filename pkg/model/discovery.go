package model

import "time"

// ServiceName defines a service name.
type ServiceName string

func (s ServiceName) String() string {
	return string(s)
}

// InstanceID defines a service instance identifier.
type InstanceID string

func (i InstanceID) String() string {
	return string(i)
}

// ServiceInstance defines a service instance with its metadata.
type ServiceInstance struct {
	ServiceName ServiceName
	InstanceID  InstanceID
	HostPort    string
	LastActive  time.Time
}
