/*
See: https://gitlab.com/OpenVASP/travel-rule-protocol/-/blob/master/extensions/discoverability.md
*/
package discoverability

import (
	"context"
	"strconv"
	"time"
)

const (
	VersionEndpoint    = "/version"
	UptimeEndpoint     = "/uptime"
	ExtensionsEndpoint = "/extensions"
)

// Client defines the service interface for interacting with the TRP discoverability endpoints.
type Client interface {
	Version(ctx context.Context, address string) (*Version, error)
	Uptime(ctx context.Context, address string) (Uptime, error)
	Extensions(ctx context.Context, address string) (*Extensions, error)
}

// Version returned on the version discoverability endpoint.
type Version struct {
	Version string `json:"version,omitempty"`
	Vendor  string `json:"vendor,omitempty"`
}

// Extensions is returned on the extensions discoverability endpoint.
type Extensions struct {
	Required  []string `json:"required,omitempty"`
	Supported []string `json:"supported,omitempty"`
}

// Uptime is returned on the uptime discoverability endpoint.
type Uptime time.Duration

// Get the uptime since the given start timestamp.
func UptimeSince(start time.Time) Uptime {
	return Uptime(time.Since(start))
}

// Marshals the uptime into a base 10 integer describing the number of seconds of the duration.
func (u Uptime) MarshalText() (text []byte, err error) {
	uptime := int64(time.Duration(u).Seconds())
	text = []byte(strconv.FormatInt(uptime, 10))
	return text, nil
}

// Unmarshals the uptime from a base 10 integer describing the number of seconds up.
func (u *Uptime) UnmarshalText(text []byte) (err error) {
	var uptime int64
	if uptime, err = strconv.ParseInt(string(text), 10, 64); err != nil {
		return err
	}

	*u = Uptime(time.Duration(uptime) * time.Second)
	return nil
}
