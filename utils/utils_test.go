package utils

import "testing"

func TestPrefixZero(t *testing.T) {
	testcases := []struct {
		name   string
		inOTP  string
		outOTP string
	}{
		{
			name:   "should prefix zero",
			inOTP:  "1234",
			outOTP: "001234",
		},
		{
			name:   "should not prefix zero",
			inOTP:  "123456",
			outOTP: "123456",
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			otp := prefix0(tc.inOTP)
			if otp != tc.outOTP {
				t.Errorf("Expected %s, got %s", tc.outOTP, otp)
			}
		})
	}

}

func TestHOTPToken(t *testing.T) {
	testcases := []struct {
		name     string
		secret   string
		interval int64
		otp      string
	}{
		{
			name:     "should generate HOTP token",
			secret:   "dummySECRETdummy",
			interval: int64(50780342),
			otp:      "971294",
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			otp := GetHOTPToken(tc.secret, tc.interval)
			if otp != tc.otp {
				t.Errorf("Expected %s, got %s", tc.otp, otp)
			}
		})
	}
}
