package main

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/superg3m/server/Utility"
)

const STOIC_MIGRATION_UP_STR = "-- StoicMigration Up"
const STOIC_MIGRATION_DOWN_STR = "-- StoicMigration Down"

type MigrationMode int

const (
	MIGRATION_MODE_UP MigrationMode = iota
	MIGRATION_MODE_DOWN
)

func getSqlCommandsFromFile(mode MigrationMode, filePath string) ([]string, error) {
	migrationStr := []string{"-- StoicMigration Up\n", "-- StoicMigration Down\n"}
	delimitor := ';'

	otherMode := int(mode)
	Utility.ToggleBit(&otherMode, 0)

	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file %s: %w", filePath, err)
	}
	string_data := strings.ReplaceAll(string(data), "\r\n", "\n")

	if !strings.Contains(string_data, migrationStr[mode]) {
		return nil, fmt.Errorf("migration file doesn't contain %s", migrationStr[mode])
	}

	lines := strings.SplitAfter(string_data, "\n")

	var ret []string
	var charAccumulator strings.Builder
	insideMigration := false

	for _, line := range lines {
		if !insideMigration && line != migrationStr[mode] {
			continue
		}

		if line == migrationStr[mode] {
			insideMigration = true
			continue
		}

		if line == migrationStr[otherMode] {
			break
		}

		for _, c := range line {
			charAccumulator.WriteByte(byte(c))
			if c == delimitor {
				ret = append(ret, charAccumulator.String())
				charAccumulator.Reset()
			}
		}
	}

	return ret, nil
}

func findFilesWithExtension(root, ext string) ([]string, error) {
	info, err := os.Stat(root)
	Utility.AssertOnErrorMsg(err, "Failed to access the root directory")
	if !info.IsDir() {
		return nil, fmt.Errorf("provided root path is not a directory: %s", root)
	}

	var files []string

	err = filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			Utility.AssertOnErrorMsg(err, fmt.Sprintf("Error accessing path: %s", path))
		}

		if d.IsDir() {
			return nil
		}

		if filepath.Ext(d.Name()) == ext {
			files = append(files, path)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return files, nil
}

// Note(Jovanni):
// If I was doing something more serious you would use some type of
// .env to hide these constants
const (
	DB_ENGINE = "mysql"
	HOST      = "localhost"
	PORT      = 3306
	USER      = "root"
	PASSWORD  = "P@55word"
	DBNAME    = "stoic"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: <program> [up|down]")
		os.Exit(1)
	}

	arg := os.Args[1]
	mode := MIGRATION_MODE_DOWN

	switch arg {
	case "up":
		mode = MIGRATION_MODE_UP
	case "down":
		mode = MIGRATION_MODE_DOWN
	default:
		fmt.Printf("Invalid argument: %s\n", arg)
		fmt.Println("Valid options are: 'up' or 'down'")
		os.Exit(1)
	}

	dsn := Utility.GetDSN(DB_ENGINE, HOST, PORT, USER, PASSWORD, DBNAME)
	db, err := sqlx.Connect(DB_ENGINE, dsn)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	files, _ := findFilesWithExtension(fmt.Sprintf("./%s", DB_ENGINE), ".sql")

	for _, file := range files {
		sqlUpCommands, err := getSqlCommandsFromFile(mode, file)
		Utility.AssertOnError(err)

		if mode == MIGRATION_MODE_UP {
			Utility.LogSuccess("Migration Up: %s", file)
		} else {
			Utility.LogDebug("Migration Down: %s", file)
		}

		for _, element := range sqlUpCommands {
			_, err := db.Exec(element)
			Utility.AssertOnError(err)
		}
	}
}
