package db

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"player_server/utils"

	goSql "github.com/go-sql-driver/mysql"
)

type Player struct {
	PlayerID     string    `json:"playerID,omitempty"`
	BirthYear    int       `json:"birthYear,omitempty"`
	BirthMonth   int       `json:"birthMonth,omitempty"`
	BirthDay     int       `json:"BirthDay,omitempty"`
	BirthCountry string    `json:"birthCountry,omitempty"`
	BirthState   string    `json:"birthState,omitempty"`
	BirthCity    string    `json:"birthCity,omitempty"`
	DeathYear    int       `json:"deathYear,omitempty"`
	DeathMonth   int       `json:"deathMonth,omitempty"`
	DeathDay     int       `json:"deathDay,omitempty"`
	DeathCountry string    `json:"deathCountry,omitempty"`
	DeathState   string    `json:"deathState,omitempty"`
	DeathCity    string    `json:"deathCity,omitempty"`
	NameFirst    string    `json:"nameFirst,omitempty"`
	NameLast     string    `json:"nameLast,omitempty"`
	NameGiven    string    `json:"nameGiven,omitempty"`
	Weight       int       `json:"weight,omitempty"`
	Height       int       `json:"height,omitempty"`
	Bats         string    `json:"bats,omitempty"`
	Throws       string    `json:"throws,omitempty"`
	Debut        time.Time `json:"debut,omitempty"`
	FinalGame    time.Time `json:"finalGame,omitempty"`
	RetroID      string    `json:"retroID,omitempty"`
	BbrefID      string    `json:"bbrefID,omitempty"`
}

var (
	db        *sql.DB = nil
	isOffline bool    = false
)

func GetDB() *sql.DB {
	return db
}

func init() {
	err := connect() // Will set db if successful.

	if db == nil || err != nil {
		waitTime := 1
		maxWaitTime := 60
		for (db == nil || err != nil) && waitTime < maxWaitTime {
			fmt.Printf("ERROR: Database connection not ready... retrying in %d seconds\n", waitTime)
			time.Sleep(time.Duration(waitTime) * time.Second)
			waitTime *= 2
			err = connect()
		}
	}

	if db == nil {
		log.Fatal("ERROR: Could not connect to the database")
	}

	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)

	createTables()
}

func connect() error {
	var err error

	DSN := os.Getenv("DSN")
	if DSN == "" {
		username := "rest_api_user"
		password := "rest_api_pw"
		databaseName := "rest_server"
		DSN = fmt.Sprintf("%s:%s@/%s?parseTime=true", username, password, databaseName)
	}
	db, err = sql.Open("mysql", DSN)
	if err != nil {
		fmt.Println("INFO: Database not connected. err:", err)
	}

	return db.Ping()
}

func Close() {
	db.Close()
}

