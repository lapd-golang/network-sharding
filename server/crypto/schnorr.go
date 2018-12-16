package crypto

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"log"

	"gopkg.in/dedis/kyber.v2"
	"gopkg.in/dedis/kyber.v2/group/edwards25519"
)

//PubKey ...
type PubKey struct {
	Key *[]byte
}

var curve = edwards25519.NewBlakeSHA256Ed25519()
var sha256 = curve.Hash()

//Curve ...
var Curve = curve

//Signature ...
type Signature struct {
	R kyber.Point
	S kyber.Scalar
}

//MarshalBinary ...
func (S *Signature) MarshalBinary() ([]byte, error) {
	r, err := S.R.MarshalBinary()
	if err != nil {
		return nil, err
	}
	s, err := S.S.MarshalBinary()
	if err != nil {
		return nil, err
	}

	binary := append(r, s...)

	return binary, nil
}

//UnmarshalBinary ...
func (S *Signature) UnmarshalBinary(data []byte) error {
	n := len(data)

	S.R.UnmarshalBinary(data[:n/2])
	S.S.UnmarshalBinary(data[n/2:])

	return nil
}

func (S Signature) String() string {
	return fmt.Sprintf("%s%s", S.R, S.S)
}

//Suite ...
func Suite() *edwards25519.SuiteEd25519 {
	return curve
}

//Hash ...
func Hash(s string) []byte {
	sha256.Reset()
	sha256.Write([]byte(s))

	return sha256.Sum(nil)
}

//Sign ...
// m: Message
// x: Private key
func Sign(m string, x kyber.Scalar) *Signature {
	if m == "" {
		panic("Error: Signing an empty message is insecure.")
	}

	if x.Equal(curve.Scalar().Zero()) {
		panic("Error: Private key cannot be zero.")
	}

	// Get the base of the curve.
	g := curve.Point().Base()

	// Pick a random k from allowed set.
	k := curve.Scalar().Pick(curve.RandomStream())

	// r = k * G (likewise, r = g^k)
	r := curve.Point().Mul(k, g)

	// Hash(m || r)
	e := curve.Scalar().SetBytes(Hash(m + r.String()))

	// s = k - e * x
	s := curve.Scalar().Sub(k, curve.Scalar().Mul(e, x))

	return &Signature{R: r, S: s}
}

//PublicKey ...
// m: Message
// S: Signature
func PublicKey(m string, S Signature) kyber.Point {
	if m == "" {
		panic("Error: Recovering an empty string is insecure.")
	}
	if S.R == nil || S.S == nil || S.R.Equal(curve.Point()) || S.S.Equal(curve.Scalar()) {
		panic("Error: Signature is malformed.")
	}

	// Create a generator.
	g := curve.Point().Base()

	// e = Hash(m || r)
	e := curve.Scalar().SetBytes(Hash(m + S.R.String()))

	// y = (r - s * G) * (1 / e)
	y := curve.Point().Sub(S.R, curve.Point().Mul(S.S, g))
	y = curve.Point().Mul(curve.Scalar().Div(curve.Scalar().One(), e), y)

	return y
}

//Verify ..
// m: Message
// s: Signature
// y: Public key
func Verify(m string, S Signature, y kyber.Point) bool {
	if m == "" {
		panic("Error: Signing an empty string is insecure.")
	}
	if S.R == nil || S.S == nil || S.R.Equal(curve.Point()) || S.S.Equal(curve.Scalar()) {
		panic("Error: Signature is malformed.")
	}
	if y.Equal(curve.Point().Null()) {
		panic("Error: Public key should not be the curve's neutral element.")
	}

	// Create a generator.
	g := curve.Point().Base()

	// e = Hash(m || r)
	e := curve.Scalar().SetBytes(Hash(m + S.R.String()))

	// Attempt to reconstruct 's * G' with a provided signature; s * G = r - e * y
	sGv := curve.Point().Sub(S.R, curve.Point().Mul(e, y))

	// Construct the actual 's * G'
	sG := curve.Point().Mul(S.S, g)

	// Equality check; ensure signature and public key outputs to s * G.
	return sG.Equal(sGv)
}

//EncodePubKey ..
// m: Message
// s: Signature
func EncodePubKey(s *Signature) []byte {
	var network bytes.Buffer        // Stand-in for a network connection
	enc := gob.NewEncoder(&network) // Will write to network.

	err := enc.Encode(s)
	if err != nil {
		log.Fatal("encode error:", err)
	}

	return network.Bytes()
}

//DecodePubKey ..
// b: Buffered Signature
func DecodePubKey(b []byte) *Signature {
	decoder := bytes.NewBuffer(b)  // Stand-in for a network connection
	dec := gob.NewDecoder(decoder) // Will read from network.

	var dsign = CreateSignature()
	err := dec.Decode(dsign)
	if err != nil {
		log.Fatal("decode error:", err)
	}

	return dsign
}

//CreateSignature ...
func CreateSignature() *Signature {
	privateKey := Curve.Scalar().Pick(Curve.RandomStream())
	message := "Secret message"

	signature := Sign(message, privateKey)

	return signature
}
