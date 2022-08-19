package engine

import (
	"transaction-matching-engine/pool"
)

type MatchEngine struct {
	pools pool.MatchPool
}

func (meg *MatchEngine) Run() {

}
