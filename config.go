package main

import (
	"github.com/kelseyhightower/envconfig"
)

type Configuration struct {
	RancherAccessKey string `required:"true" envconfig:"rancher_access_key"`
	RancherSecretKey string `required:"true" envconfig:"rancher_secret_key"`
	RancherURL       string `required:"true" envconfig:"rancher_url"`
	ServerAPI        string `required:"true" envconfig:"server_api"`
	RancherProjectID string `required:"true" envconfig:"rancher_project_id"`
	RunInterval      int    `required:"true" envconfig:"run_interval"`
}

func LoadConfig() interface{} {
	var s Configuration
	err := envconfig.Process("abattoir", &s)
	if err != nil {
		AbattoirLog.Fatal(err.Error())
	}
	return &s
}
