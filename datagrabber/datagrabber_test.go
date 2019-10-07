package datagrabber

import (
	"testing"

	"github.com/Barrett17/market_prices/types"
)

// Warning don't run this test often as create a lot
// of connections.
func TestConcurrency(t *testing.T) {
	for i := 0; i < 5; i++ {
		go GetData(types.EUR_USD)
		go pollData()
		go GetData(types.EUR_GBP)
	}
}
