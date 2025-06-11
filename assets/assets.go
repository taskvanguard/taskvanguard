package assets

import (
	"embed"
)

//go:embed templates/*
var Templates embed.FS

func Load(name string) (string, error) {
	data, err := Templates.ReadFile("templates/" + name)
	if err != nil {
		return "", err
	}
	return string(data), nil
}