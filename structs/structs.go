package structs

import (
	"time"
)

type Transaction struct {
	Payer     string
	Points    int
	Timestamp time.Time
}

type TransactionPQ []*Transaction

func (pq TransactionPQ) Len() int { return len(pq) }

func (pq TransactionPQ) Less(i, j int) bool {
	return pq[i].Timestamp.Before(pq[j].Timestamp)
}

func (pq TransactionPQ) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
}

func (pq *TransactionPQ) Push(x interface{}) {
	item := x.(*Transaction)
	*pq = append(*pq, item)
}

func (pq *TransactionPQ) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	*pq = old[0 : n-1]
	return item
}
