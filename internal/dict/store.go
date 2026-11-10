package dict

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3" // registers the "sqlite3" driver
)

// SchemaVersion is bumped whenever the krdict.db layout changes, so older
// databases can be detected and re-imported.
const SchemaVersion = 1

const schema = `
CREATE TABLE meta (
  key   TEXT PRIMARY KEY,
  value TEXT NOT NULL
);
CREATE TABLE entry (
  id            INTEGER PRIMARY KEY,
  krdict_id     INTEGER NOT NULL,   -- LexicalEntry@val (not unique)
  headword      TEXT NOT NULL,
  homonym       INTEGER,
  unit          TEXT,               -- 단어, 관용구, 속담, ...
  pos           TEXT,
  level         TEXT,
  hanja         TEXT,
  pronunciation TEXT
);
CREATE TABLE sense (
  entry_id INTEGER NOT NULL REFERENCES entry(id),
  ord      INTEGER NOT NULL,
  def_ko   TEXT,
  gloss_en TEXT,
  def_en   TEXT
);
CREATE TABLE form (                -- listed conjugated/contracted forms
  entry_id INTEGER NOT NULL REFERENCES entry(id),
  written  TEXT NOT NULL
);
`

const indexes = `
CREATE INDEX entry_headword ON entry(headword);
CREATE INDEX sense_entry    ON sense(entry_id);
CREATE INDEX form_written   ON form(written);
`

// Writer builds a new krdict.db. Call Insert for every entry, then Close.
type Writer struct {
	db                   *sql.DB
	tx                   *sql.Tx
	entry, sense, form   *sql.Stmt
	Entries, Senses, Eng int // counts, for the import summary
	Forms                int
}

// Create makes a new, empty krdict.db at path. The file must not exist.
func Create(path string) (*Writer, error) {
	db, err := sql.Open("sqlite3", "file:"+path+"?_journal_mode=OFF&_synchronous=OFF")
	if err != nil {
		return nil, err
	}
	w := &Writer{db: db}
	if _, err := db.Exec(schema); err != nil {
		db.Close()
		return nil, fmt.Errorf("dict: create schema: %w", err)
	}
	if err := w.begin(); err != nil {
		db.Close()
		return nil, err
	}
	return w, nil
}

func (w *Writer) begin() error {
	var err error
	if w.tx, err = w.db.Begin(); err != nil {
		return err
	}
	if w.entry, err = w.tx.Prepare(`INSERT INTO entry
		(krdict_id, headword, homonym, unit, pos, level, hanja, pronunciation)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)`); err != nil {
		return err
	}
	if w.sense, err = w.tx.Prepare(`INSERT INTO sense
		(entry_id, ord, def_ko, gloss_en, def_en) VALUES (?, ?, ?, ?, ?)`); err != nil {
		return err
	}
	w.form, err = w.tx.Prepare(`INSERT INTO form (entry_id, written) VALUES (?, ?)`)
	return err
}

// SetMeta records a key/value pair, e.g. the source label or import time.
func (w *Writer) SetMeta(key, value string) error {
	_, err := w.tx.Exec(`INSERT OR REPLACE INTO meta (key, value) VALUES (?, ?)`, key, value)
	return err
}

// Insert adds one entry with its senses and forms.
func (w *Writer) Insert(e Entry) error {
	res, err := w.entry.Exec(e.KrdictID, e.Headword, e.Homonym, e.Unit,
		nullable(e.POS), e.Level, nullable(e.Hanja), nullable(e.Pronunciation))
	if err != nil {
		return fmt.Errorf("dict: insert %s: %w", e.Headword, err)
	}
	id, err := res.LastInsertId()
	if err != nil {
		return err
	}
	for _, s := range e.Senses {
		if _, err := w.sense.Exec(id, s.Order, s.DefKo, nullable(s.GlossEn), nullable(s.DefEn)); err != nil {
			return fmt.Errorf("dict: insert sense of %s: %w", e.Headword, err)
		}
		w.Senses++
		if s.GlossEn != "" || s.DefEn != "" { // ~0.6% have a definition but no short gloss
			w.Eng++
		}
	}
	for _, f := range e.Forms {
		if _, err := w.form.Exec(id, f); err != nil {
			return fmt.Errorf("dict: insert form of %s: %w", e.Headword, err)
		}
		w.Forms++
	}
	w.Entries++
	return nil
}

