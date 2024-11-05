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

type Solicitation struct {
	Clock      vectorClock.VectorClock
	ServerID   int
	IsBuy      bool
	Routes     []graphs.Route
	ResponseCh chan bool
}

type RequestBuy struct {
	DataJson   []byte
	Port       string
	ResponseCh chan bool
}

var SolicitationsQueue = make(chan *Solicitation)
var RequestQueueOne = make(chan *RequestBuy)
var RequestQueueTwo = make(chan *RequestBuy)
var mutex sync.Mutex

func sendRequest(port string, jsonRoutes []byte) bool {
	resp, err := http.Post("http://localhost:"+port+"/passages/buy", "application/json", bytes.NewBuffer(jsonRoutes))
	if err != nil {
		fmt.Println("Erro:", err)
		return false
	}
	defer resp.Body.Close()
	return true
}

func StartProcessRequestQueue() {
	go func() {
		for request := range RequestQueueOne {
			request.ResponseCh <- sendRequest(request.Port, request.DataJson)
		}
	}()
	go func() {
		for request := range RequestQueueTwo {
			request.ResponseCh <- sendRequest(request.Port, request.DataJson)
		}
	}()
}

func processQueue() {
	var purchaseQueue []*Solicitation

	go func() {
		for request := range SolicitationsQueue {
			mutex.Lock()
			purchaseQueue = append(purchaseQueue, request)
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
				if isLess && !isGreater {
					return true
				} else if !isLess && isGreater {
					return false
				}

				return purchaseQueue[i].ServerID < purchaseQueue[j].ServerID
			})
			fmt.Println("Fila de solicitações: ", len(purchaseQueue))
			for position, solicitation := range purchaseQueue {
				fmt.Println("[", position, "]", "ServerID:", solicitation.ServerID, "Clock:", solicitation.Clock)
			}
			mutex.Unlock()
		}
	}()

	for {
		if len(purchaseQueue) > 0 {
			mutex.Lock()
			nextRequest := purchaseQueue[0]
			routes := nextRequest.Routes
			isBuy := nextRequest.IsBuy

			response := false
			if isBuy {
				response = graphs.BuySeats(routes)
			} else {
				response = graphs.RollBack(routes)
			}

			nextRequest.ResponseCh <- response

			graphs.SaveSeats()
			graphs.ReadRoutes()

			purchaseQueue = purchaseQueue[1:]
			mutex.Unlock()
		}
	}
}

func StartProcessQueue() {
	go processQueue()
}
