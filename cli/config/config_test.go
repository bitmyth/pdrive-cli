package config

import "testing"

func TestConfigHosts(t *testing.T) {
	config, err := NewConfig()
	if err != nil {
		t.Error(err)
	}
	strings, s := config.Hosts()
	t.Log(strings, s)
}

func TestDefaultHost(t *testing.T) {
	config, err := NewConfig()
	if err != nil {
		t.Error(err)
	}
	strings, s := config.DefaultHost()
	t.Log(strings, s)
}
