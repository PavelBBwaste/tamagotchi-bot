package pet

import (
	"database/sql"
	"log"

	_ "modernc.org/sqlite"
)

var db *sql.DB

type Pet struct {
	Hunger    int
	Happiness int
	Coins     int
}

func InitDB(dataSourceName string) {
	var err error
	db, err = sql.Open("sqlite", dataSourceName)
	if err != nil {
		log.Fatal(err)
	}

	createTable()
}

func createTable() {
	sqlStmt := `
    CREATE TABLE IF NOT EXISTS pets (
        user_id INTEGER PRIMARY KEY,
        hunger INTEGER DEFAULT 100,
        happiness INTEGER DEFAULT 100,
        coins INTEGER DEFAULT 0
    );
    `
	_, err := db.Exec(sqlStmt)
	if err != nil {
		log.Fatalf("%q: %s\n", err, sqlStmt)
	}
}

func SavePetState(userID int64, pet *Pet) error {
	stmt, err := db.Prepare("INSERT OR REPLACE INTO pets(user_id, hunger, happiness, coins) VALUES (?, ?, ?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(userID, pet.Hunger, pet.Happiness, pet.Coins)
	return err
}

func LoadPetState(userID int64) (*Pet, error) {
	pet := &Pet{}
	err := db.QueryRow("SELECT hunger, happiness, coins FROM pets WHERE user_id = ?", userID).Scan(&pet.Hunger, &pet.Happiness, &pet.Coins)
	if err != nil {
		if err == sql.ErrNoRows {
			return &Pet{Hunger: 100, Happiness: 100, Coins: 0}, nil
		}
		return nil, err
	}
	return pet, nil
}
