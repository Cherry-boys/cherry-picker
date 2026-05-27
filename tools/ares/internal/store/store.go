// Hand-authored store for ares-pp-cli transcendence features.
// Not generator-emitted; survives regen as a whole file.
package store

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	_ "modernc.org/sqlite"
)

// Store is the local SQLite layer that the stateless ARES API has no notion of.
// It backs portfolio watchlists, subject snapshots (for change diffing), and
// offline full-text search.
type Store struct {
	db *sql.DB
}

// Subject is a normalized, high-gravity view of an ARES economic subject.
type Subject struct {
	Ico           string          `json:"ico"`
	ObchodniJmeno string          `json:"obchodniJmeno"`
	Adresa        string          `json:"adresa"`
	PravniForma   string          `json:"pravniForma"`
	Raw           json.RawMessage `json:"raw,omitempty"`
	FetchedAt     time.Time       `json:"fetchedAt"`
}

// PortfolioEntry is a watchlisted IČO with an optional human label.
type PortfolioEntry struct {
	Ico     string    `json:"ico"`
	Label   string    `json:"label,omitempty"`
	AddedAt time.Time `json:"addedAt"`
}

// Change describes one detected difference between two snapshots of a subject.
type Change struct {
	Ico        string    `json:"ico"`
	Field      string    `json:"field"`
	Old        string    `json:"old"`
	New        string    `json:"new"`
	DetectedAt time.Time `json:"detectedAt"`
}

// DefaultDBPath returns the per-user SQLite path, honoring ARES_DB and XDG.
func DefaultDBPath() string {
	if p := os.Getenv("ARES_DB"); p != "" {
		return p
	}
	base, err := os.UserConfigDir()
	if err != nil || base == "" {
		home, _ := os.UserHomeDir()
		base = filepath.Join(home, ".config")
	}
	return filepath.Join(base, "ares-pp-cli", "ares.db")
}

// Open opens (and migrates) the store at path; "" uses DefaultDBPath.
func Open(ctx context.Context, path string) (*Store, error) {
	if path == "" {
		path = DefaultDBPath()
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return nil, fmt.Errorf("creating db dir: %w", err)
	}
	// busy_timeout so concurrent invocations (and the scorecard probe matrix
	// running commands in parallel) wait for the lock instead of failing with
	// SQLITE_BUSY. WAL is intentionally NOT forced: switching journal mode
	// takes a brief exclusive lock that itself contends across parallel
	// processes, which is the very failure we're avoiding.
	dsn := path + "?_pragma=busy_timeout(10000)"
	db, err := sql.Open("sqlite", dsn)
	if err != nil {
		return nil, fmt.Errorf("opening db: %w", err)
	}
	db.SetMaxOpenConns(1)
	s := &Store{db: db}
	if err := s.migrate(ctx); err != nil {
		_ = db.Close()
		return nil, err
	}
	return s, nil
}

func (s *Store) Close() error { return s.db.Close() }

func (s *Store) migrate(ctx context.Context) error {
	stmts := []string{
		`CREATE TABLE IF NOT EXISTS subjects (
			ico TEXT PRIMARY KEY,
			obchodni_jmeno TEXT NOT NULL DEFAULT '',
			adresa TEXT NOT NULL DEFAULT '',
			pravni_forma TEXT NOT NULL DEFAULT '',
			raw TEXT NOT NULL DEFAULT '{}',
			fetched_at TEXT NOT NULL
		)`,
		`CREATE TABLE IF NOT EXISTS snapshots (
			ico TEXT NOT NULL,
			obchodni_jmeno TEXT NOT NULL DEFAULT '',
			adresa TEXT NOT NULL DEFAULT '',
			pravni_forma TEXT NOT NULL DEFAULT '',
			captured_at TEXT NOT NULL
		)`,
		`CREATE INDEX IF NOT EXISTS idx_snapshots_ico ON snapshots(ico, captured_at)`,
		`CREATE TABLE IF NOT EXISTS portfolio (
			ico TEXT PRIMARY KEY,
			label TEXT NOT NULL DEFAULT '',
			added_at TEXT NOT NULL
		)`,
		`CREATE VIRTUAL TABLE IF NOT EXISTS subjects_fts USING fts5(
			ico, obchodni_jmeno, adresa
		)`,
	}
	for _, q := range stmts {
		if _, err := s.db.ExecContext(ctx, q); err != nil {
			return fmt.Errorf("migrate: %w", err)
		}
	}
	return nil
}

