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

	wallet2, err := neomobile.FromMnemonic("材 软 浅 毕 价 玉 售 均 司 彩 允 蒙 骤 沈 挂 肠 斯 露 慰 季 港 函 肃")

	assert.NoError(t, err)

	assert.Equal(t, wallet.Address(), wallet2.Address())
}
