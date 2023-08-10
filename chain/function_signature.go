package chain

import "github.com/ethereum/go-ethereum/crypto"

var (
	RequestSent               = []byte("RequestSent(bytes32,address)")
	RequestSentSignature      = crypto.Keccak256Hash(RequestSent).Hex()
	RequestFulfilled          = []byte("RequestFulfilled(bytes32,bytes,bytes)")
	RequestFulfilledSignature = crypto.Keccak256Hash(RequestFulfilled).Hex()

	OracleRequest          = []byte("OracleRequest(bytes32,address,address,bytes32,address,bytes)")
	OracleRequestSignature = crypto.Keccak256Hash(OracleRequest).Hex()
)