// UpsertSubject stores/updates a subject, appends a snapshot if a fingerprint
// field changed since the last snapshot, and refreshes the FTS row.
func (s *Store) UpsertSubject(ctx context.Context, sub Subject) error {
	if sub.FetchedAt.IsZero() {
		sub.FetchedAt = time.Now().UTC()
	}
	raw := string(sub.Raw)
	if raw == "" {
		raw = "{}"
	}
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback() }()

	var prevJmeno, prevAdresa, prevForma sql.NullString
	_ = tx.QueryRowContext(ctx,
		`SELECT obchodni_jmeno, adresa, pravni_forma FROM subjects WHERE ico = ?`, sub.Ico,
	).Scan(&prevJmeno, &prevAdresa, &prevForma)

	if _, err := tx.ExecContext(ctx,
		`INSERT INTO subjects (ico, obchodni_jmeno, adresa, pravni_forma, raw, fetched_at)
		 VALUES (?,?,?,?,?,?)
		 ON CONFLICT(ico) DO UPDATE SET
		   obchodni_jmeno=excluded.obchodni_jmeno,
		   adresa=excluded.adresa,
		   pravni_forma=excluded.pravni_forma,
		   raw=excluded.raw,
		   fetched_at=excluded.fetched_at`,
		sub.Ico, sub.ObchodniJmeno, sub.Adresa, sub.PravniForma, raw, sub.FetchedAt.Format(time.RFC3339),
	); err != nil {
		return err
	}

	changed := !prevJmeno.Valid ||
		prevJmeno.String != sub.ObchodniJmeno ||
		prevAdresa.String != sub.Adresa ||
		prevForma.String != sub.PravniForma
	if changed {
		if _, err := tx.ExecContext(ctx,
			`INSERT INTO snapshots (ico, obchodni_jmeno, adresa, pravni_forma, captured_at)
			 VALUES (?,?,?,?,?)`,
			sub.Ico, sub.ObchodniJmeno, sub.Adresa, sub.PravniForma, sub.FetchedAt.Format(time.RFC3339),
		); err != nil {
			return err
		}
	}

	if _, err := tx.ExecContext(ctx, `DELETE FROM subjects_fts WHERE ico = ?`, sub.Ico); err != nil {
		return err
	}
	if _, err := tx.ExecContext(ctx,
		`INSERT INTO subjects_fts (ico, obchodni_jmeno, adresa) VALUES (?,?,?)`,
		sub.Ico, sub.ObchodniJmeno, sub.Adresa,
	); err != nil {
		return err
	}
	return tx.Commit()
}

// GetSubject returns the stored subject for an IČO, or ok=false if absent.
func (s *Store) GetSubject(ctx context.Context, ico string) (Subject, bool, error) {
	var sub Subject
	var raw string
	var fetched string
	err := s.db.QueryRowContext(ctx,
		`SELECT ico, obchodni_jmeno, adresa, pravni_forma, raw, fetched_at FROM subjects WHERE ico = ?`, ico,
	).Scan(&sub.Ico, &sub.ObchodniJmeno, &sub.Adresa, &sub.PravniForma, &raw, &fetched)
	if err == sql.ErrNoRows {
		return Subject{}, false, nil
	}
	if err != nil {
		return Subject{}, false, err
	}
	sub.Raw = json.RawMessage(raw)
	sub.FetchedAt, _ = time.Parse(time.RFC3339, fetched)
	return sub, true, nil
}

