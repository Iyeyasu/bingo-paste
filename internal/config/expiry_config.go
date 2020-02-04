package config

import (
	"time"

	"github.com/Iyeyasu/bingo-paste/internal/util/log"
)

// ExpiryConfig contains configuration for paste expiration.
type ExpiryConfig struct {
	Enabled      bool            `yaml:"enabled"`
	Durations    []time.Duration `yaml:"-"`
	RawDurations []int           `yaml:"durations"`
}

// DefaultExpiryConfig creates a new ExpiryConfig with default values.
func DefaultExpiryConfig() ExpiryConfig {
	return ExpiryConfig{
		Enabled:      true,
		RawDurations: []int{10, 60, 1440, 10080, 43200, 525600},
	}
}

func newDurations(values []int) []time.Duration {
	durations := make([]time.Duration, len(values), len(values))
	for i := 0; i < len(values); i++ {
		durations[i] = time.Duration(values[i]) * time.Minute
	}

	log.Debugf("Using %d expiry durations", len(durations))
	log.Debugf("Used expiry durations are %v", durations)
	return durations
}
