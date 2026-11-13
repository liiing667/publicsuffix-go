package publicsuffix_test

import (
	"testing"

	wpsl "github.com/weppos/publicsuffix-go/net/publicsuffix"
	xpsl "golang.org/x/net/publicsuffix"
)

func TestPublicSuffix(t *testing.T) {
	testCases := []string{
		"example.com",
		"www.example.com",
		"example.co.uk",
		"www.example.co.uk",
		"example.blogspot.com",
		"www.example.blogspot.com",
		"parliament.uk",
		"www.parliament.uk",
		// not listed
		"www.example.test",
	}

	for _, testCase := range testCases {
		ws, wb := wpsl.PublicSuffix(testCase)
		xs, xb := xpsl.PublicSuffix(testCase)

		if ws != xs || wb != xb {
			t.Errorf("PublicSuffix(%v): x/psl -> (%v, %v) != w/psl -> (%v, %v)", testCase, xs, xb, ws, wb)
		}
	}
}

func TestEffectiveTLDPlusOne(t *testing.T) {
	testCases := []string{
		"example.com",
		"www.example.com",
		"example.co.uk",
		"www.example.co.uk",
		"example.blogspot.com",
		"www.example.blogspot.com",
		"parliament.uk",
		"www.parliament.uk",
		// not listed
		"www.example.test",
	}

	for _, testCase := range testCases {
		ws, we := wpsl.EffectiveTLDPlusOne(testCase)
		xs, xe := xpsl.EffectiveTLDPlusOne(testCase)

		if ws != xs || we != xe {
			t.Errorf("EffectiveTLDPlusOne(%v): x/psl -> (%v, %v) != w/psl -> (%v, %v)", testCase, xs, xe, ws, we)
		}
	}
}

func TestPublicSuffixNormalization(t *testing.T) {
	testCases := []struct {
		input  string
		suffix string
		icann  bool
	}{
		// mixed case and trailing dot are normalized
		{"WWW.EXAMPLE.COM", "com", true},
		{"www.example.com.", "com", true},
		{"Example.CO.UK.", "co.uk", true},
		// U-labels are normalized to A-labels
		{"example.рф", "xn--p1ai", true},
		{"sub.example.公司.cn", "xn--55qx5d.cn", true},
		// fail closed: the domain is itself a public suffix
		{"com", "com", true},
		{"blogspot.com", "blogspot.com", false},
		// fail closed: invalid domains are returned unchanged
		{".example.com", ".example.com", false},
	}

	for _, testCase := range testCases {
		suffix, icann := wpsl.PublicSuffix(testCase.input)
		if suffix != testCase.suffix || icann != testCase.icann {
			t.Errorf("PublicSuffix(%v) = (%v, %v), want (%v, %v)", testCase.input, suffix, icann, testCase.suffix, testCase.icann)
		}
	}
}

func TestEffectiveTLDPlusOneNormalization(t *testing.T) {
	testCases := map[string]string{
		"WWW.EXAMPLE.COM":   "example.com",
		"www.example.com.":  "example.com",
		"example.рф":        "example.xn--p1ai",
		"sub.example.公司.cn": "example.xn--55qx5d.cn",
	}

	for input, want := range testCases {
		if got, err := wpsl.EffectiveTLDPlusOne(input); err != nil {
			t.Errorf("EffectiveTLDPlusOne(%v) returned error: %v", input, err)
		} else if got != want {
			t.Errorf("EffectiveTLDPlusOne(%v) = %v, want %v", input, got, want)
		}
	}
}
