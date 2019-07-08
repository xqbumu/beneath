package control

import (
	"fmt"
	"log"
	"net/http"

	"github.com/beneath-core/beneath-go/control/db"
	"github.com/beneath-core/beneath-go/control/gql"
	"github.com/beneath-core/beneath-go/control/migrations"
	"github.com/beneath-core/beneath-go/control/resolver"
	"github.com/beneath-core/beneath-go/core"

	"github.com/99designs/gqlgen/handler"
	"github.com/go-chi/chi"
	"github.com/rs/cors"
)

type configSpecification struct {
	HTTPPort    int    `envconfig:"PORT" default:"4000"`
	RedisURL    string `envconfig:"REDIS_URL" required:"true"`
	PostgresURL string `envconfig:"POSTGRES_URL" required:"true"`
}

var (
	// Config for control
	Config configSpecification
)

func init() {
	// load config
	core.LoadConfig("beneath", &Config)

	// connect postgres and redis
	db.InitPostgres(Config.PostgresURL)
	db.InitRedis(Config.RedisURL)

	// run migrations
	migrations.MustRunUp(db.DB)
}

// ListenAndServeHTTP serves the GraphQL API on HTTP
func ListenAndServeHTTP(port int) error {
	router := chi.NewRouter()

	// Add CORS
	router.Use(cors.New(cors.Options{
		AllowedOrigins: []string{
			"http://localhost:8080",
		},
		AllowCredentials: true,
		Debug:            true,
	}).Handler)

	// Add playground
	router.Handle("/", handler.Playground("Beneath", "/graphql"))

	// Add graphql server
	router.Handle("/graphql",
		handler.GraphQL(gql.NewExecutableSchema(gql.Config{Resolvers: &resolver.Resolver{}})),
	)

	// Serve
	log.Printf("HTTP server running on port %d\n", port)
	return http.ListenAndServe(fmt.Sprintf(":%d", port), router)
}