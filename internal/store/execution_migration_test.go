package store

import (
	"database/sql"
	"io/fs"
	"path/filepath"
	"strings"
	"testing"
)

// A user may already have running work and correction history when this preview is updated.
func TestActionTimeMigrationPreservesRecordedWork(t *testing.T) {
	path := filepath.Join(t.TempDir(), "existing.db")
	db, err := sql.Open("sqlite", path+"?_pragma=foreign_keys(1)")
	if err != nil {
		t.Fatal(err)
	}
	db.SetMaxOpenConns(1)
	if _, err := db.Exec(`CREATE TABLE schema_migrations(name TEXT PRIMARY KEY)`); err != nil {
		t.Fatal(err)
	}
	names, err := fs.Glob(migrations, "migrations/*.sql")
	if err != nil {
		t.Fatal(err)
	}
	for _, name := range names {
		if strings.Contains(name, "0011_") {
			continue
		}
		if err := apply(db, name); err != nil {
			t.Fatal(err)
		}
	}
	const stamp = "2026-10-25T12:00:00.000000000Z"
	for _, query := range []string{
		`INSERT INTO tasks(id,title,created_by,created_at,updated_at,estimate_minutes) VALUES(1,'Existing task','user','` + stamp + `','` + stamp + `',2)`,
		`INSERT INTO steps(id,task_id,description,position) VALUES(1,1,'Existing step',0)`,
		`INSERT INTO focus_sessions(id,task_id,step_id,duration_minutes,state,segment_started_at,accumulated_ms,created_at) VALUES(1,1,1,15,'running','` + stamp + `',110000,'` + stamp + `')`,
		`INSERT INTO focus_selection VALUES(1,1)`,
		`INSERT INTO focus_requests VALUES('previous-start',1)`,
		`INSERT INTO focus_corrections(session_id,before_ms,after_ms,reason,actor,request_id,created_at) VALUES(1,200000,110000,'Existing correction','user','previous-correction','` + stamp + `')`,
	} {
		if _, err := db.Exec(query); err != nil {
			t.Fatal(err)
		}
	}
	db.Close()
	db, err = Open(path)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	var elapsed int64
	var started, state string
	if err := db.QueryRow(`SELECT accumulated_ms,segment_started_at,state FROM focus_sessions WHERE id=1`).Scan(&elapsed, &started, &state); err != nil {
		t.Fatal(err)
	}
	if elapsed != 110000 || started != stamp || state != "running" {
		t.Fatalf("migration changed work: %d %s %s", elapsed, started, state)
	}
	for _, table := range []string{"tasks", "steps", "focus_selection", "focus_requests", "focus_corrections"} {
		var count int
		if err := db.QueryRow(`SELECT COUNT(*) FROM ` + table).Scan(&count); err != nil || count != 1 {
			t.Fatalf("%s: count=%d error=%v", table, count, err)
		}
	}
	if _, err := db.Query(`SELECT duration_minutes FROM focus_sessions`); err == nil {
		t.Fatal("independent duration column still exists")
	}
	rows, err := db.Query(`PRAGMA foreign_key_check`)
	if err != nil {
		t.Fatal(err)
	}
	defer rows.Close()
	if rows.Next() {
		t.Fatal("migration broke existing references")
	}
}
