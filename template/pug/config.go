package pug

import "github.com/go-iris2/iris2/template/html"

// Pug is the 'jade', same configs as the html engine

// Config for pug template engine
type Config html.Config

// DefaultConfig returns the default configuration for the pug(jade) template engine
func DefaultConfig() Config {
	return Config(html.DefaultConfig())
}
