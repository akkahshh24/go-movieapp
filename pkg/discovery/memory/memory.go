package memory

import (
	"context"
	"errors"
	"log"
	"sync"
	"time"

	"github.com/akkahshh24/movieapp/pkg/discovery"
	"github.com/akkahshh24/movieapp/pkg/model"
)

// Registry defines an in-memory service regisry.
// Note: this registry does not perform health monitoring of active instances.
type Registry struct {
	sync.RWMutex
	serviceAddrs map[model.ServiceName]map[model.InstanceID]*model.ServiceInstance
}

// NewRegistry creates a new in-memory service registry instance.
func NewRegistry() *Registry {
	return &Registry{serviceAddrs: map[model.ServiceName]map[model.InstanceID]*model.ServiceInstance{}}
}

// Register creates a service record in the registry.
func (r *Registry) Register(ctx context.Context, instanceID model.InstanceID, serviceName model.ServiceName, hostPort string) error {
	r.Lock()
	defer r.Unlock()
	if _, ok := r.serviceAddrs[serviceName]; !ok {
		r.serviceAddrs[serviceName] = map[model.InstanceID]*model.ServiceInstance{}
	}
	r.serviceAddrs[serviceName][instanceID] = &model.ServiceInstance{HostPort: hostPort, LastActive: time.Now()}
	return nil
}

// Deregister removes a service record from the registry.
func (r *Registry) Deregister(ctx context.Context, instanceID model.InstanceID, serviceName model.ServiceName) error {
	r.Lock()
	defer r.Unlock()
	if _, ok := r.serviceAddrs[serviceName]; !ok {
		return nil
	}
	delete(r.serviceAddrs[serviceName], instanceID)
	return nil
}

// ReportHealthyState is a push mechanism for reporting healthy state to the registry.
func (r *Registry) ReportHealthyState(instanceID model.InstanceID, serviceName model.ServiceName) error {
	r.Lock()
	defer r.Unlock()
	if _, ok := r.serviceAddrs[serviceName]; !ok {
		return errors.New("instance " + instanceID.String() + " of service " + serviceName.String() + " is not registered yet")
	}
	if _, ok := r.serviceAddrs[serviceName][instanceID]; !ok {
		return errors.New("service instance is not registered yet")
	}
	r.serviceAddrs[serviceName][instanceID].LastActive = time.Now()
	return nil
}

// ServiceAddresses returns the list of addresses of active instances of the given service.
func (r *Registry) ServiceEndpoints(ctx context.Context, serviceName model.ServiceName) ([]string, error) {
	r.RLock()
	defer r.RUnlock()
	if len(r.serviceAddrs[serviceName]) == 0 {
		return nil, discovery.ErrNotFound
	}
	var res []string
	for instanceID, serviceInstance := range r.serviceAddrs[serviceName] {
		// Only return instances with a successful health check within the last 5 seconds.
		if serviceInstance.LastActive.Before(time.Now().Add(-5 * time.Second)) {
			log.Println("Instance " + instanceID.String() + " of service " + serviceName.String() + " is not active, skipping")
			continue
		}
		res = append(res, serviceInstance.HostPort)
	}
	return res, nil
}
