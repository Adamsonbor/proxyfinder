package main

import (
	"context"
	"flag"
	"os"


	"github.com/pressly/goose/v3"
	_ "modernc.org/sqlite"
)

func main() {
	flags := flag.NewFlagSet("migrations", flag.ExitOnError)
	flags.String("dir", "./migrations", "directory with migration files")
	flags.Parse(os.Args[1:])
	args := flags.Args()

	if len(args) < 3 {
		flags.Usage()
		return
	}

	dbpath, command := args[1], args[2]

	db, err := goose.OpenDBWithDriver("sqlite3", dbpath)
	if err != nil {
		panic(err)
	}
	defer func() {
		err := db.Close()
		if err != nil {
			panic(err)
		}
	}()

	var arguments []string
	if len(args) > 3 {
		arguments = args[3:]
	}

	err = goose.RunContext(context.Background(), command, db, args[0], arguments...)
	if err != nil {
		panic(err)
	}
}
