package vectorClock

import (
    "sync"
)

// Tipando o vetor de relógio vetorial
type VectorClock []int

var (
    LocalClock VectorClock // Clock do servidor
    ServerId   int         // Id do servidor
    clockMutex sync.Mutex  // Declaração do mutex
)

// Função para criar um novo relógio vetorial
func NewVectorClock(numServers int) {
    LocalClock = VectorClock(make([]int, numServers))
}

func SetServerId(id int) {
    ServerId = id
}

// Função para obter uma copia do relógio vetorial
func (vc VectorClock) Copy() VectorClock {
    copy := make(VectorClock, len(vc))
    copy.Update(vc)
    return copy
}

// Função para incrementar o relógio vetorial
func (vc VectorClock) Increment() {
    clockMutex.Lock() // Trava o mutex antes de modificar o relógio vetorial
    defer clockMutex.Unlock() // Destrava o mutex após modificar o relógio vetorial
    vc[ServerId]++
}

// Função para atualizar o relógio vetorial com base em outro relógio vetorial
func (vc VectorClock) Update(otherVC VectorClock) {
    clockMutex.Lock() // Trava o mutex antes de modificar o relógio vetorial
    defer clockMutex.Unlock() // Destrava o mutex após modificar o relógio vetorial
    for i := range vc {
        if otherVC[i] > vc[i] {
            vc[i] = otherVC[i]
        }
    }
}