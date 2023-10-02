package endpoints

import (
	"Fetch-Coding-Exercise2024/db"
	"Fetch-Coding-Exercise2024/structs"
	"container/heap"
)

/*
This method is used to clear out the cache and db from any outstanding debt,
whether in the form of a negative balance or negative transaction.
*/
func HandleDebtBalance(debt int) error {
	//makes local var positive for simplicity
	debt = -debt

	//checks if either pq is empty or debt is gone
	for debt > 0 && len(pq) > 0 {
		//if debt<pq[0] points, subtract debt from pq[0] points
		if pq[0].Points > debt {
			pq[0].Points -= debt
			err := db.UpdateTransaction(pq[0].ID, pq[0].Points)
			payers[pq[0].Payer] -= debt
			if err != nil {
				return err
			}
			return nil
			//otherwise subtract pq[0] points from debt and pop
		} else {
			head := heap.Pop(&pq)
			transaction := head.(structs.Transaction)
			err := db.RemoveTransaction(transaction.ID)
			payers[transaction.Payer] -= transaction.Points
			if err != nil {
				return err
			}
			debt -= transaction.Points
		}
	}
	return nil
}

/*
This is very similar to the method above, but used in the /spend route

There are 2 key differences:

1. the pq will never be exhausted before the pointsSpent

2. we return a paidMap to keep track of our json encoding

*/

func HandlePointsSpent(pointsSpent int) (map[string]int, error) {
	paidMap := map[string]int{}

	for pointsSpent > 0 {
		if pq[0].Points > pointsSpent {
			pq[0].Points -= pointsSpent
			err := db.UpdateTransaction(pq[0].ID, pq[0].Points)
			if err != nil {
				return nil, err
			}
			payers[pq[0].Payer] -= pointsSpent
			_, exists := paidMap[pq[0].Payer]
			if !exists {
				paidMap[pq[0].Payer] = -pointsSpent
			} else {
				paidMap[pq[0].Payer] -= pointsSpent
			}
			return paidMap, nil
		} else {
			head := heap.Pop(&pq)
			transaction := head.(structs.Transaction)
			err := db.RemoveTransaction(transaction.ID)
			if err != nil {
				return nil, err
			}
			payers[transaction.Payer] -= transaction.Points
			_, exists := paidMap[transaction.Payer]
			if !exists {
				paidMap[transaction.Payer] = -transaction.Points
			} else {
				paidMap[transaction.Payer] -= transaction.Points
			}
			pointsSpent -= transaction.Points
		}
	}

	return paidMap, nil
}

/*
This method is called to load in the cache from the db
*/
func UpdateCache() error {
	err := db.GetPayers(&payers)
	if err != nil {
		return err
	}
	err = db.GetPayersTransactionPQ(&payers, &pq)
	if err != nil {
		return err
	}

	return nil
}
