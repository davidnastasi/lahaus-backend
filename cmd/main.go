package main

import (
	"flag"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	_ "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"go.uber.org/zap"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"lahaus/adapter"
	"lahaus/config"
	ucproperties "lahaus/domain/usecases/properties"
	"lahaus/domain/usecases/ruler"
	ucusers "lahaus/domain/usecases/users"
	"lahaus/infrastructure/api"
	"lahaus/infrastructure/api/middlewares"
	"lahaus/infrastructure/storage"
	"lahaus/logger"
	"log"
	"net/http"
)

var yamlPathFlag = flag.String("config", "./config.yml", "Specify the path of config.yml file, e.g.: -config /folder/config.yml")
var migrationDir = flag.String("migration", "./infrastructure/storage/migrations", "Directory where the migration files are located")

func main() {

	flag.Parse()
	conf, err := configure(*yamlPathFlag)
	if err != nil {
		logger.GetInstance().Fatal("error loading configurations", zap.Error(err))
	}

	setLoggingLevel(conf.SystemSettings.Logger.Level)

	// CreateProperty the repository
	repositoryPostgresSQL, err := storage.NewPostgreSQLManager(conf.SystemSettings.Storage.Database)
	if err != nil {
		logger.GetInstance().Fatal("failed to instantiate postgreSQL repository", zap.Error(err))
	}

	// create the adapters
	databaseAdapter := adapter.NewPostgreSQLAdapter(repositoryPostgresSQL)

	err = databaseAdapter.MigrationsUp(migrationDir)
	if err != nil {
		logger.GetInstance().Fatal("failed to run migrations", zap.Error(err))
	}

	// Create the usecases
	rulerUserCase := ruler.NewPropertyRulerUseCase(conf)
	createPropertyUseCase := ucproperties.NewCreatePropertyUseCase(databaseAdapter, rulerUserCase)
	updatePropertyUseCase := ucproperties.NewUpdatePropertyUseCase(databaseAdapter, rulerUserCase)
	searchPropertiesUseCase := ucproperties.NewSearchPropertyUseCase(databaseAdapter)

	signInUserExecutor := ucusers.NewSignInUserUseCase(databaseAdapter)
	signUpUserExecutor := ucusers.NewSignUpUserUseCase(conf.SystemSettings.Security, databaseAdapter)
	addFavouriteUserExecutor := ucusers.NewAddFavouriteUseCase(databaseAdapter)
	listFavouriteUserExecutor := ucusers.NewListFavouriteUseCase(databaseAdapter)

	// Create handlers
	handlerProperties := api.NewPropertyHandler(createPropertyUseCase, updatePropertyUseCase, searchPropertiesUseCase)
	handlerUser := api.NewUserHandler(signInUserExecutor, signUpUserExecutor, addFavouriteUserExecutor, listFavouriteUserExecutor)

	// Create web routing
	router := chi.NewRouter()
	router.Use(render.SetContentType(render.ContentTypeJSON))
	router.Use(middleware.RequestID, middleware.NoCache, middleware.Logger, middleware.Recoverer)

	authenticationMiddleware := middlewares.NewAuthenticationMiddleware(conf.SystemSettings.Security)

	router.Route("/v1", func(r chi.Router) {
		r.Route("/properties", func(r chi.Router) {
			r.Post("/", handlerProperties.CreateProperty)
			r.Put("/{id}", handlerProperties.UpdateProperty)
			r.Get("/", handlerProperties.SearchProperties)
		})

		r.Route("/users", func(r chi.Router) {
			r.Post("/", handlerUser.SignInUser)
			r.Post("/login", handlerUser.SignUpUser)
			r.Route("/me/favourites", func(r chi.Router) {
				r.Use(authenticationMiddleware.Execute)
				r.Post("/", handlerUser.AddFavourite)
				r.Get("/", handlerUser.ListFavourites)
			})
		})
	})

	log.Fatal(http.ListenAndServe(":8080", router))

}

func configure(path string) (*config.Config, error) {
	logger.GetInstance().Info("loading configurations from file", zap.String("path", path))
	conf := &config.Config{}
	//#nosec
	file, err := ioutil.ReadFile(path)
	if err != nil {
		logger.GetInstance().Error("error reading config file", zap.Error(err))
		return nil, err
	}
	err = yaml.Unmarshal(file, conf)
	if err != nil {
		logger.GetInstance().Error("error unmarshal config file", zap.Error(err))
		return nil, err
	}
	return conf, err
}

func setLoggingLevel(level string) {
	l := zap.InfoLevel
	_ = l.Set(level)
	logger.GetAtomLevel().SetLevel(l)
}
