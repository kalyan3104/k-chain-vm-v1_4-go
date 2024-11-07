package secp256k1

import (
	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcec/v2/ecdsa"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/kalyan3104/k-chain-vm-v1_4-go/crypto/hashing"
	"github.com/kalyan3104/k-chain-vm-v1_4-go/crypto/signing"
)

type MessageHashType uint8

const (
	ECDSAPlainMsg MessageHashType = iota
	ECDSASha256
	ECDSADoubleSha256
	ECDSAKeccak256
	ECDSARipemd160
)

type secp256k1 struct {
}

func NewSecp256k1() *secp256k1 {
	return &secp256k1{}
}

// VerifySecp256k1 checks a secp256k1 signature provided in the DER encoding format.
// The hash type used over the message can also be configured using @param hashType
func (sec *secp256k1) VerifySecp256k1(key, msg, sig []byte, hashType uint8) error {
	pubKey, err := btcec.ParsePubKey(key)
	if err != nil {
		return err
	}

	messageHash, err := sec.hashMessage(msg, hashType)
	if err != nil {
		return err
	}

	signature, err := ecdsa.ParseSignature(sig)
	if err != nil {
		return err
	}

	verified := signature.Verify(messageHash, pubKey)
	if !verified {
		return signing.ErrInvalidSignature
	}

	return nil
}

// EncodeSecp256k1DERSignature creates a DER encoding of a signature provided with r and s.
// Useful when having the plain params - like in the case of ecrecover
//
//	from ethereum
func (sec *secp256k1) EncodeSecp256k1DERSignature(r, s []byte) []byte {
	rScalar := &btcec.ModNScalar{}
	rScalar.SetByteSlice(r)

	sScalar := &btcec.ModNScalar{}
	sScalar.SetByteSlice(s)

	sig := ecdsa.NewSignature(rScalar, sScalar)

	return sig.Serialize()
}

func (sec *secp256k1) hashMessage(msg []byte, hashType uint8) ([]byte, error) {
	hasher := hashing.NewHasher()

	var err error
	var hashedMsg []byte
	switch MessageHashType(hashType) {
	case ECDSASha256:
		hashedMsg, err = hasher.Sha256(msg)
	case ECDSADoubleSha256:
		hashedMsg = chainhash.DoubleHashB(msg)
	case ECDSAKeccak256:
		hashedMsg, err = hasher.Keccak256(msg)
	case ECDSARipemd160:
		hashedMsg, err = hasher.Ripemd160(msg)
	case ECDSAPlainMsg:
		hashedMsg = msg
	default:
		return nil, signing.ErrHasherNotSupported
	}

	if err != nil {
		return nil, err
	}

	return hashedMsg, nil
}