func createTables() {
	createStatements := []string{
		`CREATE TABLE IF NOT EXISTS players (
		 playerID varchar (60) primary key unique,
		 birthYear int,
		 birthMonth int,
		 birthDay int,
		 birthCountry varchar(60),
		 birthState varchar(60),
		 birthCity varchar(60),
		 deathYear int,
		 deathMonth int,
		 deathDay int,
		 deathCountry varchar(60),
		 deathState varchar(60),
		 deathCity varchar(60),
		 nameFirst varchar(60),
		 nameLast varchar(60),
		 nameGiven varchar(60),
		 weight int,
		 height int,
		 bats varchar(60),
		 throws varchar(60),
		 debut DATE,
		 finalGame DATE,
		 retroID varchar(60),
		 bbrefID varchar(60)
	    )`,
		`CREATE TABLE IF NOT EXISTS config (
		id int primary key unique,
		fileHash varchar(255)
		)`,
	}

	for _, stmt := range createStatements {
		_, err := db.Exec(stmt)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func FetchPlayer(id string) (*Player, error) {
	if isOffline {
		return nil, errors.New("database is offline. Please try again later")
	}

	q := "SELECT * from players where playerId = ?"
	rows, err := db.Query(q, id)
	if err != nil {
		fmt.Println("ERROR: Can't fetch player. err:", err)
		return nil, err
	}
	defer rows.Close()
	if rows.Next() {
		player := Player{}
		err := rows.Scan(&player.PlayerID, &player.BirthYear, &player.BirthMonth,
			&player.BirthDay, &player.BirthCountry, &player.BirthState, &player.BirthCity,
			&player.DeathYear, &player.DeathMonth, &player.DeathDay, &player.DeathCountry,
			&player.DeathState, &player.DeathCity, &player.NameFirst, &player.NameLast,
			&player.NameGiven, &player.Weight, &player.Height, &player.Bats, &player.Throws,
			&player.Debut, &player.FinalGame, &player.RetroID, &player.BbrefID)
		if err != nil {
			return nil, err
		}
		return &player, nil
	}
	return nil, nil
}

func FetchPlayers(page int, pageSize int) ([]*Player, error) {
	if isOffline {
		return nil, errors.New("database is offline. Please try again later")
	}

	q := "SELECT * from players"
	if page > 0 {
		offset := (page - 1) * pageSize
		q += fmt.Sprintf(" LIMIT %d, %d", offset, pageSize)
	}
	fmt.Printf("INFO: FetchPlayers: q = '%s'\n", q)
	rows, err := db.Query(q)
	if err != nil {
		fmt.Println("ERROR: Can't fetch players. err:", err)
		return nil, err
	}
	defer rows.Close()
	players := []*Player{}
	for rows.Next() {
		player := Player{}
		err := rows.Scan(&player.PlayerID, &player.BirthYear, &player.BirthMonth,
			&player.BirthDay, &player.BirthCountry, &player.BirthState, &player.BirthCity,
			&player.DeathYear, &player.DeathMonth, &player.DeathDay, &player.DeathCountry,
			&player.DeathState, &player.DeathCity, &player.NameFirst, &player.NameLast,
			&player.NameGiven, &player.Weight, &player.Height, &player.Bats, &player.Throws,
			&player.Debut, &player.FinalGame, &player.RetroID, &player.BbrefID)
		if err != nil {
			return nil, err
		}

		players = append(players, &player)
	}
	return players, nil
}

func getFileHash() (string, error) {
	hash := ""
	q := "SELECT fileHash from config where id = 1"
	row := db.QueryRow(q)
	err := row.Scan(&hash)
	if err != nil {
		return "", err
	}
	return hash, nil
}

func updateFileHash(hash string) error {
	q := "UPDATE config set fileHash = ? WHERE id = 1"
	res, err := db.Exec(q, hash)
	if err != nil {
		return err
	}
	numRows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if numRows == 0 {
		q = "INSERT INTO config (id, fileHash) VALUES (1, ?)"
		_, err = db.Exec(q, hash)
		if err != nil {
			return err
		}
	}
	return nil
}

func dropTable(tableName string) error {
	if tableName == "" {
		return errors.New("invalid argument. tableName must be provided")
	}
	q := fmt.Sprintf("DROP TABLE IF EXISTS `%s`", tableName)
	_, err := db.Exec(q)
	return err
}

func doDataLoad(path string, tableName string) error {
	if !utils.DoesFileExist(path) {
		return errors.New("path does not exist")
	}
	if tableName == "" {
		return errors.New("invalid argument. tableName must be provided")
	}

	goSql.RegisterLocalFile(path)

	q := fmt.Sprintf(`LOAD DATA LOCAL INFILE '%s' INTO TABLE %s
		FIELDS TERMINATED BY ','
		OPTIONALLY ENCLOSED BY '"'
		ESCAPED BY '"'
		LINES TERMINATED BY '\n'
		IGNORE 1 LINES`, path, tableName)

	_, err := db.Exec(q)

	return err
}

func PopulatePlayer(path string) error {
	if db == nil {
		return errors.New("datbase is not connected")
	}

	hash, err := utils.GetFileHash(path)
	if err != nil {
		fmt.Println("ERROR: Couldn't calculate file hash. err:", err)
		return err
	}

	prevHash, err := getFileHash()
	if err != nil {
		prevHash = ""
	}

	if hash != "" && hash == prevHash {
		fmt.Println("INFO: File hash has not changed. Not updating.")
		return nil
	}

	fmt.Print("INFO: Populating player database...")

	isOffline = true

	startTime := time.Now().UnixMilli()

	// Drop table players
	err = dropTable("players")
	if err != nil {
		fmt.Println("ERROR: dropTable failed. err:", err)
	}

	// Ensure that the players table exists.
	createTables()

	// Do mysql data load
	err = doDataLoad(path, "players")
	if err != nil {
		fmt.Println("ERROR: doDataLoad failed. err:", err)
	}

	// Allow requests again at this point.
	isOffline = false

	endTime := time.Now().UnixMilli()
	elapsedTime := endTime - startTime

	// Update file hash in our config
	err = updateFileHash(hash)
	if err != nil {
		fmt.Println("ERROR: Could not update file hash. err:", err)
		return err
	}

	fmt.Printf(" done in %d ms.\n", elapsedTime)

	return nil
}
