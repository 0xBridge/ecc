// Package signature provides an additional abstraction and modularity to digital signature schemes of built-in implementations
package signature

import (
	"crypto"
	"io"

	"github.com/bytemare/cryptotools/signature/internal"
)

// Identifier indicates the signature scheme to be used.
type Identifier string

const (
	// Ed25519 indicates usage of the Ed25519 signature scheme.
	Ed25519 Identifier = "Ed25519"

	//
	// Ed448 Identifier = "Ed448".
)

// Signature abstracts digital signature operations, wrapping built-in implementations.
type Signature interface {

	// LoadKey loads the given key. Will not fail if the key is invalid, but it might later.
	LoadKey(privateKey []byte)

	// GenerateKey generates a fresh signing key and keeps it internally.
	GenerateKey() error

	// GetPrivateKey returns the private key.
	GetPrivateKey() []byte

	// GetPublicKey returns the public key.
	GetPublicKey() []byte

	// Public implements the Signer.Public() function.
	Public() crypto.PublicKey

	// Seed re-calculates the private key from the seed for compatible schemes. Implementations can only retain a seed
	// to reduce storage size.
	Seed(seed []byte)

	// SignMessage uses the internal private key to sign the message. The message argument doesn't need to be hashed beforehand.
	SignMessage(message ...[]byte) []byte

	// Sign implements the Signer.Sign() function.
	Sign(rand io.Reader, digest []byte, opts crypto.SignerOpts) (signature []byte, err error)

	// Verify checks whether signature of the message is valid given the public key.
	Verify(publicKey, message, signature []byte) bool
}

// New returns a Signature implementation to the specified scheme.
func (i Identifier) New() Signature {
	switch i {
	case Ed25519:
		return internal.NewEd25519()
	default:
		panic("invalid identifier")
	}
}

// Sign returns the signature of message (concatenated, if using a variadic argument) using secretKey.
func (i Identifier) Sign(secretKey []byte, message ...[]byte) []byte {
	s := i.New()
	s.LoadKey(secretKey)

	return s.SignMessage(message...)
}

// Verify checks whether signature of the message is valid given the public key.
func (i Identifier) Verify(publicKey, message, signature []byte) bool {
	return i.New().Verify(publicKey, message, signature)
}
