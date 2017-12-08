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
	client := neogo.NewClient(cnf.GetString("mainnet", "xxxxx") + "/extend")

	key, err := ReadKeyStore(
		[]byte(`{"address":"Ab8vffxvjaA3JKm3weBg6ChmZMSvorMoBM","crypto":{"cipher":"aes-128-ctr","ciphertext":"c076f3f2d30dd5a0a384e8f3b4503ffa214958707be7fb163a49f0db9ce2ffa5","cipherparams":{"iv":"281ec882c08bbf3cbc6e0994e0e55aef"},"kdf":"scrypt","kdfparams":{"dklen":32,"n":4096,"p":6,"r":8,"salt":"3318a816a722528968412a150363cacea19e8da8820a5d961b66276ac7f3a362"},"mac":"16913516dacd0c54e81d3d5555a50ee87261bc97a115d2fa9b101fbae3ab15c2"},"id":"5cf8ec1a-646f-4051-8b7b-a2d5dfb6edb2","version":3}`),
		"Xiaoji123",
	)

	println(hex.EncodeToString(key.PrivateKey.D.Bytes()))

	// key, err := KeyFromWIF(cnf.GetString("wallet", "xxxxx"))

	assert.NoError(t, err)

	claims, err := client.GetClaim(key.Address)

	assert.NoError(t, err)

	printResult(claims)
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

	assert.NoError(t, err)

	claims, err := client.GetClaim(key.Address)

	assert.NoError(t, err)

	sort.Sort(testSorter(claims.Claims))

	printResult(claims)

	val, err := strconv.ParseFloat(claims.Available, 8)

	assert.NoError(t, err)

	tx, err := CreateClaimTx(val, key.Address, claims.Claims)

	assert.NoError(t, err)

	rawtx, txid, err := tx.GenerateWithSign(key)

	assert.NoError(t, err)

	logger.DebugF("txid %s with raw data:\n%s", txid, hex.EncodeToString(rawtx))

	client = neogo.NewClient(cnf.GetString("testnode", "xxxxx"))

	// rawtx, _ = hex.DecodeString("02002321ae25b7673fd7c46b0fe11420fca54d4b6db19f616e5230e5537c1e4b051f740000be0c7294e1f7713721974139882955da65da8abe93fa75f453a86c7572d7183d0000922cb93933088142995d2759d738b41d9db25011ebd0730d46630ed57da9a4ab00000e2caa4469ae9dd952cddcaa113ce0411352a6856fa9a5014a8e38a997a58b3300005f5d862da91b003b4da1dfa46b2871dcc43445b68dc7a3880bb0d71ca312f6830000a6ccd237f4720159f97129622e8c10c788e3228557b2f8fcaf4420da09021f59000088667dc8c42b455875e716d84ca53f44cf2c8fab31fd1bf0b5fb304e472e2a59000085cfdd24e4321f51b104f372633ddbf6bd4eda1c8bb62f405c1ff61792e983a40000df43890143af211026d6c364da44ed7a7795b9ce3dd153f505eda1815c29e93e000075bcdb3992100c3f6b2fc1fd64090f236981a7d7a18518272fcd60a62a4cb53c0000695581b06570b92f4ec133246c262f1ef3106bd2732d8eb922ce1a51a5a7a7e700007a067190203e810f9b921d2a3c7c9311e7616f38ff3211e99d37bca3505dae380000ce5e796bc68844a4acb17f57c297c4dbf7150e4e3e7b7bb2a7e24a444754f47200006c795f8124e719fe4df97804ab637310bcc2b6d5045437a4041244feed999ef10000a465dc96ba06fdf92e654a89dd4618b2229009ec2388d43fcbb91f3bb8ee10600000ec5cf3b328107cde59061b8994e9259a4f07c9202471e7744aa3e42b9d676df00000b685602f03dfe7083dff80b3fb9c93c2e27cfad7c693fe8a7dea5720bada606e0000f9e3e6c1a2202ea84f01e6e6c456c3d39b4e8818e18d88513c430c1c53c04d4b00002dd5062a5461a9e3780823d27beba54a65832810fedb816103163525e0b34fae0000dfea20bbdec7edc1eb398474ebf05c099d957bedc6391981a9e5c127ea2fb7900000861a2e1c34078c38fa3197b6fd7a4d3a4397da33f9b5d7abfe0e4661752ba82101009f747c6b6f6083590e84fc9cec9e89f5aa9f4fea24d8715bf8fcf2aeaf0e07200000e94c5583498cd0a5ac377970c9b21ca40d2d785a9743573207882cdacf8ca70c0000e942affbbc344a22c9c74d6af159e1117bc08b228bdbf323b0bf03a1b5418289000003c1af7130c2fb222cda360450bb4efdb133c949f7cf484ecf36490cc781ba7e0000c359942f62070a16ec18411c60d03757d0526f3ac16a8c1889a6e097f6552d630000c5b9ccaa162fd327092d28abb1b31f41836a85ee87a74a2d955045b4b621051e0000e906dd6536e035b44919ea15e5604d007bdaf882f37ebf4e2405bd7af2be42e60000da07e946820d4df9eed700344b4508213f097e590f63aa3785875d512c0e3a7300000255b6a164e81fdc9d120c0ead5ee2d81edfa61ae1c76a17fda45f01939cba9b00001475c33f35ed2e1b9aba0ddd9ccb2f44843af84579777007bfe959115a7a90590000646c8732645fe822344aabd28fe8acb1f154ac883dd0f3038834758a8bffadef000029bf383d2fe71c669a2839b62c68e8529e1d87df0bd6eebd796d82a9e78fb1ea0000113129c02c969369f3f471df3193335078314d4eaea91f48bc52c61e887dcdf8000056524f338748a7eb039ff5205c0f45b2ef108ef8efddfc826e6438474dab257e0000000001e72d286979ee6cb1b7e65dfddfb2e384100b8d148e7758de42e4168b71792c60931b9801000000007ee3a05ea28c7949b5f23d61c1fae05b754aec8c0141401c2fe8a5a63574fa3acc34ace272b807e17a8f76e209c265faa6eb5a6ea78d788af1f739a5b4d02a34ee4e3547811a9680911fa74f20402936d758e50a291782232103aa0047673b0bf10f936bb45a909bc70eeef396de934429c796ad496d94911820ac")

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
