package cli

import (
	"errors"
	"fmt"
)

var (
	ErrMissingDefault     = errors.New("invalid configuration: missing default profile")
	ErrProfileNotFound    = errors.New("could not find profiles configuration")
	ErrInvalidProfilePath = errors.New("could not load profiles from path")
	ErrNoActiveProfile    = errors.New("invalid configuration: missing active profile")
	ErrIncorrectVersion   = fmt.Errorf("invalid configuration version: current version is %s", ProfileVersion)
	ErrNoProfileDirectory = errors.New("could not find configuration directory")
	ErrInvalidEndpoint    = errors.New("valid endpoint is required")
	ErrCannotSetPath      = errors.New("cannot set path on loaded profile")
	ErrDoNotOverwrite     = errors.New("not overwriting existing configuration")
)
