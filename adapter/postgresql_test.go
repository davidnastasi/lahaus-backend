package adapter

import (
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/stretchr/testify/suite"
	"lahaus/config"
	"lahaus/domain/model"
	"lahaus/domain/usecases/properties"
	"lahaus/domain/usecases/users"
	"lahaus/infrastructure/storage"
	"testing"
)

type PostgreSQLAdapterSuite struct {
	suite.Suite
	config          *config.Config
	postgresAdapter *PostgreSQLAdapter
}

func TestPostgreSQLAdapterSuite(t *testing.T) {
	suite.Run(t, new(PostgreSQLAdapterSuite))
}

func (suite *PostgreSQLAdapterSuite) SetupTest() {
	suite.config = &config.Config{
		SystemSettings: &config.SystemSettings{
			Storage: &config.Storage{
				Database: &config.Database{
					Host:         "localhost",
					Port:         5432,
					User:         "postgres",
					Password:     "",
					DatabaseName: "lahaus",
				},
			},
		},
	}
	postgresqlManager, err := storage.NewPostgreSQLManager(suite.config.SystemSettings.Storage.Database)
	suite.NoError(err)
	suite.postgresAdapter = NewPostgreSQLAdapter(postgresqlManager)

	//Clean database
	dir := "../infrastructure/storage/migrations"
	suite.Require().NoError(err)
	err = suite.postgresAdapter.MigrationsDown(&dir)
	suite.Require().NoError(err)
	err = suite.postgresAdapter.MigrationsUp(&dir)
	suite.Require().NoError(err)
}

func (suite *PostgreSQLAdapterSuite) TestPostgreSQLAdapter() {
	property := &model.Property{
		Title: "Apartamento cerca a la estación",
		Location: model.Location{
			Longitude: -94.0665887,
			Latitude:  4.6371593,
		},
		Pricing: model.Pricing{
			SalePrice: 450000000,
		},
		PropertyType: model.HOUSE,
		Bedrooms:     3,
		Bathrooms:    2,
		ParkingSpots: nil,
		Area:         60,
		Photos: model.Photos{
			"https://cdn.pixabay.com/photo/2014/08/11/21/39/wall-416060_960_720.jpg",
			"https://cdn.pixabay.com/photo/2016/09/22/11/55/kitchen-1687121_960_720.jpg",
		},
		Status: model.INACTIVE,
	}
	propertyStored, err := suite.postgresAdapter.SaveProperty(property)
	suite.NoError(err)
	suite.NotEqual(int64(0), propertyStored.ID)

	description := "casa elegante"
	administrativeFee := 20000
	propertyStored.Description = &description
	propertyStored.Pricing.AdministrativeFee = &administrativeFee

	propertyUpdated, err := suite.postgresAdapter.UpdateProperty(propertyStored)
	suite.NoError(err)
	suite.Equal(int64(1), propertyUpdated.ID)
	suite.Equal(description, *propertyUpdated.Description)
	suite.Equal(administrativeFee, *propertyUpdated.Pricing.AdministrativeFee)

	filter, err := suite.postgresAdapter.FilterProperties(properties.PropertySearchParams{
		Status:   "INACTIVE",
		Page:     1,
		PageSize: 10,
	})
	suite.NoError(err)
	suite.Equal(int64(1), filter.Total)
	suite.Equal(int64(0), filter.TotalPages)
	suite.Len(filter.Data, 1)

	filter, err = suite.postgresAdapter.FilterProperties(properties.PropertySearchParams{
		Status:   "ACTIVE",
		Page:     1,
		PageSize: 10,
	})
	suite.NoError(err)
	suite.Equal(int64(0), filter.Total)
	suite.Equal(int64(0), filter.TotalPages)
	suite.Len(filter.Data, 0)

	user := &model.User{Email: "david@mail.com", Password: "sarasa"}
	err = suite.postgresAdapter.SaveUser(user)
	suite.NoError(err)
	userStored, found, err := suite.postgresAdapter.GetUser(user.Email)
	suite.NoError(err)
	suite.True(found)
	suite.Equal(int64(1), userStored.ID)

	userNotFound, found, err := suite.postgresAdapter.GetUser("nada@noexiste.com")
	suite.NoError(err)
	suite.False(found)
	suite.Nil(userNotFound)

	err = suite.postgresAdapter.AddFavourite(userStored.ID, propertyStored.ID)
	suite.NoError(err)

	list, err := suite.postgresAdapter.ListFavourites(
		users.FavouritesSearchParams{UserID: userStored.ID, Page: 1, PageSize: 10})
	suite.NoError(err)
	suite.Len(list.Data, 0)

	propertyUpdated.Status = model.ACTIVE
	_, err = suite.postgresAdapter.UpdateProperty(propertyUpdated)
	suite.NoError(err)
	list, err = suite.postgresAdapter.ListFavourites(
		users.FavouritesSearchParams{UserID: userStored.ID, Page: 1, PageSize: 10})
	suite.NoError(err)
	suite.Len(list.Data, 1)

}
