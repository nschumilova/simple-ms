package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	healthhttp "github.com/nschumilova/simple-ms/pkg/healthcheck/port/http"
	usersql "github.com/nschumilova/simple-ms/pkg/user/adapter/sql"
	userhttp "github.com/nschumilova/simple-ms/pkg/user/port/http"
	userusecase "github.com/nschumilova/simple-ms/pkg/user/usecase"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type connectionPool struct {
	MaxIdleConns    int    `mapstructure:"max_idle_connections"`
	MaxOpenConns    int    `mapstructure:"max_idle_connections"`
	ConnMaxLifetime string `mapstructure:"connection_max_lifetime"`
}

type postgresDb struct {
	Host     string
	Port     int
	Database string `mapstructure:"db_name"`
	User     string
	Password string
	Pool connectionPool `mapstructure:"connection_pool"`
}

func main() {
	logger, _ := zap.NewProduction()
	defer logger.Sync()
	log := logger.Sugar()

	//config
	if err := setupConfig(); err != nil {
		log.Fatal(err)
	}

	//postgres adapter
	db := configurePostgres(log)
	userRepository := usersql.NewUserRepository(db)

	//usecases
	userUC := userusecase.NewUseCase(userRepository)

	//handlers
	httpHealthHandler := &healthhttp.HeathCheckHandler{
		Log: log,
	}
	httpUserHandler := userhttp.NewHandler(userUC, log)
	router := chi.NewRouter()
	router.Get("/health/", httpHealthHandler.Health)
	router.Route("/user", func(r chi.Router) {
		r.Post("/", httpUserHandler.Create)
		r.Route("/{userID}", func(r chi.Router) {
			r.Get("/", httpUserHandler.Get)
			r.Put("/", httpUserHandler.Update)
			r.Delete("/", httpUserHandler.Delete)
		})
	})
	port := viper.GetInt("port")
	log.Error(http.ListenAndServe(fmt.Sprintf(":%d", port), router))
}

func setupConfig() error {
	pflag.String("configs", "./configs/", "path to config files directory")
	pflag.Int("port", 8000, "application port")
	pflag.Parse()
	viper.BindPFlags(pflag.CommandLine)

	configFile := func(name string) {
		viper.SetConfigName(name)
		viper.SetConfigType("yaml")
		viper.AddConfigPath(viper.GetString("configs"))
	}
	configFile("application")
	if err := viper.ReadInConfig(); err != nil {
		return fmt.Errorf("failed to read application configuration: %w", err)
	}
	configFile("secrets")
	if err := viper.MergeInConfig(); err != nil {
		return fmt.Errorf("failed to read secrets configuration: %w", err)
	}
	return nil
}

func configurePostgres(log *zap.SugaredLogger) *gorm.DB {
	dbConfig := postgresDb{}
	if err := viper.UnmarshalKey("database.postgres", &dbConfig); err != nil {
		log.Fatal("failed to load postgres database configuration", err)
	}
	dsn := fmt.Sprintf(
		"host=%v user=%v password=%v dbname=%v port=%d sslmode=disable TimeZone=Europe/Moscow",
		dbConfig.Host, dbConfig.User, dbConfig.Password, dbConfig.Database, dbConfig.Port,
	)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatal(err)
	}
	poolConfig := &dbConfig.Pool
	connMaxLifetime, err := time.ParseDuration(poolConfig.ConnMaxLifetime)
	if err != nil {
		log.Fatal("failed to parse postgres connection max lifetime", err)
	}
	sqlDB.SetMaxIdleConns(poolConfig.MaxIdleConns)
	sqlDB.SetMaxOpenConns(poolConfig.MaxOpenConns)
	sqlDB.SetConnMaxLifetime(connMaxLifetime)
	return db
}
