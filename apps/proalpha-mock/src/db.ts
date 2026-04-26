import { Database } from "bun:sqlite";

const db = new Database("proalpha.sqlite");

db.run(`
  CREATE TABLE IF NOT EXISTS materials (
    id INTEGER PRIMARY KEY,
    cis_code TEXT UNIQUE NOT NULL,
    name TEXT NOT NULL,
    unit TEXT NOT NULL,
    description TEXT,
    belimed_nr TEXT,
    tyco_nr TEXT
  )
`);

db.run(`
  CREATE TABLE IF NOT EXISTS bom_entries (
    id INTEGER PRIMARY KEY,
    cis_product_code TEXT NOT NULL,
    material_cis_code TEXT NOT NULL,
    qty REAL NOT NULL,
    unit TEXT NOT NULL,
    FOREIGN KEY (material_cis_code) REFERENCES materials(cis_code)
  )
`);

const insertMaterial = db.prepare(
  "INSERT OR IGNORE INTO materials (cis_code, name, unit, belimed_nr, tyco_nr) VALUES (?, ?, ?, ?, ?)"
);
const seedMaterials = db.transaction(() => {
  insertMaterial.run("TBD-18301", "Aufnahmegehäuse 6.3mm 2-pol transparent", "ks", "18301", "172157-1");
  insertMaterial.run("2601086",   "MATE-N-LOK Buchse 0.30-0.89mm²",          "ks", "54770", "170362-1");
  insertMaterial.run("TBD-18286", "Aderendhülse Ø0.5mm×14mm DIN 46228",      "ks", "18286", "966156-1");
  insertMaterial.run("5100560",   "Kabel 0.5mm² (schwarz/braun)",             "m",  null,    null);
  insertMaterial.run("9100008",   "Sáček",                                    "ks", null,    null);
});
seedMaterials();

const insertBomEntry = db.prepare(
  "INSERT OR IGNORE INTO bom_entries (cis_product_code, material_cis_code, qty, unit) VALUES (?, ?, ?, ?)"
);
const seedBomEntries = db.transaction(() => {
  insertBomEntry.run("74678", "TBD-18301", 1,  "ks");
  insertBomEntry.run("74678", "2601086",   4,  "ks");
  insertBomEntry.run("74678", "TBD-18286", 4,  "ks");
  insertBomEntry.run("74678", "5100560",   1,  "m");
  insertBomEntry.run("74678", "9100008",   1,  "ks");
});
seedBomEntries();

export default db;
