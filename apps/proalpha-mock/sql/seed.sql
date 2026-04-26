-- Seed data from drawing 74678
-- cis_codes marked TBD-* are placeholders; fill in once CiS assigns real codes

INSERT OR IGNORE INTO materials (cis_code, name, unit, belimed_nr, tyco_nr) VALUES
  ('TBD-18301', 'Aufnahmegehäuse 6.3mm 2-pol transparent', 'ks', '18301', '172157-1'),
  ('2601086',   'MATE-N-LOK Buchse 0.30-0.89mm²',          'ks', '54770', '170362-1'),
  ('TBD-18286', 'Aderendhülse Ø0.5mm×14mm DIN 46228',      'ks', '18286', '966156-1'),
  ('5100560',   'Kabel 0.5mm² (schwarz/braun)',             'm',  NULL,    NULL),
  ('9100008',   'Sáček',                                    'ks', NULL,    NULL);

INSERT OR IGNORE INTO bom_entries (cis_product_code, material_cis_code, qty, unit) VALUES
  ('74678', 'TBD-18301', 1,  'ks'),
  ('74678', '2601086',   4,  'ks'),
  ('74678', 'TBD-18286', 4,  'ks'),
  ('74678', '5100560',   1,  'm'),
  ('74678', '9100008',   1,  'ks');
