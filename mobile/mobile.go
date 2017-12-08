package neomobile

import (
	"encoding/hex"
	"encoding/json"

	"github.com/inwecrypto/bip39"
	"github.com/inwecrypto/neo-wallet-sdk/wallet"
	"github.com/inwecrypto/neogo"
)

// Wallet neo mobile wallet
type Wallet struct {
	key *wallet.Key
}

// Tx neo rawtx wrapper
type Tx struct {
	Data string
	ID   string
}

// FromWIF create wallet from wif
func FromWIF(wif string) (*Wallet, error) {
	key, err := wallet.KeyFromWIF(wif)

	if err != nil {
		return nil, err
	}

	return &Wallet{
		key: key,
	}, nil
}

// New create a new wallet
func New() (*Wallet, error) {
	key, err := wallet.NewKey()

	if err != nil {
		return nil, err
	}

	return &Wallet{
		key: key,
	}, nil
}

// FromMnemonic create wallet from mnemonic
func FromMnemonic(mnemonic string) (*Wallet, error) {
	dic, _ := bip39.GetDict("zh_CN")

	data, err := bip39.MnemonicToByteArray(mnemonic, dic)

	if err != nil {
		return nil, err
	}

	data = data[1 : len(data)-1]

	println(hex.EncodeToString(data))

	key, err := wallet.KeyFromPrivateKey(data)

	if err != nil {
		return nil, err
	}

	return &Wallet{
		key: key,
	}, nil
}

// FromKeyStore create wallet from keystore
func FromKeyStore(keystore string, password string) (*Wallet, error) {
	key, err := wallet.ReadKeyStore([]byte(keystore), password)

	if err != nil {
		return nil, err
	}

	return &Wallet{
		key: key,
	}, nil
}

// ToKeyStore write wallet to keystore format string
func (wrapper *Wallet) ToKeyStore(password string) (string, error) {
	keystore, err := wallet.WriteLightScryptKeyStore(wrapper.key, password)

	return string(keystore), err
}

// CreateAssertTx create assert transfer raw tx
func (wrapper *Wallet) CreateAssertTx(assert, from, to string, amount float64, unspent string) (*Tx, error) {
	var utxos []*neogo.UTXO

	if err := json.Unmarshal([]byte(unspent), &utxos); err != nil {
		return nil, err
	}

	tx, err := wallet.CreateSendAssertTx(assert, from, to, amount, utxos)

	if err != nil {
		return nil, err
	}

	rawtxdata, txid, err := tx.GenerateWithSign(wrapper.key)

	return &Tx{
		Data: hex.EncodeToString(rawtxdata),
		ID:   txid,
	}, err
}

// Address get wallet address
func (wrapper *Wallet) Address() string {
	return wrapper.key.Address
}

// Mnemonic gete mnemonic string
func (wrapper *Wallet) Mnemonic() (string, error) {
	privateKeyBytes := wrapper.key.ToBytes()

	dic, _ := bip39.GetDict("zh_CN")

	println(hex.EncodeToString(privateKeyBytes))

	data, err := bip39.NewMnemonic(privateKeyBytes, dic)

	if err != nil {
		return "", err
	}

	return data, nil
}

// CreateClaimTx create claim tx
func (wrapper *Wallet) CreateClaimTx(amount float64, address string, unspent string) (*Tx, error) {
	var utxos []*neogo.UTXO

	if err := json.Unmarshal([]byte(unspent), &utxos); err != nil {
		return nil, err
	}

	tx, err := wallet.CreateClaimTx(amount, address, utxos)

	if err != nil {
		return nil, err
	}

	rawtxdata, txid, err := tx.GenerateWithSign(wrapper.key)

	return &Tx{
		Data: hex.EncodeToString(rawtxdata),
		ID:   txid,
	}, err
}
