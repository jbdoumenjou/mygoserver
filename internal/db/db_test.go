package db

import (
	"encoding/json"
	"os"
	"reflect"
	"testing"
)

func TestDB_CreateChirp(t *testing.T) {
	const dbPath = "testdata/empty_database.json"

	db, err := NewDB(dbPath)
	if err != nil {
		t.Errorf("newDB should not have an error %v", err)
		return
	}
	defer os.Remove(dbPath)

	got, err := db.CreateChirp("I had something interesting for breakfast")
	if err != nil {
		t.Errorf("CreateChirp should not have an error %v", err)
		return
	}

	want := Chirp{ID: 0, Body: "I had something interesting for breakfast"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("CreateChirp() got = %v, want %v", got, want)
	}

	var content []byte
	content, err = os.ReadFile(db.path)
	if err != nil {
		t.Errorf("read db: %v", err)
	}

	var data DBStructure
	if err = json.Unmarshal(content, &data); err != nil {
		t.Errorf("unmarshal db: %v", err)
	}

	if !reflect.DeepEqual(data.Chirps, map[int]Chirp{0: want}) {
		t.Errorf("CreateChirp() got = %v, want %v", data, want)
	}
}

func TestDB_GetChirps(t *testing.T) {
	const dbPath = "testdata/database.json"

	db, err := NewDB(dbPath)
	if err != nil {
		t.Errorf("newDB should not have an error %v", err)
		return
	}

	got, err := db.ListChirps()
	if err != nil {
		t.Errorf("ListChirps should not have an error %v", err)
		return
	}

	want := []Chirp{
		{ID: 0, Body: "I had something interesting for breakfast"},
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("ListChirps() got = %v, want %v", got, want)
	}
}
