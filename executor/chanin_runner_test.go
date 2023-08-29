package executor

import (
	"testing"
	"time"
)

func TestChainHandler_GetWorkScore(t *testing.T) {
	//t.Skip()
	chainHandler := NewChainHandler(nil, "http://210.73.218.174:8001/api/v1/verifier/get_node_score/", "", time.Second)

	chainHandler.GetLocalWorkScore()
}
