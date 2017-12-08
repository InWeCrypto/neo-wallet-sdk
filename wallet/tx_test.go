package wallet

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"testing"

	"github.com/dynamicgo/config"
	"github.com/inwecrypto/neogo"
	"github.com/stretchr/testify/assert"
)

var cnf *config.Config

func init() {
	cnf, _ = config.NewFromFile("./test.json")
}

func TestType(t *testing.T) {
	assert.Equal(t, ClaimTransaction, byte(0x02))
	assert.Equal(t, ContractTransaction, byte(0x80))
}

func TestCreateWallet(t *testing.T) {

	orderurl := cnf.GetString("ordernode", "xxxxx")
	deviceid := cnf.GetString("deviceid", "xxxxx")

	key, err := KeyFromWIF(cnf.GetString("wallet", "xxxxx"))

	assert.NoError(t, err)

	resp, err := http.Post(
		fmt.Sprintf("%s/wallet/%s/%s", orderurl, deviceid, key.Address),
		"application/json",
		strings.NewReader("{}"))

	if assert.NoError(t, err) {
		result, _ := ioutil.ReadAll(resp.Body)
		fmt.Println(string(result))
		assert.Equal(t, 200, resp.StatusCode)
	}

}

func TestSign(t *testing.T) {

	client := neogo.NewClient(cnf.GetString("testnode", "xxxxx") + "/extend")

	// key, err := KeyFromWIF(cnf.GetString("wallet", "xxxxx"))

	privateKey, _ := hex.DecodeString("4473bf11d103deee68ca3349b0c6e1cf4e5da6ad64e5faa719ea78c77b4321f5")

	key, err := KeyFromPrivateKey(privateKey)

	assert.NoError(t, err)

	utxos, err := client.GetBalance(key.Address, "0xc56f33fc6ecfcd0c225c4ab356fee59390af8560be0e930faebe74a6daff7c9b")

	assert.NoError(t, err)

	printResult(utxos)

	tx, err := CreateSendAssertTx(
		"0xc56f33fc6ecfcd0c225c4ab356fee59390af8560be0e930faebe74a6daff7c9b",
		key.Address,
		"Ab8vffxvjaA3JKm3weBg6ChmZMSvorMoBM",
		1, utxos)

	assert.NoError(t, err)

	rawtx, txid, err := tx.GenerateWithSign(key)

	assert.NoError(t, err)

	logger.DebugF("txid %s with raw data:\n%s", txid, hex.EncodeToString(rawtx))

	// order, err := json.Marshal(&struct {
	// 	Tx    string
	// 	From  string
	// 	To    string
	// 	Asset string
	// 	Value string
	// }{
	// 	Tx:    fmt.Sprintf("0x%s", txid),
	// 	From:  key.Address,
	// 	To:    key.Address,
	// 	Asset: "0xc56f33fc6ecfcd0c225c4ab356fee59390af8560be0e930faebe74a6daff7c9b",
	// 	Value: "1",
	// })

	// fmt.Println(string(order))

	// assert.NoError(t, err)

	// orderurl := cnf.GetString("ordernode", "xxxxx")

	// // orderurl = "http://localhost:8000"

	// resp, err := http.Post(fmt.Sprintf("%s/order", orderurl), "application/json", bytes.NewReader(order))

	// if assert.NoError(t, err) {
	// 	assert.Equal(t, 200, resp.StatusCode)
	// } else {
	// 	return
	// }

	client = neogo.NewClient(cnf.GetString("testnode", "xxxxx"))

	//rawtx ,_ = hex.DecodeString("80000002c930bcec3b692760ba942f35f72bfe0c93952c2fc46c623073fa7019e5f360d70000f68270f4677b752f42af2f1187ff07670fc4b2ad6d80f44d53638fc9f70d33790100029b7cffdaa674beae0f930ebe6085af9093e5fe56b34a5c220ccdcf6efc336fc500c2eb0b000000007ee3a05ea28c7949b5f23d61c1fae05b754aec8c9b7cffdaa674beae0f930ebe6085af9093e5fe56b34a5c220ccdcf6efc336fc50065cd1d000000007ee3a05ea28c7949b5f23d61c1fae05b754aec8c01414075ba614beddb941a4b5c6c14ce8b61c6bc0f9824c07a0344d02ae7e1ce915445212ebf00f563c19b3d143c8cd5b27f1a079b3781fe0e2ab4a2a519049fbb42ff232103aa0047673b0bf10f936bb45a909bc70eeef396de934429c796ad496d94911820ac")

	status, err := client.SendRawTransaction(rawtx)

	assert.NoError(t, err)

	println(status)
}

