package kconfig

import (
	"bufio"
	"os"
	"strings"
)

// Kconfig represents a kernel configuration.
type Kconfig struct {

	// map of kernel parameters
	params map[string]string
}

// New returns a new Kconfig.
func New() Kconfig {
	return Kconfig{
		params: make(map[string]string),
	}
}

func (k *Kconfig) Read(p string) error {

	// Open file p.
	f, err := os.Open(p)
	if err != nil {
		return err
	}
	defer f.Close()

	s := bufio.NewScanner(f)
	if err := parse(s, k.params); err != nil {
		return err
	}

	return nil
}

// Write writes out the Kconfig to the given path.
func (k Kconfig) Write(p string) error {

	f, err := os.Create(p)
	if err != nil {
		return err
	}
	defer f.Close()

	w := bufio.NewWriter(f)
	if err := dump(k.params, w); err != nil {
		return err
	}

	f.Sync()

	return nil
}

// Merge merges a map of configuration values into the Kconfig.
// The keys of params can be optionally prefixed with CONFIG_,
// the prefix will be trimmed.
func (k Kconfig) Merge(params map[string]string) {
	for key, val := range params {
		k.params[strings.TrimPrefix(key, "CONFIG_")] = val
	}
}
