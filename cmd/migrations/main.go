package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"sort"
	"strings"

	"hellper/internal/config"
	"hellper/internal/model/sql"

	_ "github.com/lib/pq"
)

func main() {
	migrationFiles, err := getMigrationFiles()
	if err != nil {
		log.Fatal("Could not read migrations directory")
		os.Exit(1)
	}

	for _, migration := range migrationFiles {
		err = executeMigration(migration)
		if err != nil {
			log.Fatal(fmt.Sprintf("Could not execute migration %s", migration))
			fmt.Println(err)
			os.Exit(2)
		}
	}
}

func getMigrationFiles() ([]string, error) {
	dirName := "./internal/model/sql/postgres/schema"

	files, err := ioutil.ReadDir(dirName)
	if err != nil {
		return nil, err
	}

	fileNames := make([]string, 0, len(files))
	for _, file := range files {
		fileName := file.Name()
		if strings.HasSuffix(fileName, ".sql") {
			fullFileName := dirName + "/" + fileName
			fileNames = append(fileNames, fullFileName)
		}
	}

	sort.Strings(fileNames)

	return fileNames, nil
}

func executeMigration(migrationFile string) error {
	fileContentBytes, err := ioutil.ReadFile(migrationFile)
	if err != nil {
		log.Fatal(fmt.Sprintf("Could not read file %s", migrationFile))
		return err
	}

	fileContent := string(fileContentBytes)
	fmt.Println(config.Env.DSN)
	db := sql.NewDBWithDSN(config.Env.Database, config.Env.DSN)
	_, err = db.Query(fileContent)

	fmt.Println(err)

	return err
}
