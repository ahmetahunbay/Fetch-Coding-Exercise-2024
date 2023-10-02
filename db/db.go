package db

import (
	"Fetch-Coding-Exercise2024/structs"
	"container/heap"
	"database/sql"
	"log"
	"time"

	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
)

// hardcoded time format
var timeLayout = "2006-01-02T15:04:05Z"
var db *sql.DB

/*
initializes 3 tables:

userBalance stores the balance and its int value

payers stores a list of unique payers -- this is tracked because 0 balance payers are relevant

transactions stores transaction info with a uuid key for hashing
*/
func InitDB() {
	var err error
	db, err = sql.Open("sqlite3", "FetchExerciseData.db")
	if err != nil {
		log.Fatal(err)
	}

	//creates userBalance table
	_, err = db.Exec(`
        CREATE TABLE IF NOT EXISTS userBalance (
            name TEXT PRIMARY KEY,
            value INTEGER
        )
    `)
	if err != nil {
		log.Fatal(err)
	}

	var balanceExists bool
	err = db.QueryRow(`
		SELECT 1 FROM userBalance WHERE name = 'balance'
	`).Scan(&balanceExists)
	if err != nil && err != sql.ErrNoRows {
		log.Fatal(err)
	}

	if !balanceExists {
		_, err = db.Exec(`
			INSERT INTO userBalance (name, value)
			VALUES ('balance', 0)
		`)
		if err != nil {
			log.Fatal(err)
		}
	}

	//creates payers table
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS payers (
			name TEXT PRIMARY KEY
		)
	`)
	if err != nil {
		log.Fatal(err)
	}

	//creates transactions table
	_, err = db.Exec(`
        CREATE TABLE IF NOT EXISTS transactions (
            id TEXT PRIMARY KEY,
            payer TEXT,
            points INTEGER,
            timestamp TEXT
        )
    `)
	if err != nil {
		log.Fatal(err)
	}

}

/*
Inserts transaction into db if positive

adds to payer list and updates balance regardless
*/
func InsertTransaction(payer string, points int, timestamp string) (string, error) {

	err := AddPayer(payer)
	if err != nil {
		return "", err
	}

	err = UpdateBalance(points)
	if err != nil {
		return "", err
	}

	var id string
	if points > 0 {
		id = uuid.New().String()
		_, err = db.Exec(`
			INSERT INTO transactions (id, payer, points, timestamp)
			VALUES (?, ?, ?, ?)
		`, id, payer, points, timestamp)

	}

	return id, err
}

/*
removes transaction
*/
func RemoveTransaction(id string) error {
	_, err := db.Exec(`
        DELETE FROM transactions
        WHERE id = ?
    `, id)
	return err
}

/*
Queries map of all payers from db
*/
func GetPayers(payers *map[string]int) error {
	rows, err := db.Query(`
		SELECT name
		FROM payers
	`)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var payer string
		err = rows.Scan(&payer)
		if err != nil {
			return err
		}

		(*payers)[payer] = 0
	}
	if err := rows.Err(); err != nil {
		return err
	}

	return nil
}

/*
Updates payers and pq from memory
*/
func GetPayersTransactionPQ(payers *map[string]int, pq *structs.TransactionPQ) error {
	rows, err := db.Query(`
		SELECT id, payer, points, timestamp
		FROM transactions
	`)
	if err != nil {
		return err
	}
	defer rows.Close()

	transactions := structs.TransactionPQ{}
	for rows.Next() {
		var transaction structs.DBTransaction
		err = rows.Scan(&transaction.ID, &transaction.Payer, &transaction.Points, &transaction.Timestamp)
		if err != nil {
			return err
		}

		//already checked for error
		timeStamp, _ := time.Parse(timeLayout, transaction.Timestamp)
		heap.Push(&transactions, structs.Transaction{
			ID:        transaction.ID,
			Payer:     transaction.Payer,
			Points:    transaction.Points,
			Timestamp: timeStamp,
		})
		(*payers)[transaction.Payer] += transaction.Points

	}
	if err := rows.Err(); err != nil {
		return err
	}

	return nil
}

/*
Updates transaction points from ID
*/
func UpdateTransaction(id string, points int) error {
	_, err := db.Exec("UPDATE transactions SET points = ? WHERE id = ?", points, id)
	return err
}

/*
Updates balance
*/
func UpdateBalance(balanceChange int) error {
	currBalance, err := GetBalance()
	if err != nil {
		return err
	}

	newBalance := currBalance + balanceChange

	_, err = db.Exec("UPDATE userBalance SET value = ? WHERE name = \"balance\"", newBalance)
	return err
}

/*
Queries balance
*/
func GetBalance() (int, error) {
	var balance int
	err := db.QueryRow("SELECT value FROM userBalance WHERE name = \"balance\"").Scan(&balance)
	if err != nil {
		return 0, err
	}
	return balance, nil
}

/*
Adds payer to payer list, ignores unique constraint errors
*/
func AddPayer(payer string) error {
	_, err := db.Exec(`
        INSERT INTO payers (name)
        VALUES (?)
    `, payer)
	if err != nil && err.Error() == "UNIQUE constraint failed: payers.name" {
		return nil
	}
	return err
}

/*
Clears all tables
*/
func ClearDB() error {
	_, err := db.Exec(`
        DELETE FROM transactions
    `)
	if err != nil {
		return err
	}

	_, err = db.Exec("UPDATE userBalance SET value = ?", 0)
	if err != nil {
		return err
	}

	_, err = db.Exec(`
        DELETE FROM payers
    `)
	if err != nil {
		return err
	}
	return nil
}
