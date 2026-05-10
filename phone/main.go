package main

import (
	"bufio"
	"database/sql"
	"fmt"
	"os"
	"strings"
	"unicode"

	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = "5432"
	user     = "postgres"
	password = "root"
	dbname   = "gophercises_phone"

	phoneNumbers_file = "phone.txt"
)

func main() {
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s sslmode=disable",
		host, port, user, password)
	db, err := sql.Open("postgres", psqlInfo)
	must(err)
	must(resetDB(db, dbname))
	db.Close()

	psqlInfo = fmt.Sprintf("%s dbname=%s", psqlInfo, dbname)
	db, err = sql.Open("postgres", psqlInfo)
	must(err)
	defer db.Close()

	must(createPhoneNumberTable(db))

	pNs, err := getPhoneNumbers(phoneNumbers_file)
	must(err)

	for _, pn := range pNs {
		_, err := insertPhone(db, pn)
		must(err)
	}

	phones, err := allPhone(db)
	must(err)
	for _, p := range phones {
		fmt.Printf("Working on... %v\n", p.number)
		number := normalize(p.number)
		if p.number != number {
			existing, err := findPhone(db, number)
			must(err)
			if existing != nil {
				p.number = number
				must(updatePhone(db, p))
			} else {
				must(deletePhone(db, p.id))
			}
		} else {
			fmt.Println("No changes")
		}
	}

}

type phone struct {
	id     int
	number string
}

func getPhoneNumbers(filename string) ([]string, error) {
	var ret []string

	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		ret = append(ret, scanner.Text())
	}

	return ret, nil
}

func updatePhone(db *sql.DB, p phone) error {
	_, err := db.Exec("UPDATE phone_numbers SET value=$2 WHERE id=$1", p.id, p.number)
	return err
}

func insertPhone(db *sql.DB, phone string) (int, error) {
	statement := `INSERT INTO phone_numbers(value) VALUES($1) RETURNING id`
	var id int
	err := db.QueryRow(statement, phone).Scan(&id)
	if err != nil {
		return -1, err
	}
	return id, nil
}

func deletePhone(db *sql.DB, id int) error {
	_, err := db.Exec("DELETE FROM phone_numbers WHERE id=$1", id)
	return err
}

func allPhone(db *sql.DB) ([]phone, error) {
	var ret []phone

	rows, err := db.Query("SELECT * FROM phone_numbers")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var p phone
		if err = rows.Scan(&p.id, &p.number); err != nil {
			return nil, err
		}
		ret = append(ret, p)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return ret, nil
}

func getPhone(db *sql.DB, id int) (string, error) {
	var number string
	row := db.QueryRow("SELECT * FROM phone_numbers WHERE id = $1", id)
	err := row.Scan(&id, &number)
	if err != nil {
		return "", err
	}
	return number, nil
}

func findPhone(db *sql.DB, number string) (*phone, error) {
	var p phone
	row := db.QueryRow("SELECT * FROM phone_numbers WHERE value = $1", number)
	err := row.Scan(&p.id, &p.number)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		} else {
			return nil, err
		}
	}
	return &p, nil
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}

func resetDB(db *sql.DB, name string) error {
	_, err := db.Exec("DROP DATABASE IF EXISTS " + name)
	if err != nil {
		return err
	}
	return createDB(db, name)
}

func createDB(db *sql.DB, name string) error {
	_, err := db.Exec("CREATE DATABASE " + name)
	return err
}

func createPhoneNumberTable(db *sql.DB) error {
	statement := `
		CREATE TABLE IF NOT EXISTS phone_numbers (
			id SERIAL,
			value VARCHAR(255)
		)
	`
	_, err := db.Exec(statement)
	return err
}

func normalize(phone string) string {
	ret := strings.Builder{}
	for _, ch := range phone {
		if unicode.IsNumber(ch) {
			ret.WriteRune(ch)
		}
	}
	return ret.String()
}

// regular expression
// func normalize(phone string) string {
// 	re := regexp.MustCompile("\\D")
// 	return re.ReplaceAllString(phone, "")
// }
