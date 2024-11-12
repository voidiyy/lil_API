package config

import (
	"fmt"
	"testing"
)

func TestInitConfig(t *testing.T) {
	path := "../postgres.yaml"

	conf := InitConfig(path)

	fmt.Println("abort", conf.HostName)
}
