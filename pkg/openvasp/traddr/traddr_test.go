package traddr_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/trisacrypto/trisa/pkg/openvasp/traddr"
)

func TestMake(t *testing.T) {
	testCases := []struct {
		input    string
		expected string
	}{
		{
			"beneficiary.com/x/12345?t=i",
			"ta2W2HPKfHxgSgrzY178knqXHg1H3jfeQrwQ9JrKBs9wv",
		},
		{
			"localhost",
			"ta226rAKyTmMxzsmkda5U2Vhw1",
		},
		{
			"192.22.123.12",
			"ta42ZEr7rdvSpGnVvnttK1bWYQU4f4f",
		},
		{
			"https://beneficiary.com/x/12345",
			"ta2W2HPKfHxgSgrzY178knqXHg1H3jfeQrwQ9JrKBs9wv",
		},
		{
			"https://beneficiary.com/x/12345?t=i",
			"ta2W2HPKfHxgSgrzY178knqXHg1H3jfeQrwQ9JrKBs9wv",
		},
		{
			"http://192.22.123.12?t=i",
			"ta42ZEr7rdvSpGnVvnttK1bWYQU4f4f",
		},
		{
			"https://beneficiary.com/x/12345?t=i&foo=bar&color=red",
			"taGw1e4cjuujwyHBV51aspbLBUbcDhebW7ss8iF6dLEn19WfyQZt6HaUTStSE2YcadkFG",
		},
		{
			"https://trisa.beneficiary.com",
			"ta6HcFqwh5y4D3r7z66JyKqvSAekyXnXc2KwdXSUec",
		},
		{
			"https://trisa.beneficiary.com/x/12345",
			"ta4uMjxqDq4t7nefCdrm4ssCTewVbvBLzmUnVcYe4SXqGafPTPewG",
		},
		{
			"https://trisa.beneficiary.com/x/12345?t=i",
			"ta4uMjxqDq4t7nefCdrm4ssCTewVbvBLzmUnVcYe4SXqGafPTPewG",
		},
		{
			"trisa.beneficiary.com:4443",
			"ta3dgYSWEoAJPfU2Wdr2hvQmXem4dHbfKuBvkbKGfBrzdb61C",
		},
		{
			"https://beneficiary.com:8000",
			"ta21hULpYqwBirCCHHU13UMT1hhy3syAAPXJem8cJ",
		},
		{
			"trisa.beneficiary.com:4443/x/12345",
			"ta2wijJqYhcj5zriQyV5SVWybf71AyHEXgRCa7WWJ8XQk9g1NodLEHbKoHdg",
		},
		{
			"https://beneficiary.com:8000/x/12345",
			"takLQA9MSCjHr3kkHh4fnC5rWNbrySUJhTwQHZsML2anv95kkqe",
		},
		{
			"trisa.beneficiary.com:4443/x/12345?t=i",
			"ta2wijJqYhcj5zriQyV5SVWybf71AyHEXgRCa7WWJ8XQk9g1NodLEHbKoHdg",
		},
		{
			"https://beneficiary.com:8000/x/12345?t=i",
			"takLQA9MSCjHr3kkHh4fnC5rWNbrySUJhTwQHZsML2anv95kkqe",
		},
		{
			"envoy.local/x/12345?t=i",
			"taEi7uRVTrhfqLUzffcu7AqS17UJXUjMqhAdZ6p",
		},
		{
			"envoy.local:8100/x/12345?t=i",
			"ta7pwsXgswc9soxfnDgsBbWpDVsEwXy4G4hhPEN8WkrmaB",
		},
		{
			"envoy.local:8100",
			"taAFKaTJj4goXwwCKCmdNPtG8zeze5BeZiK",
		},
		{
			"http://envoy.local:8100?t=i",
			"taAFKaTJj4goXwwCKCmdNPtG8zeze5BeZiK",
		},
	}

	for i, tc := range testCases {
		actual, err := traddr.Make(tc.input)
		require.NoError(t, err, "received unexpected error on test case %d", i)
		require.Equal(t, tc.expected, actual, "equality mismatch on test case %d", i)
	}
}

