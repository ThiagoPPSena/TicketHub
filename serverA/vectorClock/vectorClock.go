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

// Função para obter uma copia do relógio vetorial
func (vc *VectorClock) GetClock() []int {
	clockCopy := make([]int, len(vc.clock))
	copy(clockCopy, vc.clock)
	return clockCopy
}

// Função para comparar dois relógios vetoriais, 1 = concorrencia, 0 não tem concorrencia
func (vc *VectorClock) Compare(otherVC *VectorClock) int {
	if len(vc.clock) != len(otherVC.clock) {
		panic("Os vetores devem ter o mesmo comprimento para serem comparados")
	}

	isGreater := false
	isLess := false

	for i := range vc.clock {
		if vc.clock[i] < otherVC.clock[i] {
			isLess = true
		} else if vc.clock[i] > otherVC.clock[i] {
			isGreater = true
		}
	}

	// Se ambos isGreater e isLess são verdadeiros, é concorrente (1).
	if isGreater && isLess {
		return 1
	}

	// Caso contrário, há uma relação de causalidade (0).
	return 0
}
