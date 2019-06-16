package pgproxy

import (
	"database/sql"
	"fmt"
	"math/rand"
	"os"
	"testing"
	"time"

	"github.com/crunchydata/crunchy-proxy/config"
	"github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 4444
	user     = "postgres"
	password = "password"
	dbname   = "segment"
	table    = "copy_test"
)

var db *sql.DB
var cfg = DefaultConfig

func init() {
	rand.Seed(time.Now().UnixNano())

	// add origin servers provided in environment
	if pwd, ok := os.LookupEnv("PGA_ORIGIN_PASSWORD"); ok {
		cfg.OriginPassword = pwd
	}

	if host, ok := os.LookupEnv("PGA_ORIGIN_HOST"); ok {
		cfg.OriginHost = host
	}

	// TODO: remove
	config.SetConfigPath("./config.yaml")
	config.ReadConfig()
}

func TestMain(m *testing.M) {
	defer shutdown(db)

	stopCh := make(chan struct{})
	server := Server{
		Config: cfg,
	}

	if err := server.Setup(); err != nil {
		panic("Failed to setup server: " + err.Error())
	}

	go server.Start(stopCh)

	m.Run()
	go func() {
		stopCh <- struct{}{}
	}()
}

func TestTableCreate(t *testing.T) {
	db = setupServer(t)
	defer cleanup(t, db)

	t.Logf("Creating table '%s'...", table)
	query := fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s(letters varchar, numbers integer)", table)
	if rows, err := db.Query(query); err != nil {
		t.Fatal(err)
	} else {
		rows.Close()
	}

	t.Log("Created table ", table)
}

func TestProxyCopy(t *testing.T) {
	db = setupServer(t)
	defer cleanup(t, db)

	println("Creating table...")
	cols := []string{"letters varchar", "numbers integer"}
	createTable(t, db, table, cols)

	// create example data
	rows := make([][]interface{}, 26)
	for i := range rows {
		letter := string('a' + byte(i))
		num := int64(i + 1)
		t.Logf("Data: %s, %d\n", letter, num)
		rows[i] = []interface{}{letter, num}
	}

	t.Log("Start copying example data...")
	if err := copyExampleData(t, db, table, []string{"letters", "numbers"}, rows); err != nil {
		t.Fatal("Error copying example data:", err)
	}
	t.Log("Done copying example data")
}

func TestProxyMessage(t *testing.T) {
	db = setupServer(t)
	defer cleanup(t, db)

	println("Creating table...")
	cols := []string{"letters varchar"}
	createTable(t, db, table, cols)

	rows := [][]interface{}{
		{
			"8J-Ukg34DpvtQP5SR17nFWllY7Y19eAdtq_oLkoUaZw",
		},
	}
	t.Log("Start copying example data...")
	if err := copyExampleData(t, db, table, []string{"letters"}, rows); err != nil {
		t.Fatal("Error copying example data:", err)
	}
	t.Log("Done copying example data")
}

func TestProxyCopyLarge(t *testing.T) {
	db = setupServer(t)
	defer cleanup(t, db)

	cols := []string{"a varchar", "b varchar"}
	createTable(t, db, table, cols)

	// generate example data
	minChar, maxChar := int64(3585), int64(3654)
	rows := make([][]interface{}, 800)
	for r := range rows {
		row := make([]interface{}, len(cols))
		for c := range row {
			chars := make([]rune, 200)
			for i := range chars {
				char := minChar + rand.Int63n(maxChar-minChar)
				chars[i] = rune(char)
			}
			row[c] = string(chars)
		}
		rows[r] = row
	}

	t.Log("Start copying example data...")
	if err := copyExampleData(t, db, table, []string{"a", "b"}, rows); err != nil {
		t.Fatal("Error copying example data:", err)
	}
	t.Log("Done copying example data")
}

func copyExampleData(t *testing.T, db *sql.DB, tableName string, cols []string, rows [][]interface{}) error {
	t.Log("Starting transaction on table ", tableName)

	txn, err := db.Begin()
	if err != nil {
		return fmt.Errorf("couldn't begin transaction: %v", err)
	}

	t.Log("Preparing COPY FROM query...")
	copyQuery := pq.CopyIn(tableName, cols...)
	stmt, err := txn.Prepare(copyQuery)
	if err != nil {
		return fmt.Errorf("couldn't prepare query '%s': %v", copyQuery, err)
	}

	t.Log("Send data in statements...")
	for _, data := range rows {
		if _, err = stmt.Exec(data...); err != nil {
			txn.Rollback()
			return fmt.Errorf("encountered error '%v' when processing data %v", err, data)
		}
	}

	t.Log("Executing statement")
	if _, err = stmt.Exec(); err != nil {
		txn.Rollback()
		return fmt.Errorf("error executing transaction: %v", err)
	}

	t.Log("Closing statement")
	if err = stmt.Close(); err != nil {
		return fmt.Errorf("error closing statement: %v", err)
	}

	t.Log("Closing transaction")
	if err = txn.Commit(); err != nil {
		return fmt.Errorf("error closing transaction: %v", err)
	}

	return nil
}

func setupServer(t *testing.T) *sql.DB {
	if db == nil {
		// skip if no password configuration
		if _, ok := os.LookupEnv("PGA_ORIGIN_PASSWORD"); !ok {
			t.SkipNow()
		}

		// wait for pgproxy to start listening
		db = waitForDB(t)
		t.Log("Successfully connected!")

		// check for proper connection handling
		db.SetMaxOpenConns(1)
	}
	return db
}

func waitForDB(t *testing.T) *sql.DB {
	sleepTime := 2 * time.Second
	maxTries := 4

	psqlConnStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	for i := 1; i <= maxTries; i++ {
		t.Log("Attempting to connect...")
		if db, err := sql.Open("postgres", psqlConnStr); err != nil {
			t.Log("Failed to connect to database:", err)
		} else if err = db.Ping(); err != nil {
			t.Log("Failed to ping db:", err)
		} else {
			return db
		}

		sleepTime = sleepTime + time.Second
		t.Log("Failed to open pgproxy port")
		time.Sleep(sleepTime)
	}

	t.Fatalf("Could not reach pgproxy after %d times", maxTries)
	return nil
}

// cols should be formatted "<name> <type>"
func createTable(t *testing.T, db *sql.DB, tableName string, cols []string) {
	var colStr string
	for i, col := range cols {
		if i != 0 {
			colStr += ","
		}
		colStr += col
	}

	t.Logf("Creating table '%s'...", tableName)
	query := fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s(%s)", tableName, colStr)
	if rows, err := db.Query(query); err != nil {
		t.Fatal(err)
	} else {
		rows.Close()
	}
	t.Log("Created table ", tableName)
}

func cleanup(t *testing.T, db *sql.DB) {
	if rows, err := db.Query("DROP TABLE " + table); err != nil {
		t.Fatalf("couldn't DROP table '%s': %v", table, err)
	} else {
		rows.Close()
	}

}

func shutdown(db *sql.DB) {
	if db != nil {
		fmt.Println("Closing DB...")
		if err := db.Close(); err != nil {
			panic(fmt.Sprint("Failed to close DB connection:", err))
		}
		fmt.Println("DB closed.")
	}
}
