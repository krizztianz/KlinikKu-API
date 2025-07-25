-- +migrate Up
-- Dokter 1: Dokter Umum (tanpa spesialisasi)
INSERT INTO dokter (
    no_izin_praktek, nama, tanggal_lahir, jenis_kelamin, alamat, no_hp, no_telepon, ktp, email, created_by
) VALUES (
    'SIP-0001', 'dr. Sari Yuliani', '1990-01-01', 'P',
    'Jl. Sehat No.1', '081234567891', '0211234567',
    '3175010101010001', 'sari.yuliani@klinikku.com', 'MIGRATION'
);

INSERT INTO users (
    username, password, role, dokter_id, created_at
) VALUES (
    'dokter_umum1',
    '$2a$10$Nhl3gJtZFxCdGv3dGnN6qOoyBi0sSbAg3ept81qPfcKvcAQXI3YWa',
    'dokter',
    (SELECT dokter_id FROM dokter WHERE no_izin_praktek = 'SIP-0001'),
    CURRENT_TIMESTAMP
);

-- Dokter 2: Dokter THT dengan mapping ke Sp.THT-KL
INSERT INTO dokter (
    no_izin_praktek, nama, tanggal_lahir, jenis_kelamin, alamat, no_hp, no_telepon, ktp, email, created_by
) VALUES (
    'SIP-0002', 'dr. Rudi Kurniawan', '1985-06-15', 'L',
    'Jl. Telinga No.2', '081234567892', '0217654321',
    '3175010202020002', 'rudi.kurniawan@klinikku.com', 'MIGRATION'
);

INSERT INTO users (
    username, password, role, dokter_id, created_at
) VALUES (
    'dokter_tht1',
    '$2a$10$Nhl3gJtZFxCdGv3dGnN6qOoyBi0sSbAg3ept81qPfcKvcAQXI3YWa',
    'dokter',
    (SELECT dokter_id FROM dokter WHERE no_izin_praktek = 'SIP-0002'),
    CURRENT_TIMESTAMP
);

-- Mapping ke spesialisasi Sp.THT-KL
INSERT INTO dokter_spesialisasi (
    dokter_id, spesialisasi_id, created_by
) VALUES (
    (SELECT dokter_id FROM dokter WHERE no_izin_praktek = 'SIP-0002'),
    (SELECT spesialisasi_id FROM spesialisasi WHERE kode_spesialisasi = 'Sp.THT-KL'),
    'MIGRATION'
);

-- +migrate Down
DELETE FROM dokter_spesialisasi
WHERE dokter_id = (SELECT dokter_id FROM dokter WHERE no_izin_praktek = 'SIP-0002')
  AND spesialisasi_id = (SELECT spesialisasi_id FROM spesialisasi WHERE kode_spesialisasi = 'Sp.THT-KL');

DELETE FROM users WHERE username IN ('dokter_umum', 'dokter_tht');

DELETE FROM dokter WHERE no_izin_praktek IN ('SIP-0001', 'SIP-0002');