package adapter

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	"github.com/lib/pq"
	"go.uber.org/zap"
	"lahaus/domain/model"
	"lahaus/domain/usecases/properties"

	"lahaus/domain/usecases/users"
	"lahaus/infrastructure/storage"
	"lahaus/logger"
	"time"
)

// PostgreSQLAdapter represents a postgres database.
type PostgreSQLAdapter struct {
	postgres *storage.PostgreSQLManager
}

// NewPostgreSQLAdapter creates a new PostgreSQLAdapter
func NewPostgreSQLAdapter(postgres *storage.PostgreSQLManager) *PostgreSQLAdapter {
	return &PostgreSQLAdapter{
		postgres: postgres,
	}
}

func (adapter *PostgreSQLAdapter) MigrationsUp(migrationDir *string) error {
	// run migrations
	m, err := migrate.New(
		fmt.Sprintf("file://%s", *migrationDir), adapter.postgres.ConnectionString)

	if err != nil {
		logger.GetInstance().Fatal("failed to creating migrations", zap.Error(err))
		return err
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		logger.GetInstance().Fatal("an error occurred while syncing the database", zap.Error(err))
		return err
	}

	return nil
}

func (adapter *PostgreSQLAdapter) MigrationsDown(migrationDir *string) error {
	// run migrations
	m, err := migrate.New(
		fmt.Sprintf("file://%s", *migrationDir), adapter.postgres.ConnectionString)

	if err != nil {
		logger.GetInstance().Fatal("failed to delete migrations", zap.Error(err))
		return err
	}

	if err := m.Down(); err != nil && err != migrate.ErrNoChange {
		logger.GetInstance().Fatal("an error occurred while syncing the database", zap.Error(err))
		return err
	}

	return nil
}

func (adapter *PostgreSQLAdapter) SaveUser(user *model.User) error {
	_, err := adapter.postgres.Conn.Exec(`INSERT INTO users(email, password) VALUES($1, $2)`, user.Email, user.Password)
	if err != nil {
		logger.GetInstance().Error("fail to save user", zap.Error(err))
		return err
	}
	return nil
}

func (adapter *PostgreSQLAdapter) GetUser(email string) (*model.User, bool, error) {
	row := adapter.postgres.Conn.QueryRow(`SELECT id, email, password FROM users WHERE email = $1 `, email)
	return mapRowsToUser(row)
}

func (adapter *PostgreSQLAdapter) AddFavourite(userID, propertyID int64) error {
	_, err := adapter.postgres.Conn.Exec(`INSERT INTO favourites(user_id, property_id) VALUES($1, $2) `, userID, propertyID)
	if err != nil {
		return err
	}
	return nil
}

func (adapter *PostgreSQLAdapter) ListFavourites(search users.FavouritesSearchParams) (*model.PropertiesPaging, error) {

	pagingResult := &model.PropertiesPaging{
		Page:     search.Page,
		PageSize: search.PageSize,
	}

	offset := search.PageSize * (search.Page - 1)

	query := fmt.Sprintf(`SELECT r.*, count(*) OVER() AS full_count FROM properties r 
	INNER JOIN favourites f ON  f.property_id = r.id 
	INNER JOIN users u ON u.id = f.user_id 
	WHERE u.id = %v AND status = 'ACTIVE' 
	ORDER BY updated_at DESC OFFSET %d LIMIT %d`, search.UserID, offset, search.PageSize)

	rows, err := adapter.postgres.Conn.Query(query)
	if err != nil {
		logger.GetInstance().Error("error listing favourites", zap.Error(err))
		return nil, err
	}

	for rows.Next() {
		property, count, err := mapRowsToProperty(rows)
		if err != nil {
			return nil, err
		}
		pagingResult.Data = append(pagingResult.Data, property)
		pagingResult.Total = count
	}

	return pagingResult, nil

}

func mapRowsToUser(row *sql.Row) (*model.User, bool, error) {
	if row.Err() != nil {
		return nil, false, row.Err()
	}

	var email, password string
	var id int64
	err := row.Scan(&id, &email, &password)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, false, nil
		}
		return nil, false, err
	}
	return &model.User{
		ID:       id,
		Email:    email,
		Password: password,
	}, true, nil

}

