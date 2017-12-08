package neomobiletest

import (
	"testing"

	"github.com/inwecrypto/neo-wallet-sdk/mobile"
	"github.com/stretchr/testify/assert"
)

func TestMem(t *testing.T) {
	wallet, err := neomobile.New()

	assert.NoError(t, err)

	mne, err := wallet.Mnemonic()

	println(mne)

	assert.NoError(t, err)

	wallet2, err := neomobile.FromMnemonic(mne)

	assert.NoError(t, err)

	assert.Equal(t, wallet.Address(), wallet2.Address())
}
