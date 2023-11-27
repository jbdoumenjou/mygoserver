package db

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"sync"
	"time"
)

type DBStructure struct {
	Chirps       map[int]Chirp        `json:"chirps"`
	Users        map[string]User      `json:"users"`
	RevokedToken map[string]time.Time `json:"revokedToken"`
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
			Chirps:       map[int]Chirp{},
			Users:        map[string]User{},
			RevokedToken: map[string]time.Time{},
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
	ID       int    `json:"id"`
	AuthorID int    `json:"author_id"`
	Body     string `json:"body"`
}

// CreateChirp creates a new chirp and saves it to disk
func (db *DB) CreateChirp(body string, authorID int) (Chirp, error) {
	db.mux.Lock()
	defer db.mux.Unlock()

	id := len(db.data.Chirps) + 1
	chirp := Chirp{ID: id, Body: body, AuthorID: authorID}
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

// DeleteChirp returns a single chirp.
func (db *DB) DeleteChirp(id int) {
	db.mux.RLock()
	defer db.mux.RUnlock()

	if _, ok := db.data.Chirps[id]; ok {
		delete(db.data.Chirps, id)
	}
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

// UpdateUser creates a new user and saves it to disk
func (db *DB) UpdateUser(id int, email, password string) (User, error) {
	db.mux.Lock()
	defer db.mux.Unlock()

	for _, user := range db.data.Users {
		if user.ID == id {
			if user.Email != email {
				delete(db.data.Users, email)
			}
			user.Email = email
			user.Password = password
			db.data.Users[email] = user
			if err := db.writeDB(db.data); err != nil {
				return User{}, fmt.Errorf("write db: %w", err)
			}
			return user, nil
		}
	}

	return User{}, errors.New("user not found")
}

// GetUserByEmail returns a single user.
func (db *DB) GetUserByEmail(email string) (*User, error) {
	db.mux.RLock()
	defer db.mux.RUnlock()

	user, ok := db.data.Users[email]
	if !ok {
		return nil, errors.New("not found")
	}

	return &user, nil
}

// GetUser returns a single user.
func (db *DB) GetUser(id int) (*User, error) {
	db.mux.RLock()
	defer db.mux.RUnlock()

	for _, user := range db.data.Users {
		if user.ID == id {
			return &user, nil
		}
	}

	return nil, errors.New("not found")
}

func (db *DB) RevokeToken(token string) string {
	db.mux.Lock()
	defer db.mux.Unlock()

	// naive approach, we could keep the previous revoked token
	db.data.RevokedToken[token] = time.Now().UTC()
	if err := db.writeDB(db.data); err != nil {
		return ""
	}

	return token
}

func (db *DB) IsTokenRevoked(token string) bool {
	db.mux.RLock()
	defer db.mux.RUnlock()

	_, ok := db.data.RevokedToken[token]
	return ok
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
