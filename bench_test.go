package jwt

import (
	"crypto/ecdsa"
	"crypto/rsa"
	"fmt"
	"testing"
	"time"
)

var benchClaims = &Claims{
	Registered: Registered{
		Issuer: "benchmark",
		Issued: NewNumericTime(time.Now()),
	},
}

func BenchmarkECDSA(b *testing.B) {
	tests := []struct {
		key *ecdsa.PrivateKey
		alg string
	}{
		{testKeyEC256, ES256},
		{testKeyEC384, ES384},
		{testKeyEC521, ES512},
	}
	for _, test := range tests {
		b.Run("sign-"+test.alg, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_, err := benchClaims.ECDSASign(test.alg, test.key)
				if err != nil {
					b.Fatal(err)
				}
			}
		})
	}

	for _, test := range tests {
		token, err := benchClaims.ECDSASign(test.alg, test.key)
		if err != nil {
			b.Fatal(err)
		}

		b.Run("check-"+test.alg, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_, err := ECDSACheck(token, &test.key.PublicKey)
				if err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}

func BenchmarkHMAC(b *testing.B) {
	// 512-bit key
	secret := make([]byte, 64)

	// all supported algorithms in ascending order
	algs := []string{HS256, HS384, HS512}

	for _, alg := range algs {
		b.Run("sign-"+alg, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_, err := benchClaims.HMACSign(alg, secret)
				if err != nil {
					b.Fatal(err)
				}
			}
		})
	}

	for _, alg := range algs {
		token, err := benchClaims.HMACSign(alg, secret)
		if err != nil {
			b.Fatal(err)
		}

		b.Run("check-"+alg, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_, err := HMACCheck(token, secret)
				if err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}

func BenchmarkRSA(b *testing.B) {
	keys := []*rsa.PrivateKey{testKeyRSA1024, testKeyRSA2048, testKeyRSA4096}

	for _, key := range keys {
		b.Run(fmt.Sprintf("sign-%d-bit", key.Size()*8), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_, err := benchClaims.RSASign(RS384, key)
				if err != nil {
					b.Fatal(err)
				}
			}
		})
	}

	for _, key := range keys {
		token, err := benchClaims.RSASign(RS384, key)
		if err != nil {
			b.Fatal(err)
		}

		b.Run(fmt.Sprintf("check-%d-bit", key.Size()*8), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_, err := RSACheck(token, &key.PublicKey)
				if err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}
