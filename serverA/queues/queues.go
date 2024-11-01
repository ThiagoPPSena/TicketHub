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

	go func() {
		for request := range SolicitationsQueue {
			purchaseQueue = append(purchaseQueue, request) // Adiciona cada solicitação à fila slice
			fmt.Printf("Quantidade de solicitações na fila: %d\n", len(purchaseQueue))
		}
	}()

	// Processamento das solicitações na fila ilimitada
	for {
		if len(purchaseQueue) > 0 {
			// Ordena a fila slice com base no relógio vetorial
			sort.Slice(purchaseQueue, func(i, j int) bool {
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

			// Processa a solicitação mais antiga
			nextRequest := purchaseQueue[0]
			fmt.Printf("Processando solicitação %v do servidor %d com relógio %v\n", nextRequest.Routes, nextRequest.ServerID, nextRequest.Clock)
			nextRequest.ResponseCh <- true

			// Remove a solicitação processada da fila
			purchaseQueue = purchaseQueue[1:]
		}
	}
}

func StartProcessQueue() {
	go processQueue()
}
