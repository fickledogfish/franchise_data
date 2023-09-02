package db

import (
	"database/sql"

	"example.com/franchises/domain"
	_ "github.com/mattn/go-sqlite3"
)

const (
	SqliteDatabaseFileName = "./locations.sqlite"

	sqliteLocationTableName = "Locations"
)

type sqliteDb struct {
	db *sql.DB
}

func NewSqliteDb() (sqliteDb, error) {
	db, err := sql.Open("sqlite3", SqliteDatabaseFileName)
	if err != nil {
		return sqliteDb{}, err
	}

	_, err = db.Exec(
		"CREATE TABLE IF NOT EXISTS " + sqliteLocationTableName + "(" +
			"id          INTEGER PRIMARY KEY," +
			"origin      TEXT    NOT NULL," +
			"name        TEXT    NOT NULL," +
			"street      TEXT    NOT NULL," +
			"city        TEXT    NOT NULL," +
			"state       TEXT    NOT NULL," +
			"country     TEXT    NOT NULL," +
			"postal_code TEXT    NOT NULL" +
			")",
	)
	if err != nil {
		return sqliteDb{}, err
	}

	return sqliteDb{
		db: db,
	}, nil
}

func (self sqliteDb) SaveLocation(location domain.Location) error {
	_, err := self.db.Exec(
		"INSERT OR IGNORE INTO "+sqliteLocationTableName+"("+
			"id,"+
			"origin,"+
			"name,"+
			"street,"+
			"city,"+
			"state,"+
			"country,"+
			"postal_code"+
			") "+
			"VALUES(?, ?, ?, ?, ?, ?, ?, ?)",
		location.Id,
		location.Origin,
		location.Name,
		location.Address.Street,
		location.Address.City,
		location.Address.State,
		location.Address.Country,
		location.Address.PostalCode,
	)

	return err
}

func (self sqliteDb) GetSavedLocations() ([]domain.Location, error) {
	rows, err := self.db.Query("SELECT " +
		"id, origin, name, street, city, state, country, postal_code " +
		"FROM " + sqliteLocationTableName,
	)
	if err != nil {
		return []domain.Location{}, nil
	}

	var locations []domain.Location
	for rows.Next() {
		var row domain.Location
		err := rows.Scan(
			&row.Id,
			&row.Origin,
			&row.Name,
			&row.Address.Street,
			&row.Address.City,
			&row.Address.State,
			&row.Address.Country,
			&row.Address.PostalCode,
		)
		if err != nil {
			return []domain.Location{}, nil
		}

		locations = append(locations, row)
	}

	return locations, nil
}

func (self sqliteDb) GetSavedLocationsFrom(origin string) ([]domain.Location, error) {
	rows, err := self.db.Query("SELECT "+
		"id, origin, name, street, city, state, country, postal_code "+
		"FROM "+sqliteLocationTableName+
		" WHERE origin = ?",
		origin,
	)
	if err != nil {
		return []domain.Location{}, nil
	}

	var locations []domain.Location
	for rows.Next() {
		var row domain.Location
		err := rows.Scan(
			&row.Id,
			&row.Origin,
			&row.Name,
			&row.Address.Street,
			&row.Address.City,
			&row.Address.State,
			&row.Address.Country,
			&row.Address.PostalCode,
		)
		if err != nil {
			return []domain.Location{}, nil
		}

		locations = append(locations, row)
	}

	return locations, nil
}
