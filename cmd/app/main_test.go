package main

import (
	"errors"
	"flag"
	"testing"

	"github.com/example/subscriptions-service/internal/config"
)

func TestMainErrorPath(t *testing.T) {
	oldLoad := configLoad
	oldFlags := flag.CommandLine
	t.Cleanup(func() {
		configLoad = oldLoad
		flag.CommandLine = oldFlags
	})

	flag.CommandLine = flag.NewFlagSet("test", flag.ContinueOnError)
	configLoad = func(string) (*config.Config, error) {
		return nil, errors.New("boom")
	}

	main()
}
