package swldata

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	_ "golang.org/x/text/number"
	"os"
	"strconv"
)

type SwlLine struct {
	Station     string
	CountryCode string
	Language    string
}

func GetByFrequency(hz string) []SwlLine {
	fmt.Printf("Hz: %s, kHz: %s\n", hz, getKhz(hz))

	userConfigDir, _ := os.UserConfigDir()
	pathToDb := userConfigDir + "/swlbuddy/eibi.sqlite"
	db, err := sql.Open("sqlite3", "file:"+pathToDb+"?cache=shared")

	if err != nil {
		panic(fmt.Sprintf("Error accessing DB: %v", err))
	}

	query := `
		SELECT station, itu_code, language FROM eibi WHERE khz LIKE '%' || $1 || '%'
	`

	lines := []SwlLine{}
	rows, err := db.Query(query, getKhz(hz))

	if err != nil {
		panic(fmt.Sprintf("Error querying eibi table: %v", err))
	}

	for rows.Next() {
		var line SwlLine
		err = rows.Scan(&line.Station, &line.CountryCode, &line.Language)

		fmt.Println("Line:%v", line)

		if err != nil {
			panic("Error processing lines")
		}

		lines = append(lines, line)
	}

	return lines
}

func getKhz(hz string) string {
	intHz, err := strconv.Atoi(hz)

	if err != nil {
		panic(fmt.Sprintf("Error converting Hz to int: %v", err))
	}

	return strconv.Itoa(intHz / 1000)
}
