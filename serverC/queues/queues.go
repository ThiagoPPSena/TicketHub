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
	contador := 25

	go func() {
		for request := range SolicitationsQueue {
			purchaseQueue = append(purchaseQueue, request) // Adiciona cada solicitação à fila slice
		}
	}()

	// Processamento das solicitações na fila ilimitada
	for {
		// time.Sleep(500 * time.Millisecond)
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
			//Printar a fila como tá o estado atual
			fmt.Printf("Tamanho da FILA: %d\n", len(purchaseQueue))
			for position,request := range purchaseQueue {
				fmt.Println("[",position,"] ", request.Clock, request.ServerID)
			}
			// fmt.Printf("Tamanho da FILA: %d\n", len(purchaseQueue))
			// fmt.Printf("Processando solicitação %v \nSERVIDOR: %d \nRELÓGIO: %v\n", nextRequest.Routes, nextRequest.ServerID, nextRequest.Clock)
			// diminui contador
			if contador > 0 {
				contador--
				nextRequest.ResponseCh <- true
			} else {
				nextRequest.ResponseCh <- false
			}

			// Remove a solicitação processada da fila
			purchaseQueue = purchaseQueue[1:]
		}
	}
}

func StartProcessQueue() {
	go processQueue()
}
