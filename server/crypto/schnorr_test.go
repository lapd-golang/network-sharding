package crypto

import (
	"fmt"
	"testing"
)

func TestSchnorr(t *testing.T) {
	privateKey := curve.Scalar().Pick(curve.RandomStream())

	fmt.Println(privateKey)

	publicKey := curve.Point().Mul(privateKey, curve.Point().Base())
	fmt.Println(publicKey)
	message := "Secret message to be signed."
	signature := *Sign(message, privateKey)

	derivedPublicKey := PublicKey(message, signature)
	if !derivedPublicKey.Equal(publicKey) {
		t.Errorf("Derived public key is incorrect.")
	}
	if !Verify(message, signature, publicKey) {
		t.Errorf("Signature verification is incorrect. Signature does not yield correct public key.")
	}

	fakePublicKey := curve.Point().Mul(curve.Scalar().Neg(curve.Scalar().One()), publicKey)
	if Verify(message, signature, fakePublicKey) {
		t.Errorf("Signature verification is incorrect. Signature claims to belong to an incorrect public key.")
	}
}
