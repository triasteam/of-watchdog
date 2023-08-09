package config

import (
	"fmt"
	"github.com/ethereum/go-ethereum/crypto"
	"testing"
)

func TestChain_FuncName(t *testing.T) {
	fmt.Println(crypto.Keccak256Hash([]byte("test1")))
	fmt.Println(crypto.Keccak256Hash([]byte("test2")))
	fmt.Println(crypto.Keccak256Hash([]byte("test3")))
	//out:
	//0x6d255fc3390ee6b41191da315958b7d6a1e5b17904cc7683558f98acc57977b4
	//0x4da432f1ecd4c0ac028ebde3a3f78510a21d54087b161590a63080d33b702b8d
	//0x204558076efb2042ebc9b034aab36d85d672d8ac1fa809288da5b453a4714aae
}
