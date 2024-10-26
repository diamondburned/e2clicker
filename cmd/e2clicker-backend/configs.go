package main

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"reflect"

	"github.com/samber/do/v2"
	"libdb.so/e2clicker/services/api"
	"libdb.so/e2clicker/services/storage/postgresql"
)

type parseAndRegisterConfigFunc func(i do.Injector, blob json.RawMessage) error

func parseAndRegisterConfig[T any](p ...func(do.Injector)) parseAndRegisterConfigFunc {
	return func(i do.Injector, blob json.RawMessage) error {
		logger := do.MustInvoke[*slog.Logger](i)

		rt := reflect.TypeFor[T]()
		rv := reflect.New(rt.Elem()).Interface()
		if err := json.Unmarshal(blob, rv); err != nil {
			return fmt.Errorf("cannot unmarshal config: %w", err)
		}

		logger.Info(
			"registering config",
			"type", rt)

		do.ProvideValue(i, rv.(T))
		for _, f := range p {
			f(i)
		}

		return nil
	}
}

var knownConfigs = map[string]parseAndRegisterConfigFunc{
	"api":        parseAndRegisterConfig[*api.ServerConfig](),
	"postgresql": parseAndRegisterConfig[*postgresql.Config](),
}

func parseFile(i do.Injector, file string) error {
	cfgs := make(map[string]json.RawMessage)
	if err := readJSONFile(file, &cfgs); err != nil {
		return err
	}

	for name, blob := range cfgs {
		add, ok := knownConfigs[name]
		if !ok {
			return fmt.Errorf("unknown config %q", name)
		}

		if err := add(i, blob); err != nil {
			return fmt.Errorf("cannot register config %q: %w", name, err)
		}
	}

	return nil
}

func readJSONFile(path string, dst any) error {
	b, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("cannot read file: %w", err)
	}

	if err := json.Unmarshal(b, dst); err != nil {
		return fmt.Errorf("cannot unmarshal config: %w", err)
	}

	return nil
}
