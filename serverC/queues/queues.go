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
			sort.Slice(purchaseQueue, func(i, j int) bool {
				clockI := purchaseQueue[i].Clock
				clockJ := purchaseQueue[j].Clock

				isLess := false
				isGreater := false
				
				for k := range clockI {
					if clockI[k] < clockJ[k] {
						isLess = true
					} else if clockI[k] > clockJ[k] {
						isGreater = true
					}
				}
				// Se o relógio de i for menor que o de j, retorna true
				if isLess && !isGreater {
					return true
				} else if !isLess && isGreater {
					return false
				}
		
				return purchaseQueue[i].ServerID < purchaseQueue[j].ServerID
			})
		}
	}()

	// Processamento das solicitações na fila ilimitada
	for {
		// time.Sleep(500 * time.Millisecond)
		if len(purchaseQueue) > 0 {
			// Processa a solicitação mais antiga
			nextRequest := purchaseQueue[0]
			//Printar a fila como tá o estado atual
			fmt.Printf("Tamanho da FILA: %d\n", len(purchaseQueue))
			for position,request := range purchaseQueue {
				fmt.Println("[",position,"] ", request.Clock, request.ServerID)
			}
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
