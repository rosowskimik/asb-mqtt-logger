package database

import (
	"database/sql"
	"errors"
	"fmt"
	"log/slog"

	_ "github.com/mattn/go-sqlite3"

	"rpi-data-logger/model"
)

// Error for data from unknown sensor
var ErrUnknownSensor = errors.New("sensor ID not found in the database")

// DB wraps the sql.DB connection pool.
type DB struct {
	Conn *sql.DB
}

// Connect to the SQLite DB and ensures the schema is created
func New(dbPath string) (*DB, error) {
	conn, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// SQL for creating the sensors metadata table
	createSensorsTableSQL := `
    CREATE TABLE IF NOT EXISTS sensors (
        sensor_id TEXT PRIMARY KEY,
        name TEXT NOT NULL,
        created_at DATETIME DEFAULT CURRENT_TIMESTAMP
    );`

	// SQL for creating the readings table
	createReadingsTableSQL := `
    CREATE TABLE IF NOT EXISTS sensor_readings (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        timestamp DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
        sensor_id TEXT NOT NULL,
        temp_min REAL, temp_max REAL, temp_median REAL,
        pressure_min REAL, pressure_max REAL, pressure_median REAL,
        humidity_min REAL, humidity_max REAL, humidity_median REAL,
        FOREIGN KEY(sensor_id) REFERENCES sensors(sensor_id)
    );`

	if _, err := conn.Exec(createSensorsTableSQL); err != nil {
		return nil, fmt.Errorf("failed to create sensors table: %w", err)
	}
	if _, err := conn.Exec(createReadingsTableSQL); err != nil {
		return nil, fmt.Errorf("failed to create readings table: %w", err)
	}

	slog.Info("Database is ready", "path", dbPath)
	return &DB{Conn: conn}, nil
}

// Validate the sensor ID and insert a new reading
func (db *DB) InsertReading(sensorID string, data model.SensorData) error {
	tx, err := db.Conn.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// 1. Check if the sensor is registered
	var exists int
	err = tx.QueryRow("SELECT 1 FROM sensors WHERE sensor_id = ?", sensorID).Scan(&exists)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrUnknownSensor // Return our specific error
		}
		return fmt.Errorf("failed to query for sensor: %w", err)
	}

	// 2. Insert the reading
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