func TestDecodeURL(t *testing.T) {
	testCases := []struct {
		input    string
		expected string
	}{
		{
			"ta2W2HPKfHxgSgrzY178knqXHg1H3jfeQrwQ9JrKBs9wv",
			"https://beneficiary.com/x/12345?t=i",
		},
		{
			"ta226rAKyTmMxzsmkda5U2Vhw1",
			"https://localhost?t=i",
		},
		{
			"ta42ZEr7rdvSpGnVvnttK1bWYQU4f4f",
			"https://192.22.123.12?t=i",
		},
		{
			"ta2W2HPKfHxgSgrzY178knqXHg1H3jfeQrwQ9JrKBs9wv",
			"https://beneficiary.com/x/12345?t=i",
		},
		{
			"ta2W2HPKfHxgSgrzY178knqXHg1H3jfeQrwQ9JrKBs9wv",
			"https://beneficiary.com/x/12345?t=i",
		},
		{
			"ta42ZEr7rdvSpGnVvnttK1bWYQU4f4f",
			"https://192.22.123.12?t=i",
		},
		{
			"taGw1e4cjuujwyHBV51aspbLBUbcDhebW7ss8iF6dLEn19WfyQZt6HaUTStSE2YcadkFG",
			"https://beneficiary.com/x/12345?color=red&foo=bar&t=i",
		},
		{
			"ta6HcFqwh5y4D3r7z66JyKqvSAekyXnXc2KwdXSUec",
			"https://trisa.beneficiary.com?t=i",
		},
		{
			"ta4uMjxqDq4t7nefCdrm4ssCTewVbvBLzmUnVcYe4SXqGafPTPewG",
			"https://trisa.beneficiary.com/x/12345?t=i",
		},
		{
			"ta3dgYSWEoAJPfU2Wdr2hvQmXem4dHbfKuBvkbKGfBrzdb61C",
			"https://trisa.beneficiary.com:4443?t=i",
		},
		{
			"ta21hULpYqwBirCCHHU13UMT1hhy3syAAPXJem8cJ",
			"https://beneficiary.com:8000?t=i",
		},
		{
			"ta2wijJqYhcj5zriQyV5SVWybf71AyHEXgRCa7WWJ8XQk9g1NodLEHbKoHdg",
			"https://trisa.beneficiary.com:4443/x/12345?t=i",
		},
		{
			"takLQA9MSCjHr3kkHh4fnC5rWNbrySUJhTwQHZsML2anv95kkqe",
			"https://beneficiary.com:8000/x/12345?t=i",
		},
		{
			"taEi7uRVTrhfqLUzffcu7AqS17UJXUjMqhAdZ6p",
			"https://envoy.local/x/12345?t=i",
		},
		{
			"ta7pwsXgswc9soxfnDgsBbWpDVsEwXy4G4hhPEN8WkrmaB",
			"https://envoy.local:8100/x/12345?t=i",
		},
		{
			"taAFKaTJj4goXwwCKCmdNPtG8zeze5BeZiK",
			"https://envoy.local:8100?t=i",
		},
	}

	for i, tc := range testCases {
		actual, err := traddr.DecodeURL(tc.input)
		require.NoError(t, err, "received unexpected error on test case %d", i)
		require.Equal(t, tc.expected, actual, "equality mismatch on test case %d", i)
	}
}

func TestEncoding(t *testing.T) {
	testCases := []struct {
		input    string
		expected string
		err      error
		msg      string
	}{
		{
			"beneficiary.com/x/12345?t=i",
			"ta2W2HPKfHxgSgrzY178knqXHg1H3jfeQrwQ9JrKBs9wv",
			nil,
			"happy case",
		},
		{
			"trisa.beneficiary.com:8000?t=i",
			"ta3dgYSWEoAJPfU2Wdr2hvQmXem4dHbfKyYrisbBTnXjL2Kub",
			nil,
			"happy case with port",
		},
		{
			"trisa.beneficiary.com:8000/x/12345?t=i",
			"ta2wijJqYhcj5zriQyV5SVWybf71AyHEXjdvtJEwNdQYykjqYkVT6dUXAfTf",
			nil,
			"happy case with port and path",
		},
		{
			"api.testing.21analytics.xyz/transfers/01H11RHY53MBMEQB1VJ4KJF39Z?t=i",
			"ta2SRCMhxtdtKVfoV2MTziaN7F9WD7WKyGvynQa7w3as7ziTr583ZRU1DeKarHjFrLSxuQAZLciZMqry4Yk797qKErWYnCmj8sLaH",
			nil,
			"mildred tilcott",
		},
		{
			"localhost.com/travelRule/12345?t=i",
			"taH3Md9Z9jkNtB1Es9Tvmq9gcYYQRyQkLvJe89Cp8pEmXLAaG5zpRp",
			nil,
			"localhost travel address",
		},
		{
			"10.20.1.231/x/12345?t=i",
			"ta7eawSJpSSSqZScV8b8hDD6hLMtLKC4S6tfKc3",
			nil,
			"ipv4 travel address",
		},
		{
			"envoy.local/x/12345?t=i",
			"taEi7uRVTrhfqLUzffcu7AqS17UJXUjMqhAdZ6p",
			nil,
			".local domain",
		},
		{
			"envoy.local:8100/x/12345?t=i",
			"ta7pwsXgswc9soxfnDgsBbWpDVsEwXy4G4hhPEN8WkrmaB",
			nil,
			".local domain with port",
		},
		{
			"beneficiary.com/x/12345", "",
			traddr.ErrMissingQueryString,
			"error due to missing query string \"t=i\"",
		},
		{
			"https://beneficiary.com/x/12345?t=i", "",
			traddr.ErrURIScheme,
			"error due to URI scheme presence",
		},
		{
			"beneficiary/x/12345?t=i", "",
			traddr.ErrInvalidTLD,
			"error due to unresolvable TLD",
		},
	}

	for i, tc := range testCases {
		addr, err := traddr.Encode(tc.input)
		if tc.err != nil {
			require.Error(t, err, "expected error on test case %d: %s", i, tc.msg)
			require.ErrorIs(t, err, tc.err, "error did not match in test case %d: %s", i, tc.msg)
		} else {
			require.NoError(t, err, "expected no error on test case %d: %s", i, tc.msg)
			require.Equal(t, tc.expected, addr, "equality mismatch on test case %d: %s", i, tc.msg)
		}
	}
}

