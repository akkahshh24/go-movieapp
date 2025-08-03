package discovery

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"time"

	"github.com/akkahshh24/movieapp/pkg/model"
)

// Registry defines a service registry.
type Registry interface {
	// Register creates a service instance record in the registry.
	Register(ctx context.Context, instanceID model.InstanceID, serviceName model.ServiceName, hostPort string) error
	// Deregister removes a service insttance record from the registry.
	Deregister(ctx context.Context, instanceID model.InstanceID, serviceName model.ServiceName) error
	// ServiceEndpoints returns the list of addresses of active instances of the given service.
	ServiceEndpoints(ctx context.Context, serviceName model.ServiceName) ([]string, error)
	// ReportHealthyState is a push mechanism for reporting healthy state to the registry.
	ReportHealthyState(instanceID model.InstanceID, serviceName model.ServiceName) error
}

// ErrNotFound is returned when no service addresses are found.
var ErrNotFound = errors.New("no service addresses found")

// GenerateInstanceID generates a pseudo-unique service instance identifier, using a service name
// suffixed by dash and a random number.
func GenerateInstanceID(serviceName model.ServiceName) model.InstanceID {
	return model.InstanceID(fmt.Sprintf("%s-%d", serviceName, rand.New(rand.NewSource(time.Now().UnixNano())).Int()))
}