func TestGetClaim(t *testing.T) {

	_, err := ReadKeyStore([]byte(`{"address":"Ab8vffxvjaA3JKm3weBg6ChmZMSvorMoBM","crypto":{"cipher":"aes-128-ctr","ciphertext":"c076f3f2d30dd5a0a384e8f3b4503ffa214958707be7fb163a49f0db9ce2ffa5","cipherparams":{"iv":"281ec882c08bbf3cbc6e0994e0e55aef"},"kdf":"scrypt","kdfparams":{"dklen":32,"n":4096,"p":6,"r":8,"salt":"3318a816a722528968412a150363cacea19e8da8820a5d961b66276ac7f3a362"},"mac":"16913516dacd0c54e81d3d5555a50ee87261bc97a115d2fa9b101fbae3ab15c2"},"id":"5cf8ec1a-646f-4051-8b7b-a2d5dfb6edb2","version":3}`), "LEIwenting0411")

	assert.NoError(t, err)

	// client := neogo.NewClient(cnf.GetString("testnode", "xxxxx") + "/extend")

	// privateKey, _ := hex.DecodeString("4473bf11d103deee68ca3349b0c6e1cf4e5da6ad64e5faa719ea78c77b4321f5")

	// key, err := KeyFromPrivateKey(privateKey)

	// println(hex.EncodeToString(key.PrivateKey.D.Bytes()))

	// // key, err := KeyFromWIF(cnf.GetString("wallet", "xxxxx"))

	// assert.NoError(t, err)

	// claims, err := client.GetClaim(key.Address)

	// assert.NoError(t, err)

	// printResult(claims)
}

type testSorter []*neogo.UTXO

func (s testSorter) Len() int      { return len(s) }
func (s testSorter) Swap(i, j int) { s[i], s[j] = s[j], s[i] }

func (s testSorter) Less(i, j int) bool {

	vi, _ := strconv.ParseFloat(s[i].Gas, 64)
	vj, _ := strconv.ParseFloat(s[j].Gas, 64)

	return vi < vj
}

func TestClaim(t *testing.T) {
	client := neogo.NewClient(cnf.GetString("testnode", "xxxxx") + "/extend")

	privateKey, _ := hex.DecodeString("4473bf11d103deee68ca3349b0c6e1cf4e5da6ad64e5faa719ea78c77b4321f5")

	key, err := KeyFromPrivateKey(privateKey)

	// key, err := KeyFromWIF(cnf.GetString("wallet", "xxxxx"))

	assert.NoError(t, err)

	claims, err := client.GetClaim(key.Address)

	assert.NoError(t, err)

	sort.Sort(claimSorter(claims.Claims))
	// claims.Available = "0.00000056"

	printResult(claims)

	val, err := strconv.ParseFloat(claims.Available, 8)

	assert.NoError(t, err)

	tx, err := CreateClaimTx(val, key.Address, claims.Claims)

	assert.NoError(t, err)

	rawtx, txid, err := tx.GenerateWithSign(key)

	assert.NoError(t, err)

	logger.DebugF("txid %s with raw data:\n%s", txid, hex.EncodeToString(rawtx))

	client = neogo.NewClient(cnf.GetString("testnode", "xxxxx"))

	// rawtx, _ = hex.DecodeString("0200058a8feef3d13ce6af23638938648a1ff581b9575fb0931b039040e9215870031b0000569b5bea2bde8a7bc61bd1aa48ec6ae9f961428427521d1e9742bd417fbc67960000ed9e2e7e0d670d6458c0bd87b268fb88087fac9419e8a4055701b07f75473a5000003aef6410399174c6d649b29aebcb8980eb9e1e2109dee80cccc16e0c8d3099700000206b4a3116fb6fdc9034b92eb6fa16f87d71088c0edbae4e263e5425b1c9a2290000000001e72d286979ee6cb1b7e65dfddfb2e384100b8d148e7758de42e4168b71792c6038cf000000000000d465a5718eed11a162ce008009bbbb20ac20d1190141406959a94cadce967e97bedb0280114afd097ad403941e8f148ed7705a495099460a37a06c1f3cdcd7b597c7cd14c662664658a3f4efe229426f60ed6f00846180232102bf23f2a852547dbe4dae698da8e49980da41a8bb3c8353c60aa6decbe4b637d9ac")

	status, err := client.SendRawTransaction(rawtx)

	assert.NoError(t, err)

	println(status)
}

func printResult(result interface{}) {

	data, _ := json.MarshalIndent(result, "", "\t")

	fmt.Println(string(data))
}

func TestDecodeAddress(t *testing.T) {
	address, err := decodeAddress("AMpupnF6QweQXLfCtF4dR45FDdKbTXkLsr")

	assert.NoError(t, err)

	logger.Debug(hex.EncodeToString(address))
}
