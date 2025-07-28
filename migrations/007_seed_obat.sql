-- +migrate Up

INSERT INTO obat (nama_obat, deskripsi, stok, harga, satuan, created_by)
VALUES 
-- Obat Umum
('Paracetamol', 'Obat penurun demam dan nyeri', 500, 1500, 'tablet', 'seeder'),
('Ibuprofen', 'Anti inflamasi non steroid', 300, 2000, 'tablet', 'seeder'),
('Asam Mefenamat', 'Pereda nyeri haid dan nyeri ringan', 250, 1800, 'tablet', 'seeder'),
('Amoxicillin', 'Antibiotik golongan penisilin', 300, 2500, 'kapsul', 'seeder'),
('Ciprofloxacin', 'Antibiotik spektrum luas', 200, 3000, 'tablet', 'seeder'),
('Dexamethasone', 'Anti-inflamasi dan alergi', 100, 1000, 'tablet', 'seeder'),
('Loperamide', 'Obat diare', 150, 1500, 'tablet', 'seeder'),
('Ranitidine', 'Obat maag H2 blocker', 180, 1700, 'tablet', 'seeder'),
('Antasida Doen', 'Obat maag dan gangguan lambung', 200, 2000, 'tablet kunyah', 'seeder'),
('Metformin', 'Obat diabetes tipe 2', 100, 3500, 'tablet', 'seeder'),

-- Obat Batuk/Pilek
('Ambroxol', 'Ekspektoran untuk batuk berdahak', 100, 2200, 'sirup', 'seeder'),
('Salbutamol', 'Bronkodilator untuk asma', 150, 3000, 'tablet', 'seeder'),
('CTM', 'Antihistamin untuk alergi dan flu', 250, 1200, 'tablet', 'seeder'),
('Pseudoefedrin', 'Dekongestan hidung', 100, 2000, 'tablet', 'seeder'),
('Dekstrometorfan', 'Antitusif untuk batuk kering', 130, 1600, 'tablet', 'seeder'),

-- Obat Hipertensi
('Amlodipine', 'Obat tekanan darah tinggi', 300, 2700, 'tablet', 'seeder'),
('Captopril', 'ACE inhibitor', 250, 2400, 'tablet', 'seeder'),
('Bisoprolol', 'Beta blocker', 180, 2900, 'tablet', 'seeder'),

-- Obat Kolesterol
('Simvastatin', 'Obat penurun kolesterol', 150, 3000, 'tablet', 'seeder'),
('Atorvastatin', 'Obat penurun kolesterol', 200, 3500, 'tablet', 'seeder'),

-- Obat Pencernaan
('Domperidone', 'Obat mual dan muntah', 160, 2000, 'tablet', 'seeder'),
('Omeprazole', 'PPI untuk tukak lambung', 220, 2700, 'tablet', 'seeder'),

-- Obat Kulit
('Betametason', 'Kortikosteroid topikal', 100, 5000, 'salep', 'seeder'),
('Miconazole', 'Antijamur topikal', 120, 4800, 'salep', 'seeder'),
('Asiklovir', 'Obat herpes topikal', 90, 6000, 'salep', 'seeder'),

-- Obat Mata
('Tetrahydrozoline', 'Tetes mata merah', 80, 3000, 'tetes mata', 'seeder'),
('Chloramphenicol', 'Antibiotik tetes mata', 100, 3500, 'tetes mata', 'seeder'),

-- Obat Anak
('Paracetamol Drops', 'Penurun demam bayi', 90, 4000, 'drops', 'seeder'),
('Ibuprofen Sirup', 'Penurun demam anak', 70, 4200, 'sirup', 'seeder'),

-- Vitamin
('Vitamin C 500mg', 'Suplemen harian', 300, 1200, 'tablet', 'seeder'),
('Vitamin B Complex', 'Suplemen B kompleks', 250, 1500, 'tablet', 'seeder'),
('Zinc', 'Untuk daya tahan tubuh', 150, 1700, 'tablet', 'seeder'),

-- Lain-lain
('ORS', 'Oralit untuk diare', 200, 800, 'sachet', 'seeder'),
('Tranexamic Acid', 'Penghenti perdarahan', 120, 2700, 'tablet', 'seeder'),
('Ketoprofen', 'Obat nyeri sendi/topikal', 100, 6000, 'patch', 'seeder'),
('Clobetasol', 'Steroid topikal poten', 90, 7000, 'salep', 'seeder'),
('Albendazole', 'Obat cacing', 180, 1800, 'tablet', 'seeder'),
('Neurobion', 'Vitamin neurotropik', 140, 3500, 'tablet', 'seeder'),
('Magnesium Hydroxide', 'Laksatif/osmotik', 130, 2000, 'sirup', 'seeder'),
('Antalgin', 'Pereda nyeri', 160, 1800, 'tablet', 'seeder'),
('Diazepam', 'Obat penenang', 80, 4000, 'tablet', 'seeder'),
('Lansoprazole', 'PPI alternatif', 70, 3800, 'tablet', 'seeder'),
('Levofloxacin', 'Antibiotik generasi baru', 100, 5000, 'tablet', 'seeder'),
('Doxycycline', 'Antibiotik tetrasiklin', 90, 4500, 'kapsul', 'seeder'),
('Furosemide', 'Diuretik untuk edema', 110, 2500, 'tablet', 'seeder'),
('Allopurinol', 'Obat asam urat', 130, 3000, 'tablet', 'seeder'),
('Clemastine', 'Antihistamin', 100, 2200, 'tablet', 'seeder');

-- +migrate Down

DELETE FROM obat;