package queues

import (
    "fmt"
    "sharedPass/graphs"
    "sharedPass/vectorClock"
    "sort"
    "sync"
)

type Solicitation struct {
    Clock      vectorClock.VectorClock
    ServerID   int
    Routes     []graphs.Route
    ResponseCh chan bool
}

var SolicitationsQueue = make(chan *Solicitation)
var mutex sync.Mutex // Declaração do mutex

func processQueue() {
    var purchaseQueue []*Solicitation
    contador := 25

    go func() {
        for request := range SolicitationsQueue {
            mutex.Lock() // Trava o mutex antes de modificar a lista
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
			//Printar a fila de solicitações
			fmt.Println("Fila de solicitações: ", len(purchaseQueue))
			for position, solicitation := range purchaseQueue {
					fmt.Println("[", position, "]", "ServerID:", solicitation.ServerID, "Clock:", solicitation.Clock)
				}
            mutex.Unlock() // Destrava o mutex após modificar a lista
        }
    }()

    // Processamento das solicitações na fila ilimitada
    for {
        if len(purchaseQueue) > 0 {
            mutex.Lock() // Trava o mutex antes de acessar a lista
            // Processa a solicitação mais antiga
            nextRequest := purchaseQueue[0]
            
            if contador > 0 {
                contador--
                nextRequest.ResponseCh <- true
            } else {
                nextRequest.ResponseCh <- false
            }
            
            // Remove a solicitação processada da fila
            purchaseQueue = purchaseQueue[1:]
            mutex.Unlock() // Destrava o mutex após modificar a lista
            fmt.Println("Contador:", contador)
        }
    }
}

func StartProcessQueue() {
    go processQueue()
}