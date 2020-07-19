package wallet

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"log"

	"golang.org/x/crypto/ripemd160"
)

const (
	checksumLength = 4
	version        = byte(0x00)
)

// Wallet : struct to handle generate private and public key
type Wallet struct {
	PrivateKey ecdsa.PrivateKey
	PublicKey  []byte
}

// NewKeyPair returns generated private and public key
func NewKeyPair() (ecdsa.PrivateKey, []byte) {
	curve := elliptic.P256()

	private, err := ecdsa.GenerateKey(curve, rand.Reader)
	if err != nil {
		log.Panic(err)
	}

	pubKey := append(private.PublicKey.X.Bytes(), private.PublicKey.Y.Bytes()...)

	return *private, pubKey
}

// MakeWallet creates a new wallet
func MakeWallet() *Wallet {
	private, pubkey := NewKeyPair()

	wallet := Wallet{private, pubkey}

	return &wallet
}

// PublicKeyHash creates public key Hash from public key created using ecdsa
// and uses ripemd160 to create public key Hash of length 160
func PublicKeyHash(pubKey []byte) []byte {
	pubHash := sha256.Sum256(pubKey)

	hasher := ripemd160.New()
	_, err := hasher.Write(pubHash[:])
	if err != nil {
		log.Panic(err)
	}

	publicRipMD := hasher.Sum(nil)

	return publicRipMD
}

// CheckSum applies sha256 twice for PublicKeyHash to generate a address from the
// public key. Returns checksum of length checksumLength
func CheckSum(payload []byte) []byte {
	firstHash := sha256.Sum256(payload)
	secondHash := sha256.Sum256(firstHash[:])

	return secondHash[:checksumLength]
}

// Address returns the address generated for the publickey.
// 1. First a public key is generated from the private key.
// 2. Ripemd160 is applied to the private key to generate the publickeyhash
// 3. sha256 is applied twice to publickeyhash and encode using Base58Encode
// 4. publickeyhash, version, Base58Encode is combine to create a address
func (w Wallet) Address() []byte {
	pubHash := PublicKeyHash(w.PublicKey)

	versionedHash := append([]byte{version}, pubHash...)
	checksum := CheckSum(versionedHash)

	fullHash := append(versionedHash, checksum...)
	address := Base58Encode(fullHash)

	return address
}