func (adapter *PostgreSQLAdapter) SaveProperty(property *model.Property) (*model.Property, error) {
	row := adapter.postgres.Conn.QueryRow(`INSERT INTO properties(title, description, longitude, latitude, sale_price, administrative_fee, property_type,  bedrooms, bathrooms, parking_spots, area, photos, status) 
			VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13) RETURNING *`, property.Title, property.Description, property.Location.Longitude, property.Location.Latitude,
		property.Pricing.SalePrice, property.Pricing.AdministrativeFee, property.PropertyType, property.Bedrooms, property.Bathrooms, property.ParkingSpots, property.Area, pq.Array(property.Photos), property.Status)

	propertyStored, found, err := mapRowToProperty(row)
	if err != nil {
		return nil, err
	}
	if !found {
		return nil, errors.New("property not found")
	}
	return propertyStored, nil
}

func (adapter *PostgreSQLAdapter) UpdateProperty(property *model.Property) (*model.Property, error) {
	row := adapter.postgres.Conn.QueryRow(`UPDATE properties SET  
                      title = $2, 
                      description = $3, 
                      longitude = $4, 
                      latitude = $5, 
                      sale_price = $6, 
                      administrative_fee = $7, 
                      property_type = $8,  
                      bedrooms = $9, 
                      bathrooms = $10, 
                      parking_spots = $11, 
                      area = $12, 
                      photos = $13, 
                      status = $14 WHERE id = $1
			RETURNING *`, property.ID, property.Title, property.Description, property.Location.Longitude, property.Location.Latitude,
		property.Pricing.SalePrice, property.Pricing.AdministrativeFee, property.PropertyType, property.Bedrooms, property.Bathrooms, property.ParkingSpots, property.Area, pq.Array(property.Photos), property.Status)

	propertyStored, found, err := mapRowToProperty(row)
	if err != nil {
		return nil, err
	}
	if !found {
		return nil, errors.New("property not found")
	}
	return propertyStored, nil
}

func (adapter *PostgreSQLAdapter) GetProperty(propertyID int64) (*model.Property, bool, error) {
	row := adapter.postgres.Conn.QueryRow(`SELECT * FROM properties WHERE id = $1`, propertyID)
	return mapRowToProperty(row)
}

func (adapter *PostgreSQLAdapter) FilterProperties(search properties.PropertySearchParams) (*model.PropertiesPaging, error) {
	pagingResult := &model.PropertiesPaging{
		Page:     search.Page,
		PageSize: search.PageSize,
	}

	whereClause := ""

	if search.Status != "ALL" {
		whereClause += fmt.Sprintf(" WHERE (status = '%s') ", search.Status)
	}

	if search.Bbox != nil {
		bboxClause := fmt.Sprintf(" (latitude >= %v AND latitude <= %v AND longitude >= %v AND longitude <= %v) ", search.Bbox.MinLatitude, search.Bbox.MaxLatitude,
			search.Bbox.MinLongitude, search.Bbox.MaxLongitude)

		if whereClause != "" {
			whereClause += " AND " + bboxClause
		} else {
			whereClause += " WHERE " + bboxClause
		}
	}

	offset := search.PageSize * (search.Page - 1)

	query := fmt.Sprintf(`SELECT *, count(*) OVER() AS full_count FROM properties %s ORDER BY updated_at DESC OFFSET %d LIMIT %d`, whereClause, offset, search.PageSize)

	rows, err := adapter.postgres.Conn.Query(query)
	if err != nil {
		logger.GetInstance().Error("error executing filtering properties query", zap.Error(err))
		return nil, err
	}

	for rows.Next() {
		property, count, err := mapRowsToProperty(rows)
		if err != nil {
			return nil, err
		}
		pagingResult.Data = append(pagingResult.Data, property)
		pagingResult.Total = count
	}

	return pagingResult, nil
}

