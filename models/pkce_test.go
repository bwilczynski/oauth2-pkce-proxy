package models

import "testing"

func TestCreateCodeVerifier(t *testing.T) {
	v1, _ := CreateCodeVerifier()
	v2, _ := CreateCodeVerifier()

	if v1 == v2 {
		t.Errorf("code verifiers doesn't appear to be random: %v == %v", v1, v2)
	}
}

func TestVerify(t *testing.T) {
	v := CodeVerifier("pdBvOtTb0nK8D7jjvx_ZvQQk8iiG89R7uWpjbXAaepk")
	valid := "O0qzciRkB7ekLKb8kTos2cH-6Rv_H_2X9Mf_c7Gvlq4"
	notValid := "whatever"

	if !v.Verify(valid) {
		t.Errorf("code challenge did not passed validation but it should, verifier: %v challenge: %v", v, valid)
	}
	if v.Verify(notValid) {
		t.Errorf("code challenge passed validation but it should not, verifier: %v challenge %v", v, notValid)
	}
}
