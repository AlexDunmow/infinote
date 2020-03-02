package boilerplate

import (
	"github.com/kelseyhightower/envconfig"
)

// PlatformConfig for the Platform
type PlatformConfig struct {
	API          *API
	LoadBalancer *LoadBalancer
	Database     *Database
	UserAuth     *UserAuth
}

// API for the API service
type API struct {
	Addr string `desc:"host:port to run the API" default:":8081"`
}

// LoadBalancer for the LoadBalancer service
type LoadBalancer struct {
	Addr     string `desc:"host:port to run caddy" default:":8080"`
	RootPath string `desc:"folder path of index.html" default:"../web/dist"`
}

// Database for the Database service
type Database struct {
	User string `desc:"Postgres username" default:"boilerplate"`
	Pass string `desc:"Postgres password" default:"dev"`
	Host string `desc:"Postgres host" default:"localhost"`
	Port string `desc:"Postgres port" default:"5438"`
	Name string `desc:"Postgres database name" default:"boilerplate"`
}

// UserAuth holds variables for user auth config, such as token expiry
type UserAuth struct {
	JWTSecret             string `desc:"JWT secret" default:"872ab3df-d7c7-4eb6-a052-4146d0f4dd15"`
	TokenExpiryDays       int    `desc:"How many days before the token expires" default:"30"`
	BlacklistRefreshHours int    `desc:"How often should the issued_tokens list be cleared of expired tokens in hours" default:"1"`
}

// PrintConfigVars of the platform config struct
func PrintConfigVars() error {
	return envconfig.Usage("boilerplate", &PlatformConfig{})
}
