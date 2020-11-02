package migrationsutil

import (
	"context"
	"fmt"

	"github.com/go-pg/migrations/v7"
	"github.com/go-pg/pg/v9"
	"github.com/spf13/cobra"

	"gitlab.com/beneath-hq/beneath/infrastructure/db"
	"gitlab.com/beneath-hq/beneath/pkg/envutil"
	"gitlab.com/beneath-hq/beneath/pkg/log"
)

// Migrator wraps go-pg/migrations with useful extra functionality
type Migrator struct {
	*migrations.Collection
	TableName string
}

// New creates a new migrator
func New(tableName string) *Migrator {
	collection := migrations.NewCollection()
	if tableName != "" {
		collection.SetTableName(tableName)
	}
	collection.DisableSQLAutodiscover(true)
	return &Migrator{
		Collection: collection,
		TableName:  tableName,
	}
}

// AddCmd registers a migration CLI with Cobra.
// The callback carries args that can be supplied outright to RunWithArgs.
func (m *Migrator) AddCmd(root *cobra.Command, name string, fn func(args []string)) {
	runWithoutName := func(cmd *cobra.Command, args []string) { fn(args) }
	runWithName := func(cmd *cobra.Command, args []string) { fn(append([]string{cmd.Name()}, args...)) }

	cmd := &cobra.Command{
		Use:   name,
		Short: "Manages migrations on the db (runs \"up\" if no subcommand provided)",
		Args:  cobra.NoArgs,
		Run:   runWithoutName,
	}
	root.AddCommand(cmd)

	cmd.AddCommand(&cobra.Command{
		Use:   "init",
		Short: "Creates version info table in the database",
		Args:  cobra.NoArgs,
		Run:   runWithName,
	})

	cmd.AddCommand(&cobra.Command{
		Use:   "up [target]",
		Short: "Runs available migrations up to the target (or all if target not set)",
		Args:  cobra.MaximumNArgs(1),
		Run:   runWithName,
	})

	cmd.AddCommand(&cobra.Command{
		Use:   "down",
		Short: "Reverts last migration",
		Args:  cobra.NoArgs,
		Run:   runWithName,
	})

	cmd.AddCommand(&cobra.Command{
		Use:   "reset",
		Short: "Reverts all migrations",
		Args:  cobra.NoArgs,
		Run:   runWithName,
	})

	cmd.AddCommand(&cobra.Command{
		Use:   "version",
		Short: "Prints latest migration version",
		Args:  cobra.NoArgs,
		Run:   runWithName,
	})

	cmd.AddCommand(&cobra.Command{
		Use:   "set_version [version]",
		Short: "",
		Args:  cobra.ExactArgs(1),
		Run:   runWithName,
	})
}

// RunWithArgs runs migrations with command-line args (as registered with RegisterCmd)
func (m *Migrator) RunWithArgs(db db.DB, args ...string) {
	pg := db.GetDB(context.Background()).(*pg.DB)
	oldVersion, newVersion, err := m.Collection.Run(pg, args...)
	m.log(oldVersion, newVersion, err)
}

// Automigrate automatically applies every new migration. If reset is true, it resets migrations
// before automigrating.
func (m *Migrator) Automigrate(db db.DB, reset bool) (oldVersion, newVersion int64, err error) {
	pgDb := db.GetDB(context.Background()).(*pg.DB)

	// safety check on reset in production
	if reset && envutil.GetEnv() == envutil.Production {
		return 0, 0, fmt.Errorf("Cannot automigrate with reset=true in production")
	}

	// run init only if necessary
	_, err = m.Version(pgDb)
	if err != nil {
		if err.Error() != fmt.Sprintf(`ERROR #42P01 relation "%s" does not exist`, m.TableName) {
			return 0, 0, err
		}

		_, _, err := m.Collection.Run(pgDb, "init")
		if err != nil {
			return 0, 0, err
		}
	}

	// run reset if requested
	if reset {
		_, _, err := m.Collection.Run(pgDb, "reset")
		if err != nil {
			return 0, 0, err
		}
	}

	// run up migrations
	oldVersion, newVersion, err = m.Collection.Run(pgDb, "up")
	if err != nil {
		return oldVersion, newVersion, err
	}

	return oldVersion, newVersion, err
}

// AutomigrateAndLog runs m.Automigrate and logs the result, returning only an error if applicable
func (m *Migrator) AutomigrateAndLog(db db.DB, reset bool) error {
	oldVersion, newVersion, err := m.Automigrate(db, reset)
	m.log(oldVersion, newVersion, err)
	return err
}

func (m *Migrator) log(oldVersion, newVersion int64, err error) {
	if err != nil {
		log.S.Errorf("failed running migrations on '%s': %s", m.TableName, err.Error())
	} else if newVersion != oldVersion {
		log.S.Infof("migrated '%s' from version %d to %d", m.TableName, oldVersion, newVersion)
	} else {
		log.S.Infof("migration version is %d on '%s'", newVersion, m.TableName)
	}
}