// Close commits everything, builds the indexes and closes the database.
func (w *Writer) Close() error {
	defer w.db.Close()
	if err := w.SetMeta("schema_version", fmt.Sprint(SchemaVersion)); err != nil {
		return err
	}
	if err := w.tx.Commit(); err != nil {
		return err
	}
	if _, err := w.db.Exec(indexes); err != nil {
		return fmt.Errorf("dict: create indexes: %w", err)
	}
	_, err := w.db.Exec(`ANALYZE`)
	return err
}

// Abort discards the import.
func (w *Writer) Abort() {
	w.tx.Rollback()
	w.db.Close()
}

func nullable(s string) any {
	if s == "" {
		return nil
	}
	return s
}

// Store reads an existing krdict.db.
type Store struct{ db *sql.DB }

// Open opens krdict.db read-only.
func Open(path string) (*Store, error) {
	db, err := sql.Open("sqlite3", "file:"+path+"?mode=ro&immutable=1")
	if err != nil {
		return nil, err
	}
	var v string
	if err := db.QueryRow(`SELECT value FROM meta WHERE key = 'schema_version'`).Scan(&v); err != nil {
		db.Close()
		return nil, fmt.Errorf("dict: %s is not a Hannote dictionary (re-run hannote-import): %w", path, err)
	}
	if v != fmt.Sprint(SchemaVersion) {
		db.Close()
		return nil, fmt.Errorf("dict: %s has schema %s, want %d (re-run hannote-import)", path, v, SchemaVersion)
	}
	return &Store{db: db}, nil
}

// Close closes the database.
func (s *Store) Close() error { return s.db.Close() }

// Meta returns a value recorded at import time, or "".
func (s *Store) Meta(key string) string {
	var v string
	s.db.QueryRow(`SELECT value FROM meta WHERE key = ?`, key).Scan(&v)
	return v
}

// Lookup returns every entry whose headword is word, followed by entries
// that list word as a conjugated or contracted form. Homographs are ordered
// by homonym number.
func (s *Store) Lookup(word string) ([]Entry, error) {
	rows, err := s.db.Query(`
		SELECT id, krdict_id, headword, homonym, unit, pos, level, hanja, pronunciation, 0 AS by_form
		  FROM entry WHERE headword = ?1
		UNION
		SELECT e.id, e.krdict_id, e.headword, e.homonym, e.unit, e.pos, e.level, e.hanja, e.pronunciation, 1
		  FROM form f JOIN entry e ON e.id = f.entry_id
		 WHERE f.written = ?1 AND e.headword <> ?1
		ORDER BY by_form, homonym, id`, word)
	if err != nil {
		return nil, err
	}
	var out []Entry
	for rows.Next() {
		var e Entry
		var pos, hanja, pron sql.NullString
		var byForm int
		if err := rows.Scan(&e.ID, &e.KrdictID, &e.Headword, &e.Homonym, &e.Unit,
			&pos, &e.Level, &hanja, &pron, &byForm); err != nil {
			rows.Close()
			return nil, err
		}
		e.POS, e.Hanja, e.Pronunciation = pos.String, hanja.String, pron.String
		out = append(out, e)
	}
	rows.Close()
	if err := rows.Err(); err != nil {
		return nil, err
	}
	for i := range out {
		if out[i].Senses, err = s.senses(out[i].ID); err != nil {
			return nil, err
		}
	}
	return out, nil
}

func (s *Store) senses(entryID int64) ([]Sense, error) {
	rows, err := s.db.Query(`SELECT ord, def_ko, gloss_en, def_en FROM sense
		WHERE entry_id = ? ORDER BY ord`, entryID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []Sense
	for rows.Next() {
		var sn Sense
		var ko, gloss, def sql.NullString
		if err := rows.Scan(&sn.Order, &ko, &gloss, &def); err != nil {
			return nil, err
		}
		sn.DefKo, sn.GlossEn, sn.DefEn = ko.String, gloss.String, def.String
		out = append(out, sn)
	}
	return out, rows.Err()
}
