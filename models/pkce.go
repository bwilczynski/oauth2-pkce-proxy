package models

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
)

const codeChallengeLength = 32

type CodeVerifier string

func CreateCodeVerifier() (res CodeVerifier, err error) {
	b := make([]byte, codeChallengeLength)
	_, err = rand.Read(b)
	if err != nil {
		return
	}

	res = CodeVerifier(base64.RawURLEncoding.EncodeToString(b))
	return
}

func (v *CodeVerifier) codeChallengeS256() string {
	b := sha256.Sum256([]byte(string(*v)))
	return base64.RawURLEncoding.EncodeToString(b[:])
}

func (v *CodeVerifier) Verify(challenge string) bool {
	return v.codeChallengeS256() == challenge
}
