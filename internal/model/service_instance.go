package model

// ServiceInstance represents an instance of a service that is running somewhere
type ServiceInstance struct {
	ID   string `db:"id,omitempty"`
	Name string `db:"name,omitempty"`
}
