package main

import (
	infinote "boilerplate"
	"boilerplate/api"
	"boilerplate/bindata"
	"boilerplate/seed"
	"boilerplate/store"
	"log"

	"github.com/kelseyhightower/envconfig"
	"github.com/volatiletech/sqlboiler/boil"

	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	migrate_bindata "github.com/golang-migrate/migrate/v4/source/go_bindata"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/oklog/run"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

const version = "v0.0.1"

var app *kingpin.Application

const dbVersionCmd = "db-version"
const dbUpCmd = "db-up"
const dbDownCmd = "db-down"
const dbMigrateCmd = "db-migrate"
const dbDropCmd = "db-drop"
const dbSeedCmd = "db-seed"
const configCmd = "config"
const serveCmd = "serve"

func init() {
	app = kingpin.New("boilerplate", "boilerplate hosting platform")
	app.Version(version)
	app.Command(serveCmd, "begin webserver").Default()
	app.Command(dbVersionCmd, "step-up database")
	app.Command(dbUpCmd, "step-up database")
	app.Command(dbDownCmd, "step-down database")
	app.Command(dbMigrateCmd, "migrate database all the way up")
	app.Command(dbDropCmd, "drop database")
	app.Command(dbSeedCmd, "seed the database")
	app.Command(configCmd, "print environment variables to override config")
}
func main() {
	cmd := kingpin.MustParse(app.Parse(os.Args[1:]))
	config := &infinote.PlatformConfig{}
	err := envconfig.Process("boilerplate", config)
	if err != nil {
		log.Fatal(err)
	}

	switch cmd {
	case dbVersionCmd:
		conn := connect(
			config.Database.User,
			config.Database.Pass,
			config.Database.Host,
			config.Database.Port,
			config.Database.Name,
		)
		m, err := newMigrateInstance(conn)
		if err != nil {
			fmt.Println(err)
			return
		}
		dbVersion, dirty, err := m.Version()
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println("version:", dbVersion, "dirty", dirty)
		return
	case dbUpCmd:
		conn := connect(
			config.Database.User,
			config.Database.Pass,
			config.Database.Host,
			config.Database.Port,
			config.Database.Name,
		)
		m, err := newMigrateInstance(conn)
		if err != nil {
			fmt.Println(err)
			return
		}
		err = m.Steps(1)
		if err != nil {
			fmt.Println(err)
			return
		}
		return
	case dbDownCmd:
		conn := connect(
			config.Database.User,
			config.Database.Pass,
			config.Database.Host,
			config.Database.Port,
			config.Database.Name,
		)
		m, err := newMigrateInstance(conn)
		if err != nil {
			fmt.Println(err)
			return
		}
		err = m.Steps(-1)
		if err != nil {
			fmt.Println(err)
			return
		}
		return
	case dbMigrateCmd:
		conn := connect(
			config.Database.User,
			config.Database.Pass,
			config.Database.Host,
			config.Database.Port,
			config.Database.Name,
		)
		m, err := newMigrateInstance(conn)
		if err != nil {
			fmt.Println(err)
			return
		}
		err = m.Up()
		if err != nil {
			fmt.Println(err)
			return
		}
		return
	case dbDropCmd:
		conn := connect(
			config.Database.User,
			config.Database.Pass,
			config.Database.Host,
			config.Database.Port,
			config.Database.Name,
		)
		m, err := newMigrateInstance(conn)
		if err != nil {
			fmt.Println(err)
			return
		}
		err = m.Drop()
		if err != nil {
			fmt.Println(err)
			return
		}
		return

	case configCmd:
		err = infinote.PrintConfigVars()
		if err != nil {
			fmt.Println(err)
		}
		return
	case dbSeedCmd:
		conn := connect(
			config.Database.User,
			config.Database.Pass,
			config.Database.Host,
			config.Database.Port,
			config.Database.Name,
		)
		err = seed.Run(conn)
		if err != nil {
			fmt.Println(err)
		}
		return
	default:
		conn := connect(
			config.Database.User,
			config.Database.Pass,
			config.Database.Host,
			config.Database.Port,
			config.Database.Name,
		)
		g := &run.Group{}
		ctx, cancel := context.WithCancel(context.Background())

		g.Add(func() error {
			logger := infinote.NewLogToStdOut("api", version, false)
			userStore := store.NewUserStore(conn)
			companyStore := store.NewCompanyStore(conn)
			notesStore := store.NewNoteStore(conn)
			roleStorer := store.NewRoleStore(conn)
			tokenStore := store.NewTokenStore(conn)
			blacklistRefreshHours := config.UserAuth.BlacklistRefreshHours
			blistProvider := infinote.NewBlacklister(logger, tokenStore, blacklistRefreshHours)

			//			subResolver := boilerplate.NewSubHub()

			jwtSecret := config.UserAuth.JWTSecret
			auther := infinote.NewAuther(jwtSecret, userStore, blistProvider, tokenStore, config.UserAuth.TokenExpiryDays)
			APIController := api.NewAPIController(&api.ControllerOpts{
				NoteStorer:        notesStore,
				UserStorer:        userStore,
				CompanyStorer:     companyStore,
				BlacklistProvider: blistProvider,
				TokenStorer:       tokenStore,
				RoleStorer:        roleStorer,
				//SubscriptionResolver: subResolver,
				JWTSecret: jwtSecret,
				Auther:    auther,
				Logger:    logger,
			})

			server := &infinote.APIService{
				Log:  logger,
				Addr: config.API.Addr,
			}
			return server.Run(ctx, APIController)
		}, func(err error) {
			fmt.Println(err)
			cancel()
		})
		g.Add(func() error {
			l := infinote.NewLogToStdOut("loadbalancer", version, false)
			lb := infinote.LoadbalancerService{Addr: config.LoadBalancer.Addr, Log: l}
			return lb.Run(ctx, config.LoadBalancer.Addr, config.API.Addr, config.LoadBalancer.RootPath)
		}, func(error) {
			fmt.Println(err)
			cancel()
			return
		})
		g.Add(func() error {
			c := make(chan os.Signal, 1)
			signal.Notify(c, os.Interrupt)

			select {
			case <-c:
				return errors.New("ctrl-c caught, exiting gracefully")
			}

		}, func(error) {
			fmt.Println(err)
			cancel()
			return
		})

		err = g.Run()
		if err != nil {
			fmt.Println(err)
			return
		}
	}

}

func newMigrateInstance(conn *sqlx.DB) (*migrate.Migrate, error) {
	s := migrate_bindata.Resource(bindata.AssetNames(),
		func(name string) ([]byte, error) {
			return bindata.Asset(name)
		})
	d, err := migrate_bindata.WithInstance(s)
	if err != nil {
		return nil, fmt.Errorf("bindata instance: %w", err)
	}
	dbDriver, err := postgres.WithInstance(conn.DB, &postgres.Config{})
	if err != nil {
		return nil, fmt.Errorf("db instance: %w", err)
	}
	m, err := migrate.NewWithInstance("go-bindata", d, "postgres", dbDriver)
	if err != nil {
		return nil, fmt.Errorf("migrate instance: %w", err)
	}
	return m, nil
}

func connect(
	DatabaseUser string,
	DatabasePass string,
	DatabaseHost string,
	DatabasePort string,
	DatabaseName string,

) *sqlx.DB {
	connString := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		DatabaseUser,
		DatabasePass,
		DatabaseHost,
		DatabasePort,
		DatabaseName,
	)
	conn, err := sqlx.Connect("postgres", connString)
	if err != nil {
		log.Fatal("could not initialise database:", err)
	}
	if conn == nil {
		panic("conn is nil")
	}

	boil.SetDB(conn)
	return conn
}
