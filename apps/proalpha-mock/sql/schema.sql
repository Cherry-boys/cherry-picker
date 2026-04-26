CREATE TABLE IF NOT EXISTS materials (
  id INTEGER PRIMARY KEY,
  cis_code TEXT UNIQUE NOT NULL,
  name TEXT NOT NULL,
  unit TEXT NOT NULL,
  description TEXT,
  belimed_nr TEXT,
  tyco_nr TEXT
);

CREATE TABLE IF NOT EXISTS bom_entries (
  id INTEGER PRIMARY KEY,
  cis_product_code TEXT NOT NULL,
  material_cis_code TEXT NOT NULL,
  qty REAL NOT NULL,
  unit TEXT NOT NULL,
  FOREIGN KEY (material_cis_code) REFERENCES materials(cis_code)
);
