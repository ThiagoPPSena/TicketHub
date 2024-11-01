package vectorClock

// Tipando o vetor de relógio vetorial
type VectorClock []int

var LocalClock VectorClock //Clock do servidor

var ServerId int //Id do servidor

// Função para criar um novo relógio vetorial
func NewVectorClock(numServers int) {
	LocalClock = VectorClock(make([]int, numServers))
}

func SetServerId(id int) {
	ServerId = id
}

// Função para incrementar o relógio vetorial
func (vc VectorClock) Increment() {
	vc[ServerId]++
}

// Função para atualizar o relógio vetorial com base em outro relógio vetorial
func (vc VectorClock) Update(otherVC VectorClock) {
	for i := range vc {
		if otherVC[i] > vc[i] {
			vc[i] = otherVC[i]
		}
	}
}
