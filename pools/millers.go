package pools

import (
	"fmt"
	"time"

	"github.com/vertcoin-project/one-click-miner-vnext/util"
)

var _ Pool = &Millers{}

type Millers struct {
	Address           string
	LastFetchedPayout time.Time
	LastPayout        uint64
}

func NewMillers(addr string) *Millers {
	return &Millers{Address: addr}
}

func (p *Millers) GetPendingPayout() uint64 {
	jsonPayload := map[string]interface{}{}
	err := util.GetJson(fmt.Sprintf("https://millerspool.com:4001/api/pools/vtc1/miners/%s", p.Address), &jsonPayload)
	if err != nil {
		return 0
	}
	vtc, ok := jsonPayload["pendingBalance"].(float64)
	if !ok {
		return 0
	}
	vtc *= 100000000
	return uint64(vtc)
}

func (p *Millers) GetStratumUrl() string {
	return "stratum+tcp://millerspool.com:3042"
}

func (p *Millers) GetUsername() string {
	return p.Address
}

func (p *Millers) GetPassword() string {
	return "x"
}

func (p *Millers) GetID() int {
	return 0
}

func (p *Millers) GetName() string {
	return "MillersPool.com"
}

func (p *Millers) GetFee() float64 {
	return 1.5
}

func (p *Millers) OpenBrowserPayoutInfo(addr string) {
	util.OpenBrowser(fmt.Sprintf("https://millerspool.com/?#vtc1/dashboard?address=%s", addr))
}
