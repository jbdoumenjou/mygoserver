package db

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
)

type DBStructure struct {
	Chirps map[int]Chirp `json:"chirps"`
}

// DB is a simple file database.
type DB struct {
	path string
	data DBStructure
	mux  *sync.RWMutex
}

// NewDB creates a new database connection
// and creates the database file if it doesn't exist
func NewDB(path string) (*DB, error) {
	db := &DB{path: path, mux: &sync.RWMutex{}}
	if _, err := os.Stat(path); err != nil {
		if !os.IsNotExist(err) {
			return nil, fmt.Errorf("stat %s: %w", path, err)
		}
		if err := db.writeDB(DBStructure{Chirps: map[int]Chirp{}}); err != nil {
			return nil, fmt.Errorf("write db: %w", err)
		}
	}

	if err := db.loadDB(); err != nil {
		return nil, fmt.Errorf("load db: %w", err)
	}

	return db, nil
}

// Chirp is a single chirp.
type Chirp struct {
	ID   int    `json:"id"`
	Body string `json:"body"`
}

// CreateChirp creates a new chirp and saves it to disk
func (db *DB) CreateChirp(body string) (Chirp, error) {
	db.mux.Lock()
	defer db.mux.Unlock()

	id := len(db.data.Chirps) + 1
	chirp := Chirp{ID: id, Body: body}
	db.data.Chirps[id] = chirp
	if err := db.writeDB(db.data); err != nil {
		return Chirp{}, fmt.Errorf("write db: %w", err)
	}

	return chirp, nil
}

// ListChirps returns all chirps in the database
func (db *DB) ListChirps() ([]Chirp, error) {
	db.mux.RLock()
	defer db.mux.RUnlock()

	var chirps []Chirp
	for _, chirp := range db.data.Chirps {
		chirps = append(chirps, chirp)
	}

	return chirps, nil
}

// loadDB reads the database file into memory
func (db *DB) loadDB() error {
	db.mux.Lock()
	defer db.mux.Unlock()

	content, err := os.ReadFile(db.path)
	if err != nil {
		return fmt.Errorf("read db: %w", err)
	}

	if err = json.Unmarshal(content, &db.data); err != nil {
		return fmt.Errorf("unmarshal db: %w", err)
	}

	return nil
}

// writeDB writes the database file to disk
func (db *DB) writeDB(dbStructure DBStructure) error {
	data, err := json.Marshal(dbStructure)
	if err != nil {
		return fmt.Errorf("marshal db: %w", err)
	}

	if err = os.WriteFile(db.path, data, 0644); err != nil {
		return fmt.Errorf("write db: %w", err)
	}

	return nil
}
