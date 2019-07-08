package migrations

import (
	"log"

	"github.com/go-pg/migrations"
	"github.com/go-pg/pg"
	"github.com/go-pg/pg/orm"
)

var (
	defaultCreateOptions *orm.CreateTableOptions
	defaultDropOptions   *orm.DropTableOptions
)

func init() {
	defaultCreateOptions = &orm.CreateTableOptions{
		FKConstraints: true,
	}

	defaultDropOptions = &orm.DropTableOptions{
		IfExists: false,
		Cascade:  true,
	}
}

// Run forwards args to https://godoc.org/github.com/go-pg/migrations#Run
func Run(db *pg.DB, a ...string) (oldVersion, newVersion int64, err error) {
	return migrations.Run(db, a...)
}

// MustRunUp initializes migrations (if necessary) then applies all new migrations; it panics on error
func MustRunUp(db *pg.DB) {
	// init migrations if not already initialized
	_, _, err := migrations.Run(db, "init")
	if err != nil {
		// ignore error if migration table already exists
		if err.Error() != "ERROR #42P07 relation \"gopg_migrations\" already exists" {
			log.Fatalf("migrations: %s", err.Error())
		}
	}

	// run migrations
	oldVersion, newVersion, err := migrations.Run(db, "up")
	if err != nil {
		log.Fatalf("migrations: %s", err.Error())
	}

	// log version status
	if newVersion != oldVersion {
		log.Printf("migrated from version %d to %d", oldVersion, newVersion)
	} else {
		log.Printf("running at migration version %d", oldVersion)
	}
}