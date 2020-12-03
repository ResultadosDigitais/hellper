package model

// ServiceInstance represents an instance of a service that is running somewhere
type ServiceInstance struct {
	ID   int64  `db:"id,omitempty"`
	Name string `db:"name,omitempty"`
}
