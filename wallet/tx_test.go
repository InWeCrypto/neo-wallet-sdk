package wallet

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strconv"
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
	assert.Equal(t, (ClaimTransaction), byte(0x02))
	assert.Equal(t, (ContractTransaction), byte(0x80))
}

func TestSign(t *testing.T) {

	client := neogo.NewClient(cnf.GetString("testnode", "xxxxx") + "/extend")

	key, err := KeyFromWIF(cnf.GetString("wallet", "xxxxx"))

	assert.NoError(t, err)

	utxos, err := client.GetBalance(key.Address, "0xc56f33fc6ecfcd0c225c4ab356fee59390af8560be0e930faebe74a6daff7c9b")

	assert.NoError(t, err)

	printResult(utxos)

	tx, err := CreateSendAssertTx(
		"0xc56f33fc6ecfcd0c225c4ab356fee59390af8560be0e930faebe74a6daff7c9b",
		key.Address,
		key.Address,
		1, utxos)

	assert.NoError(t, err)

	rawtx, txid, err := tx.GenerateWithSign(key)

	assert.NoError(t, err)

	logger.DebugF("txid %s with raw data:\n%s", txid, hex.EncodeToString(rawtx))

	client = neogo.NewClient(cnf.GetString("testnode", "xxxxx"))

	status, err := client.SendRawTransaction(rawtx)

	assert.NoError(t, err)

	println(status)
}

func TestGetClaim(t *testing.T) {
	client := neogo.NewClient(cnf.GetString("testnode", "xxxxx") + "/extend")

	key, err := KeyFromWIF(cnf.GetString("wallet", "xxxxx"))

	assert.NoError(t, err)

	claims, err := client.GetClaim(key.Address)

	assert.NoError(t, err)

	printResult(claims)
}

func TestClaim(t *testing.T) {
	client := neogo.NewClient(cnf.GetString("testnode", "xxxxx") + "/extend")

	key, err := KeyFromWIF(cnf.GetString("wallet", "xxxxx"))

	assert.NoError(t, err)

	claims, err := client.GetClaim(key.Address)

	assert.NoError(t, err)

	printResult(claims)

	val, err := strconv.ParseFloat(claims.Available, 8)

	assert.NoError(t, err)

	tx, err := CreateClaimTx(val, key.Address, claims.Claims)

	assert.NoError(t, err)

	rawtx, txid, err := tx.GenerateWithSign(key)

	assert.NoError(t, err)

	logger.DebugF("txid %s with raw data:\n%s", txid, hex.EncodeToString(rawtx))

	client = neogo.NewClient(cnf.GetString("testnode", "xxxxx"))

	// rawtx, _ = hex.DecodeString("02000a889c1b256da418f238562c17d409eb4954f3c7d5da66b18862f15cb359ca51b20000d74ed730ba02d2fe80c020b39b57e76c3fbe667242ba0367050251245db7d67c00001ed2b9d7faf54bd635ddab6b40c5d5502c711cd3cecca36dedf9dd1d0d3b109a000064c73796d6ad5b73842a15ecd95e2899a174d2b28bd52013ee53952892bb7c9e0000d75408e3069905c478a6e51da2e01c054d033253711d2b5f95519576b9ca50fc00002c18e4ef9ae145ed4ed3feaa8d4d2b47f1e0c5b56a0fee07154b00709e9222830000fcc1a57a937233e870662cec33746704850bd851e05f600a4ba845790cef1d1700005a437b1b17050ddc35987e0f9511f17a492dce594a7ce600e03a340a8d7c24d40000e7179fe0a9fde5b017d898f174af5fae3feb16dca76327befbfc3d2357cb4b780000409b3aa156f8ad3bc1edc07aa68a9e56469dabf0f2bec60b2f069647934a53510000000001e72d286979ee6cb1b7e65dfddfb2e384100b8d148e7758de42e4168b71792c608e6e0800000000004263d1f1b124778d66d847801fe7cb73dd4bef500141409340333c08fa204d3b6cbf62bbf0fa8bd8a5cbeb5986c2f1a19eb3cf800fbf2d6ae21539a4c1d46f417dafadffc1454c851af685abf0e8c35e9c295c8c79d9e923210398b8d209365a197311d1b288424eaea556f6235f5730598dede5647f6a11d99aac")

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
