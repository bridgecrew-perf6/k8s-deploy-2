package kube

type IResource interface {
	Create() error
	Update() error
	Delete() error
	Get() error
}