func TestDecoding(t *testing.T) {
	testCases := []struct {
		input    string
		expected string
		err      error
		msg      string
	}{
		{
			"ta2W2HPKfHxgSgrzY178knqXHg1H3jfeQrwQ9JrKBs9wv",
			"beneficiary.com/x/12345?t=i",
			nil,
			"happy case",
		},
		{
			"ta3dgYSWEoAJPfU2Wdr2hvQmXem4dHbfKyYrisbBTnXjL2Kub",
			"trisa.beneficiary.com:8000?t=i",
			nil,
			"happy case with port",
		},
		{
			"ta2wijJqYhcj5zriQyV5SVWybf71AyHEXjdvtJEwNdQYykjqYkVT6dUXAfTf",
			"trisa.beneficiary.com:8000/x/12345?t=i",
			nil,
			"happy case with port and path",
		},
		{
			"ta2SRCMhxtdtKVfoV2MTziaN7F9WD7WKyGvynQa7w3as7ziTr583ZRU1DeKarHjFrLSxuQAZLciZMqry4Yk797qKErWYnCmj8sLaH",
			"api.testing.21analytics.xyz/transfers/01H11RHY53MBMEQB1VJ4KJF39Z?t=i",
			nil,
			"mildred tilcott",
		},
		{
			"taH3Md9Z9jkNtB1Es9Tvmq9gcYYQRyQkLvJe89Cp8pEmXLAaG5zpRp",
			"localhost.com/travelRule/12345?t=i",
			nil,
			"localhost travel address",
		},
		{
			"ta7eawSJpSSSqZScV8b8hDD6hLMtLKC4S6tfKc3",
			"10.20.1.231/x/12345?t=i",
			nil,
			"ipv4 travel address",
		},
		{
			"2W2HPKfHxgSgrzY178knqXHg1H3jfeQrwQ9JrKBs9wv", "",
			traddr.ErrUnhandledScheme,
			"error due to missing ta prefix",
		},
		{
			"taEJKtAQyrS5x6i59GBS2fcbcUxoR14dYiW9cZu", "",
			traddr.ErrMissingQueryString,
			"error due to missing query string \"t=i\"",
		},
		{
			"ta2BCfkBRHmbhyZuKHmEdHpypmo7HG4RJVgaWYR4LkKGLyAtJQkDtJrK", "",
			traddr.ErrURIScheme,
			"error due to URI scheme presence",
		},
		{
			"taEJKtAQyrS5x6i59GKB6iMPx1iDzs8HXGNhbk1", "",
			traddr.ErrInvalidTLD,
			"error due to unresolvable TLD",
		},
		{
			"ta123412", "",
			traddr.ErrInvalidFormat,
			"error due to invalid format",
		},
		{
			"ta3MNQE1Y", "",
			traddr.ErrChecksum,
			"error due to invalid checksum",
		},
	}

	for i, tc := range testCases {
		addr, err := traddr.Decode(tc.input)
		if tc.err != nil {
			require.Error(t, err, "expected error on test case %d: %s", i, tc.msg)
			require.ErrorIs(t, err, tc.err, "error did not match in test case %d: %s", i, tc.msg)
		} else {
			require.NoError(t, err, "expected no error on test case %d: %s", i, tc.msg)
			require.Equal(t, tc.expected, addr, "equality mismatch on test case %d: %s", i, tc.msg)
		}
	}
}

func TestValidTLD(t *testing.T) {
	t.Run("Valid", func(t *testing.T) {
		tests := []string{
			"https://beneficiary.com",
			"beneficiary.org",
			"trisa.beneficiary.co",
			"https://beneficiary.co.uk:8000",
			"beneficiary.us:8000",
			"trisa.beneficiary.net:8000",
			"testnet.trisa.beneficiary.io:8000",
			"envoy.local",
			"envoy.local:8100",
			"http://envoy.local",
			"http://envoy.local:8100",
		}

		for i, tc := range tests {
			u, _ := traddr.Parse(tc)
			require.NoError(t, u.ValidTLD(), "received unexpected error for test case %d", i)
		}
	})

	t.Run("Invalid", func(t *testing.T) {
		tests := []string{
			"https://beneficiary.bear",
			"beneficiary.bear",
			"trisa.beneficiary.bear",
			"https://beneficiary.bear:8000",
			"beneficiary.bear:8000",
			"trisa.beneficiary.bear:8000",
			"testnet.trisa.beneficiary.bear:8000",
		}

		for i, tc := range tests {
			u, _ := traddr.Parse(tc)
			require.ErrorIs(t, u.ValidTLD(), traddr.ErrInvalidTLD, "did not receive expected error for test case %d", i)
		}

	})
}