func mapRowsToProperty(rows *sql.Rows) (*model.Property, int64, error) {
	var title, propertyType, status string
	var description sql.NullString
	var longitude, latitude float64
	var bedrooms, bathrooms, area, salePrice int
	var parkingSpots sql.NullInt64
	var id int64
	var createdAt, updateAt time.Time
	var photos []string
	var administrativeFee sql.NullInt64
	var fullCount int64

	err := rows.Scan(&id, &title, &description, &longitude, &latitude, &salePrice, &administrativeFee, &propertyType, &bedrooms, &bathrooms,
		&parkingSpots, &area, pq.Array(&photos), &status, &createdAt, &updateAt, &fullCount)
	if err != nil {
		logger.GetInstance().Error("error mapping property rows", zap.Error(err))
		return nil, 0, err
	}

	var descriptionValue model.Description
	if description.Valid {
		v := string(description.String)
		descriptionValue = &v
	}

	var parkingSpotsValue model.ParkingSpots
	if parkingSpots.Valid {
		v := int(parkingSpots.Int64)
		parkingSpotsValue = &v
	}

	var administrativeFeeValue model.AdministrativeFee
	if administrativeFee.Valid {
		v := int(administrativeFee.Int64)
		administrativeFeeValue = &v
	}

	return &model.Property{
		ID:          id,
		Title:       title,
		Description: descriptionValue,
		Location: model.Location{
			Longitude: longitude,
			Latitude:  latitude,
		},
		Pricing: model.Pricing{
			SalePrice:         salePrice,
			AdministrativeFee: administrativeFeeValue,
		},
		PropertyType: model.PropertyType(propertyType),
		Bedrooms:     bedrooms,
		Bathrooms:    bathrooms,
		ParkingSpots: parkingSpotsValue,
		Area:         area,
		Photos:       photos,
		CreatedAt:    createdAt,
		UpdatedAt:    updateAt,
		Status:       model.PropertyStatus(status),
	}, fullCount, nil

}

func mapRowToProperty(row *sql.Row) (*model.Property, bool, error) {
	if row.Err() != nil {
		logger.GetInstance().Error("error executing operation on property ", zap.Error(row.Err()))
		return nil, false, row.Err()
	}
	var title, propertyType, status string
	var longitude, latitude float64
	var bedrooms, bathrooms, area, salePrice int
	var parkingSpots sql.NullInt64
	var description sql.NullString
	var id int64
	var createdAt, updateAt time.Time
	var photos []string
	var administrativeFee sql.NullInt64

	err := row.Scan(&id, &title, &description, &longitude, &latitude, &salePrice, &administrativeFee, &propertyType, &bedrooms, &bathrooms,
		&parkingSpots, &area, pq.Array(&photos), &status, &createdAt, &updateAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, false, nil
		}
		logger.GetInstance().Error("error scanning property row", zap.Error(err))
		return nil, false, err
	}

	var descriptionValue model.Description
	if description.Valid {
		v := string(description.String)
		descriptionValue = &v
	}

	var parkingSpotsValue model.ParkingSpots
	if parkingSpots.Valid {
		v := int(parkingSpots.Int64)
		parkingSpotsValue = &v
	}

	var administrativeFeeValue model.AdministrativeFee
	if administrativeFee.Valid {
		v := int(administrativeFee.Int64)
		administrativeFeeValue = &v
	}

	return &model.Property{
		ID:          id,
		Title:       title,
		Description: descriptionValue,
		Location: model.Location{
			Longitude: longitude,
			Latitude:  latitude,
		},
		Pricing: model.Pricing{
			SalePrice:         salePrice,
			AdministrativeFee: administrativeFeeValue,
		},
		PropertyType: model.PropertyType(propertyType),
		Bedrooms:     bedrooms,
		Bathrooms:    bathrooms,
		ParkingSpots: parkingSpotsValue,
		Area:         area,
		Photos:       photos,
		CreatedAt:    createdAt,
		UpdatedAt:    updateAt,
		Status:       model.PropertyStatus(status),
	}, true, nil

}
