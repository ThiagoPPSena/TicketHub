package queues

import (
	"bytes"
	"fmt"
	"net/http"
	"sharedPass/graphs"
	"sharedPass/vectorClock"
	"sort"
	"sync"
)
// Estrutura para a solicitação de compra
type Solicitation struct {
	Clock      vectorClock.VectorClock
	ServerID   int
	IsBuy      bool
	Routes     []graphs.Route
	ResponseCh chan bool
}
// Estrutura para a requisição de compra
type RequestBuy struct {
	DataJson     []byte
	ServerAddres string
	Port         string
	ResponseCh   chan bool
}
// Filas de solicitações
var SolicitationsQueue = make(chan *Solicitation)
// Filas de solicitações para o servidor
var RequestQueueOne = make(chan *RequestBuy)
var RequestQueueTwo = make(chan *RequestBuy)
// Mutex para controle de acesso a fila de solicitações
var mutex sync.Mutex

// Mandar a solicitação para o servidor
func sendRequest(serverAddres string, port string, jsonRoutes []byte) bool {
	resp, err := http.Post("http://"+serverAddres+":"+port+"/passages/buy", "application/json", bytes.NewBuffer(jsonRoutes))
	if err != nil {
		fmt.Println("Erro:", err)
		return false
	}
	defer resp.Body.Close()
	statusCode := resp.StatusCode
	return statusCode == 201
}
// Inicia o processamento da fila de solicitações
func StartProcessRequestQueue() {
	// Inicia a rotina para processar a fila de solicitações
	go func() {
		for request := range RequestQueueOne {
			request.ResponseCh <- sendRequest(request.ServerAddres, request.Port, request.DataJson)
		}
	}()
	go func() {
		for request := range RequestQueueTwo {
			request.ResponseCh <- sendRequest(request.ServerAddres, request.Port, request.DataJson)
		}
	}()
}
// Processa a fila de solicitações de compra e rollback de passagens 
func processQueue() {
	var purchaseQueue []*Solicitation

	go func() {
		for request := range SolicitationsQueue {
			mutex.Lock() // Trava o mutex antes de adicionar a solicitação à fila
			purchaseQueue = append(purchaseQueue, request) 
			sort.Slice(purchaseQueue, func(i, j int) bool {
				clockI := purchaseQueue[i].Clock
				clockJ := purchaseQueue[j].Clock
				// Variável para verificar se todos os indices do vetor de relógio vetorial são menores
				isLess := false 
				// Variável para verificar se todos os indices do vetor de relógio vetorial são maiores
				isGreater := false 
				// Faz a comparação de cada indice do vetor de relógio vetorial
				// Para saber se um é menor que o outro ou se um é maior que o outro
				for k := range clockI {
					if clockI[k] < clockJ[k] {
						isLess = true
					} else if clockI[k] > clockJ[k] {
						isGreater = true
					}
				}
				// Verifica a partir das comparações feitas anteriormente se um é menor que o outro
				if isLess && !isGreater {
					return true
				} else if !isLess && isGreater {
					return false
				}

				return purchaseQueue[i].ServerID < purchaseQueue[j].ServerID
			})
			// fmt.Println("Fila de solicitações: ", len(purchaseQueue))
			// for position, solicitation := range purchaseQueue {
			// 	fmt.Println("[", position, "]", "ServerID:", solicitation.ServerID, "Clock:", solicitation.Clock)
			// }
			mutex.Unlock() // Destrava o mutex após adicionar a solicitação à fila
		}
	}()

	for {
		if len(purchaseQueue) > 0 {
			mutex.Lock() // Trava o mutex antes de processar a fila
			nextRequest := purchaseQueue[0]
			routes := nextRequest.Routes
			isBuy := nextRequest.IsBuy

			response := false
			// Se for uma solicitação de compra, tenta comprar os assentos
			// Se for uma solicitação de rollback, tenta restaurar os assentos
			if isBuy {
				response = graphs.BuySeats(routes)
			} else {
				response = graphs.RollBack(routes)
			}

			nextRequest.ResponseCh <- response // Envia a resposta da solicitação

			graphs.SaveSeats()
			graphs.ReadRoutes()

			purchaseQueue = purchaseQueue[1:]
			mutex.Unlock() // Destrava o mutex após processar a fila
		}
	}
}

func StartProcessQueue() {
	go processQueue()
}
