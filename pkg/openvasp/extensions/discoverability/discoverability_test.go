package discoverability_test

import (
	"math/rand/v2"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	. "github.com/trisacrypto/trisa/pkg/openvasp/extensions/discoverability"
)

func TestUptime(t *testing.T) {
	t.Run("Random", func(t *testing.T) {
		for i := 0; i < 128; i++ {
			uptime := Uptime(randDuration())
			text, err := uptime.MarshalText()
			require.NoError(t, err)

			var uptime2 Uptime
			require.NoError(t, uptime2.UnmarshalText(text))
			require.Equal(t, uptime, uptime2)
		}
	})

	t.Run("Fixed", func(t *testing.T) {
		testCases := map[string]time.Duration{
			"0":       0,
			"1":       1 * time.Second,
			"30":      30 * time.Second,
			"60":      1 * time.Minute,
			"3600":    1 * time.Hour,
			"86400":   24 * time.Hour,
			"172800":  48 * time.Hour,
			"259200":  72 * time.Hour,
			"360000":  100 * time.Hour,
			"604800":  168 * time.Hour,
			"3600000": 1000 * time.Hour,
		}

		for text, duration := range testCases {
			var uptime Uptime
			require.NoError(t, uptime.UnmarshalText([]byte(text)))
			require.Equal(t, duration, time.Duration(uptime))

			marshaledText, err := uptime.MarshalText()
			require.NoError(t, err)
			require.Equal(t, text, string(marshaledText))
		}
	})

	t.Run("Error", func(t *testing.T) {
		var uptime Uptime
		err := uptime.UnmarshalText([]byte("not a number"))
		require.Error(t, err)
	})
}

func TestUptimeSince(t *testing.T) {
	start := time.Now()
	time.Sleep(200 * time.Millisecond)
	uptime := UptimeSince(start)
	require.GreaterOrEqual(t, time.Duration(uptime), 200*time.Millisecond)
}

func randDuration() time.Duration {
	dur := time.Duration(rand.Int64N(1e9)) * time.Millisecond
	return dur.Truncate(time.Second)
}
