package chain

import "github.com/ethereum/go-ethereum/crypto"

var (
	RequestSent               = []byte("RequestSent(bytes32,bytes32,address)")
	RequestSentSignature      = crypto.Keccak256Hash(RequestSent).Hex()
	RequestFulfilled          = []byte("RequestFulfilled(bytes32,address,uint,bytes,bytes)")
	RequestFulfilledSignature = crypto.Keccak256Hash(RequestFulfilled).Hex()

	OracleRequest           = []byte("OracleRequest(bytes32,address,address,bytes32,address,bytes)")
	OracleRequestSignature  = crypto.Keccak256Hash(OracleRequest).Hex()
	OracleResponse          = []byte("OracleResponse(bytes32)")
	OracleResponseSignature = crypto.Keccak256Hash(OracleResponse).Hex()

	OracleRequestTimeout          = []byte("OracleRequestTimeout(bytes32,string)")
	OracleRequestTimeoutSignature = crypto.Keccak256Hash(OracleRequestTimeout).Hex()
)
