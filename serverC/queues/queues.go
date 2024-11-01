package queues

import (
	"fmt"
	"sharedPass/graphs"
	"sharedPass/vectorClock"
	"sort"
)

type Solicitation struct {
	Clock      vectorClock.VectorClock
	ServerID   int
	Routes     []graphs.Route
	ResponseCh chan bool
}

var SolicitationsQueue = make(chan *Solicitation)

func processQueue() {
	var purchaseQueue []*Solicitation

	for request := range SolicitationsQueue {
		purchaseQueue = append(purchaseQueue, request)

		sort.Slice(purchaseQueue, func(i, j int) bool {
			// Lógica de concorrencia e casualidade
			clockI := purchaseQueue[i].Clock
			clockJ := purchaseQueue[j].Clock
			for k := range clockI {
				if clockI[k] < clockJ[k] {
					return true
				} else if clockI[k] > clockJ[k] {
					return false
				}
			}
			return purchaseQueue[i].ServerID < purchaseQueue[j].ServerID
		})

		if len(purchaseQueue) > 0 {
			nextRequest := purchaseQueue[0]
			fmt.Printf("Processando solitação %v do servidor %d com relógio %v\n", nextRequest.Routes, nextRequest.ServerID, nextRequest.Clock)
			// Lógica de efetuação de compra
			nextRequest.ResponseCh <- true

			purchaseQueue = purchaseQueue[1:]
		}
	}
}

func StartProcessQueue() {
	go processQueue()
}
