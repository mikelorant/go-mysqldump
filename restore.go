package mysqldump

import (
	"bufio"
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"
)

func Restore(filepath string, db *sql.DB) error {
	file, err := os.Open(filepath)

	if err != nil {
		log.Fatalf("failed to open")

	}

	// The bufio.NewScanner() function is called in which the
	// object os.File passed as its parameter and this returns a
	// object bufio.Scanner which is further used on the
	// bufio.Scanner.Split() method.
	scanner := bufio.NewScanner(file)

	// The bufio.ScanLines is used as an
	// input to the method bufio.Scanner.Split()
	// and then the scanning forwards to each
	// new line using the bufio.Scanner.Scan()
	// method.
	buf := make([]byte, 0, 1024*1024*1024)
	scanner.Buffer(buf, 1024*1024)
	scanner.Split(bufio.ScanLines)
	var query string
	tx, err := db.Begin()
	for scanner.Scan() {
		text := strings.Trim(scanner.Text(), " ")
		if text != "" && !strings.HasPrefix(text, "-") {
			if query == "" {
				query = text
			} else {
				query = strings.Join([]string{query, text}, " ")
			}
		}
		if strings.HasSuffix(text, ";") {
			r, err := tx.Exec(query)
			if err != nil {
				tx.Rollback()
				return err
			} else {
				fmt.Println(r.RowsAffected())
			}
			query = ""
		}
	}
	if err := scanner.Err(); err != nil {
		tx.Rollback()
		panic(err)
	}
	tx.Commit()
	// The method os.File.Close() is called
	// on the os.File object to close the file
	file.Close()
	return nil
}
