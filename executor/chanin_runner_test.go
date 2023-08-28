package executor

import (
	"testing"
	"time"
)

func TestChainHandler_GetWorkScore(t *testing.T) {
	t.Skip()
	chainHandler := NewChainHandler(nil, "", time.Second)
	chainHandler.WorkScoreURL = "http://210.73.218.174:8001/api/v1/worker/ranking/"
	chainHandler.GetWorkScore("9")
}
