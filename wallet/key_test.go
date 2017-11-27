package wallet

import (
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNEOAddress(t *testing.T) {
	key, err := KeyFromWIF("L4Ns4Uh4WegsHxgDG49hohAYxuhj41hhxG6owjjTWg95GSrRRbLL")

	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t,
		hex.EncodeToString(toBytes(key.PrivateKey)),
		"d59208b9228bff23009a666262a800f20f9dad38b0d9291f445215a0d4542beb")

	assert.Equal(t, hex.EncodeToString(publicKeyToBytes(&key.PrivateKey.PublicKey)), "0398b8d209365a197311d1b288424eaea556f6235f5730598dede5647f6a11d99a")
	assert.Equal(t, key.Address, "AMpupnF6QweQXLfCtF4dR45FDdKbTXkLsr")

	ks, err := WriteLightScryptKeyStore(key, "test")

	assert.NoError(t, err)

	key2, err := ReadKeyStore(ks, "test")

	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t,
		hex.EncodeToString(toBytes(key2.PrivateKey)),
		"d59208b9228bff23009a666262a800f20f9dad38b0d9291f445215a0d4542beb")

	assert.Equal(t, hex.EncodeToString(publicKeyToBytes(&key2.PrivateKey.PublicKey)), "0398b8d209365a197311d1b288424eaea556f6235f5730598dede5647f6a11d99a")
	assert.Equal(t, key2.Address, "AMpupnF6QweQXLfCtF4dR45FDdKbTXkLsr")

}
