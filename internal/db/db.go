package db

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"sync"
)

type DBStructure struct {
	Chirps map[int]Chirp   `json:"chirps"`
	Users  map[string]User `json:"users"`
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
		structure := DBStructure{
			Chirps: map[int]Chirp{},
			Users:  map[string]User{},
		}
		if err := db.writeDB(structure); err != nil {
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

// GetChirp returns a single chirp.
func (db *DB) GetChirp(id int) (*Chirp, error) {
	db.mux.RLock()
	defer db.mux.RUnlock()

	chirp, ok := db.data.Chirps[id]
	if !ok {
		return nil, errors.New("not found")
	}

	return &chirp, nil
}

// User is a single user.
type User struct {
	ID       int    `json:"id"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

// CreateUser creates a new user and saves it to disk
func (db *DB) CreateUser(email, password string) (User, error) {
	db.mux.Lock()
	defer db.mux.Unlock()
	if _, ok := db.data.Users[email]; ok {
		return User{}, errors.New("user already exists")
	}

	id := len(db.data.Users) + 1

	user := User{
		ID:       id,
		Password: password,
		Email:    email,
	}
	db.data.Users[email] = user
	if err := db.writeDB(db.data); err != nil {
		return User{}, fmt.Errorf("write db: %w", err)
	}

	return user, nil
}

// GetUser returns a single user.
func (db *DB) GetUSer(email string) (*User, error) {
	db.mux.RLock()
	defer db.mux.RUnlock()

	for _, user := range db.data.Users {
		if user.Email == email {
			return &user, nil
		}
	}

	return nil, errors.New("not found")
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