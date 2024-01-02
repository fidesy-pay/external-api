package config

import "fmt"

var (
	ErrConfigNotFoundByKey = func(key int) error {
		return fmt.Errorf("config not found by key = %q", key)
	}
)
