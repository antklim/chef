package cli

type Project interface {
	Add() error
	Bootstrap() error
	Location() (string, error)
	Name() string
}
