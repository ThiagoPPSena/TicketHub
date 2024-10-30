package vectorClock

// Tipando o vetor de relógio vetorial
type VectorClock struct {
	clock []int
}

// Função para criar um novo relógio vetorial
func NewVectorClock(numServers int) *VectorClock {
	return &VectorClock{make([]int, numServers)}
}

// Função para incrementar o relógio vetorial
func (vc *VectorClock) Increment(serverIndex int) {
	vc.clock[serverIndex]++
}

// Função para atualizar o relógio vetorial com base em outro relógio vetorial
func (vc *VectorClock) Update(otherVC *VectorClock) {
	for i := range vc.clock {
		if otherVC.clock[i] > vc.clock[i] {
			vc.clock[i] = otherVC.clock[i]
		}
	}
}

// Função para obter o relógio vetorial
func (vc *VectorClock) GetClock() []int {
	clockCopy := make([]int, len(vc.clock))
	copy(clockCopy, vc.clock)
	return clockCopy
}

// Função para comparar dois relógios vetoriais
func (vc *VectorClock) Compare(otherVC *VectorClock) int {
	for i := range vc.clock {
		if vc.clock[i] < otherVC.clock[i] {
			return -1
		} else if vc.clock[i] > otherVC.clock[i] {
			return 1
		}
	}

	return 0
}
