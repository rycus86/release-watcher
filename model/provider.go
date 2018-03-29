package model

type Provider interface {
	Initialize()
	GetName() string
	Parse(interface{}) GenericProject
}
