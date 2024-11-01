package collections

import (
	"sharedPass/graphs"
	"sharedPass/vectorClock"
)

type Body struct {
	Routes   []graphs.Route           `json:"routes"`
	Clock    *vectorClock.VectorClock `json:"clock"`
	ServerId *int                     `json:"serverId"`
}
