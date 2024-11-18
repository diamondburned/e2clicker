-- Language: postgresql
--
CREATE TABLE IF NOT EXISTS delivery_methods (
  id text PRIMARY KEY, units text NOT NULL, name text NOT NULL
);

ALTER TABLE delivery_methods
  ADD COLUMN IF NOT EXISTS description text NOT NULL DEFAULT '';

INSERT INTO delivery_methods (id, units, name, description)
  VALUES --
    ('EB im', 'mg', 'Estradiol Benzoate, Intramuscular', ''),
    ('EV im', 'mg', 'Estradiol Valerate, Intramuscular', ''),
    ('EEn im', 'mg', 'Estradiol Enanthate, Intramuscular', ''),
    ('EC im', 'mg', 'Estradiol Cypionate, Intramuscular', ''),
    ('EUn im', 'mg', 'Estradiol Undecylate, Intramuscular', ''),
    ('EUn casubq', 'mg', 'Estradiol Undecylate in Castor oil, Subcutaneous', ''),
    ('patch ow', 'mcg/day', 'Patch (once weekly)', ''),
    ('patch tw', 'mcg/day', 'Patch (twice weekly)', '')
  ON CONFLICT (id)
    DO UPDATE SET
      units = excluded.units, name = excluded.name, description = excluded.description;
