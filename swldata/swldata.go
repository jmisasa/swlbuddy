package swldata

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"os"
	"strconv"
	"time"
)

type SwlLine struct {
	Station     string
	CountryCode string
	Language    string
}

func GetByFrequency(hz string) []SwlLine {
	userConfigDir, _ := os.UserConfigDir()
	pathToDb := userConfigDir + "/swlbuddy/eibi.sqlite"
	db, err := sql.Open("sqlite3", "file:"+pathToDb+"?cache=shared")

	if err != nil {
		panic(fmt.Sprintf("Error accessing DB: %v", err))
	}

	_, utcHour := getCurrentUTC()

	query := `
		SELECT station, itu_code, language 
		FROM eibi 
		WHERE 
			khz LIKE $1 || '%'
			AND utc_start <= $2 AND utc_end >= $3
	`

	lines := []SwlLine{}
	rows, err := db.Query(query, getKhz(hz), utcHour, utcHour)

	if err != nil {
		panic(fmt.Sprintf("Error querying eibi table: %v", err))
	}

	for rows.Next() {
		var line SwlLine
		err = rows.Scan(&line.Station, &line.CountryCode, &line.Language)

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

func getCurrentUTC() (ymd string, his string) {
	currentUTC := time.Now().UTC()

	return currentUTC.Format("2006-01-02"), currentUTC.Format("1504")
}
