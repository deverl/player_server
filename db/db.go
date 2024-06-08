package db

import (
	"database/sql"
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"time"

	"player_server/utils"

	_ "github.com/go-sql-driver/mysql"
	_ "modernc.org/sqlite"
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

var db *sql.DB = nil

func GetDB() *sql.DB {
	return db
}

func init() {
	useMysql := utils.GetEnvVarAsBool("USE_MYSQL", false)

	if useMysql {
		fmt.Println("INFO: Using mysql")
	} else {
		fmt.Println("INFO: Using sqlite")
	}

	err := connect(useMysql) // Will set db if successful.

	if db == nil || err != nil {
		waitTime := 1
		maxWaitTime := 60
		for (db == nil || err != nil) && waitTime < maxWaitTime {
			fmt.Printf("ERROR: Database connection not ready... retrying in %d seconds\n", waitTime)
			time.Sleep(time.Duration(waitTime) * time.Second)
			waitTime *= 2
			err = connect(useMysql)
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

func connect(useMysql bool) error {
	var err error

	DSN := os.Getenv("DSN")

	if useMysql {
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
	} else {
		err := utils.EnsureDirectory("data")
		if err != nil {
			log.Fatal("ERROR: Couldn't create data directory")
		}
		dbFilePath := "./data/player.db"
		db, err = sql.Open("sqlite", dbFilePath)
		if err != nil {
			log.Fatal("ERROR: Couldn't open sqlite database")
		}
	}

	fmt.Println("INFO: Calling db.Ping()")

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

func InsertPlayer(player *Player) error {
	stmt := `
	INSERT INTO players (
		playerId,
		birthYear,
		birthMonth,
		birthDay,
		birthCountry,
		birthState,
		birthCity,
		deathYear,
		deathMonth,
		deathDay,
		deathCountry,
		deathState,
		deathCity,
		nameFirst,
	    nameLast,
		nameGiven,
		weight,
		height,
		bats,
		throws,
		debut,
		finalGame,
		retroID,
		bbrefID
	)
	VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	_, err := db.Exec(stmt, player.PlayerID, player.BirthYear, player.BirthMonth,
		player.BirthDay, player.BirthCountry, player.BirthState, player.BirthCity,
		player.DeathYear, player.DeathMonth, player.DeathDay, player.DeathCountry,
		player.DeathState, player.DeathCity, player.NameFirst, player.NameLast, player.NameGiven,
		player.Weight, player.Height, player.Bats, player.Throws, player.Debut, player.FinalGame,
		player.RetroID, player.BbrefID)
	if err != nil {
		fmt.Println("ERROR: Can't insert player. err:", err)
		return err
	}

	return nil
}

func UpdatePlayer(player *Player) error {
	stmt := `
	UPDATE players
		SET playerId = ?,
			birthYear = ?,
			birthMonth = ?,
			birthDay = ?,
			birthCountry = ?,
			birthState = ?,
			birthCity = ?,
			deathYear = ?,
			deathMonth = ?,
			deathDay = ?,
			deathCountry = ?,
			deathState = ?,
			deathCity = ?,
			nameFirst = ?,
			nameLast = ?,
			nameGiven = ?,
			weight = ?,
			height = ?,
			bats = ?,
			throws = ?,
			debut = ?,
			finalGame = ?,
			retroID = ?,
			bbrefID = ?
	WHERE playerId = ?`

	_, err := db.Exec(stmt, player.PlayerID, player.BirthYear, player.BirthMonth,
		player.BirthDay, player.BirthCountry, player.BirthState, player.BirthCity,
		player.DeathYear, player.DeathMonth, player.DeathDay, player.DeathCountry,
		player.DeathState, player.DeathCity, player.NameFirst, player.NameLast, player.NameGiven,
		player.Weight, player.Height, player.Bats, player.Throws, player.Debut, player.FinalGame,
		player.RetroID, player.BbrefID, player.PlayerID)
	if err != nil {
		fmt.Println("ERROR: Can't update player. err:", err)
		return err
	}

	return nil
}

func DoesPlayerExist(id string) (bool, error) {
	q := "SELECT playerId FROM players WHERE playerId = ?"
	err := db.QueryRow(q, id).Scan(&id)
	if err == sql.ErrNoRows {
		return false, nil
	}
	if err != nil {
		fmt.Println("ERROR: Can't test if player exists. err:", err)
		return false, err
	}
	return true, nil
}

func DeletePlayer(id string) error {
	q := "DELETE from players where playerId = ?"
	_, err := db.Exec(q, id)
	return err
}

func FetchPlayer(id string) (*Player, error) {
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

func GetFileHash() (string, error) {
	hash := ""
	q := "SELECT fileHash from config where id = 1"
	row := db.QueryRow(q)
	err := row.Scan(&hash)
	if err != nil {
		return "", err
	}
	return hash, nil
}

func UpdateFileHash(hash string) error {
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

func PopulatePlayer(path string) error {
	if db == nil {
		return errors.New("datbase is not connected")
	}

	dir, err := os.Getwd()
	if err != nil {
		fmt.Println("ERROR: Couldn't get working directory. err:", err)
	} else {
		fmt.Printf("INFO: Working directory: '%s'\n", dir)
	}

	hash, err := utils.GetFileHash(path)
	if err != nil {
		fmt.Println("ERROR: Couldn't calculate file hash. err:", err)
		return err
	}

	prevHash, err := GetFileHash()
	if err != nil {
		prevHash = ""
	}

	if hash == prevHash {
		fmt.Println("INFO: File hash has not changed. Not updating.")
		return nil
	}

	err = UpdateFileHash(hash)
	if err != nil {
		fmt.Println("ERROR: Could not update file hash. err:", err)
		return err
	}

	fmt.Print("INFO: Populating player database...")
	startTime := time.Now().Unix()

	f, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
		return err
	}
	defer f.Close()

	line := 0

	r := csv.NewReader(f)
	for {
		line++
		rec, err := r.Read()
		if line == 1 {
			continue
		}
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		player := playerFromCsvLine(rec)

		exists, err := DoesPlayerExist(player.PlayerID)
		if err != nil {
			fmt.Println("ERROR: DoesPlayerExist failed with err:", err)
			continue
		}

		if exists {
			err = UpdatePlayer(&player)
			if err != nil {
				fmt.Println("ERROR: UpdatePlayer failed with err:", err)
				continue
			}
		} else {
			err = InsertPlayer(&player)
			if err != nil {
				fmt.Println("ERROR: InsertPlayer failed with err:", err)
			}
		}
	}
	endTime := time.Now().Unix()
	elapsedTime := endTime - startTime

	fmt.Printf(" done in %d seconds\n.", elapsedTime)

	return nil
}

func intFromString(s string, defaultValue int) int {
	n, err := strconv.Atoi(s)
	if err != nil {
		fmt.Printf("WARNING: Using default int value for '%s'\n", s)
		n = defaultValue
	}
	return n
}

func dateFromString(s string, defaultDate time.Time) time.Time {
	t, err := time.Parse("2006-01-02", s)
	if err != nil {
		fmt.Printf("WARNING: Using default date value for '%s'\n", s)
		t = defaultDate
	}
	return t
}

func playerFromCsvLine(rec []string) Player {
	thirtyYears := 30 * 365 * 24 * time.Hour
	fiveYears := 5 * 365 * 24 * time.Hour

	defaultDebut := time.Now().Add(-thirtyYears)
	defaultFinal := time.Now().Add(-fiveYears)

	player := Player{
		PlayerID: rec[0],
	}

	if rec[1] != "" {
		player.BirthYear = intFromString(rec[1], 1970)
	}
	if rec[2] != "" {
		player.BirthMonth = intFromString(rec[2], 7)
	}
	if rec[3] != "" {
		player.BirthDay = intFromString(rec[3], 11)
	}
	if rec[4] != "" {
		player.BirthCountry = rec[4]
	}
	if rec[5] != "" {
		player.BirthState = rec[5]
	}
	if rec[6] != "" {
		player.BirthCity = rec[6]
	}
	if rec[7] != "" {
		player.DeathYear = intFromString(rec[7], 1998)
	}
	if rec[8] != "" {
		player.DeathMonth = intFromString(rec[8], 4)
	}
	if rec[9] != "" {
		player.DeathDay = intFromString(rec[9], 21)
	}
	if rec[10] != "" {
		player.DeathCountry = rec[10]
	}
	if rec[11] != "" {
		player.DeathState = rec[11]
	}
	if rec[12] != "" {
		player.DeathCity = rec[12]
	}
	if rec[13] != "" {
		player.NameFirst = rec[13]
	}
	if rec[14] != "" {
		player.NameLast = rec[14]
	}
	if rec[15] != "" {
		player.NameGiven = rec[15]
	}
	if rec[16] != "" {
		player.Weight = intFromString(rec[16], 180)
	}
	if rec[17] != "" {
		player.Height = intFromString(rec[17], 75)
	}
	if rec[18] != "" {
		player.Bats = rec[18]
	}
	if rec[19] != "" {
		player.Throws = rec[19]
	}
	if rec[20] != "" {
		player.Debut = dateFromString(rec[20], defaultDebut)
	}
	if rec[21] != "" {
		player.FinalGame = dateFromString(rec[21], defaultFinal)
	}
	if rec[22] != "" {
		player.RetroID = rec[22]
	}
	if rec[23] != "" {
		player.BbrefID = rec[23]
	}

	return player
}
