package database

import (
	"database/sql"
	"fmt"
	"log/slog"

	_ "modernc.org/sqlite"

	"rpi-data-logger/model"
)

// DB wraps the sql.DB connection pool.
type DB struct {
	Conn *sql.DB
}

// Connect to the SQLite DB and ensures the schema is created
func New(dbPath string) (*DB, error) {
	conn, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// SQL for creating the readings table
	createReadingsTableSQL := `
    CREATE TABLE IF NOT EXISTS sensor_readings (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        timestamp DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
        sensor_id INTEGER NOT NULL,
        temp_min REAL, temp_max REAL, temp_median REAL,
        pressure_min REAL, pressure_max REAL, pressure_median REAL,
        humidity_min REAL, humidity_max REAL, humidity_median REAL
    );`

	if _, err := conn.Exec(createReadingsTableSQL); err != nil {
		return nil, fmt.Errorf("failed to create readings table: %w", err)
	}

	slog.Info("Database is ready", "path", dbPath)
	return &DB{Conn: conn}, nil
}

// Validate the sensor ID and insert a new reading
func (db *DB) InsertReading(sensorID int, data model.SensorData) error {
	tx, err := db.Conn.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	_, err = tx.Exec(`
        INSERT INTO sensor_readings (
            sensor_id,
            temp_min, temp_max, temp_median,
            pressure_min, pressure_max, pressure_median,
            humidity_min, humidity_max, humidity_median
        ) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
    `,
		sensorID,
		data.Temp.Min, data.Temp.Max, data.Temp.Med,
		data.Pressure.Min, data.Pressure.Max, data.Pressure.Med,
		data.Humidity.Min, data.Humidity.Max, data.Humidity.Med,
	)
	if err != nil {
		return fmt.Errorf("failed to insert sensor reading: %w", err)
	}

	return tx.Commit()
}

// Gracefully close the database connection
func (db *DB) Close() {
	db.Conn.Close()
}