// Search runs an FTS5 query over name + address, newest first.
func (s *Store) Search(ctx context.Context, query string, limit int) ([]Subject, error) {
	if limit <= 0 {
		limit = 25
	}
	rows, err := s.db.QueryContext(ctx,
		`SELECT f.ico, s.obchodni_jmeno, s.adresa, s.pravni_forma, s.fetched_at
		 FROM subjects_fts f JOIN subjects s ON s.ico = f.ico
		 WHERE subjects_fts MATCH ? ORDER BY s.fetched_at DESC LIMIT ?`,
		query, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []Subject
	for rows.Next() {
		var sub Subject
		var fetched sql.NullString
		if err := rows.Scan(&sub.Ico, &sub.ObchodniJmeno, &sub.Adresa, &sub.PravniForma, &fetched); err != nil {
			continue
		}
		sub.FetchedAt, _ = time.Parse(time.RFC3339, fetched.String)
		out = append(out, sub)
	}
	return out, rows.Err()
}

// AddPortfolio adds or relabels a watchlisted IČO.
func (s *Store) AddPortfolio(ctx context.Context, ico, label string) error {
	_, err := s.db.ExecContext(ctx,
		`INSERT INTO portfolio (ico, label, added_at) VALUES (?,?,?)
		 ON CONFLICT(ico) DO UPDATE SET label=excluded.label`,
		ico, label, time.Now().UTC().Format(time.RFC3339))
	return err
}

// RemovePortfolio drops an IČO from the watchlist.
func (s *Store) RemovePortfolio(ctx context.Context, ico string) error {
	_, err := s.db.ExecContext(ctx, `DELETE FROM portfolio WHERE ico = ?`, ico)
	return err
}

// ListPortfolio returns all watchlisted IČOs, oldest first.
func (s *Store) ListPortfolio(ctx context.Context) ([]PortfolioEntry, error) {
	rows, err := s.db.QueryContext(ctx,
		`SELECT ico, label, added_at FROM portfolio ORDER BY added_at ASC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []PortfolioEntry
	for rows.Next() {
		var e PortfolioEntry
		var added sql.NullString
		if err := rows.Scan(&e.Ico, &e.Label, &added); err != nil {
			continue
		}
		e.AddedAt, _ = time.Parse(time.RFC3339, added.String)
		out = append(out, e)
	}
	return out, rows.Err()
}

// ChangesSince returns snapshot-to-snapshot field changes detected at or after
// the cutoff, across all subjects, newest first.
func (s *Store) ChangesSince(ctx context.Context, since time.Time) ([]Change, error) {
	rows, err := s.db.QueryContext(ctx,
		`SELECT ico, obchodni_jmeno, adresa, pravni_forma, captured_at
		 FROM snapshots ORDER BY ico ASC, captured_at ASC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	type snap struct {
		jmeno, adresa, forma string
		at                   time.Time
	}
	prev := map[string]snap{}
	var out []Change
	for rows.Next() {
		var ico string
		var jmeno, adresa, forma, at sql.NullString
		if err := rows.Scan(&ico, &jmeno, &adresa, &forma, &at); err != nil {
			continue
		}
		t, _ := time.Parse(time.RFC3339, at.String)
		cur := snap{jmeno.String, adresa.String, forma.String, t}
		if p, ok := prev[ico]; ok && !t.Before(since) {
			if p.jmeno != cur.jmeno {
				out = append(out, Change{ico, "obchodniJmeno", p.jmeno, cur.jmeno, t})
			}
			if p.adresa != cur.adresa {
				out = append(out, Change{ico, "adresa", p.adresa, cur.adresa, t})
			}
			if p.forma != cur.forma {
				out = append(out, Change{ico, "pravniForma", p.forma, cur.forma, t})
			}
		}
		prev[ico] = cur
	}
	return out, rows.Err()
}
