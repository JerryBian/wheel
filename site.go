package main

type Site struct {
	Name string `yaml:"name"`
	Description string
	ExternalPort int
	InternalPort int
	State string
}