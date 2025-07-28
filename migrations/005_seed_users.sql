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
    username, nama_lengkap, password, role, dokter_id, created_at
) VALUES (
    'dokter_umum1',
    'dr. Sari Yuliani',
    '$2a$14$9YPWgHyOj/s/ZPqv5cltOOPbXoKKPFnLC.wQEwiy42feLJAKnAPsq',
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
    username, nama_lengkap, password, role, dokter_id, created_at
) VALUES (
    'dokter_tht1',
    'dr. Rudi Kurniawan',
    '$2a$14$9YPWgHyOj/s/ZPqv5cltOOPbXoKKPFnLC.wQEwiy42feLJAKnAPsq',
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

-- Frontliner User
INSERT INTO users (
    username, nama_lengkap, password, role, created_at
) VALUES (
    'frontliner1',
    'Fitri Frontliner',
    '$2a$14$9YPWgHyOj/s/ZPqv5cltOOPbXoKKPFnLC.wQEwiy42feLJAKnAPsq',
    'frontliner',
    CURRENT_TIMESTAMP
);

-- Admin User
INSERT INTO users (
    username, nama_lengkap, password, role, created_at
) VALUES (
    'admin1',
    'Andi Admin',
    '$2a$14$9YPWgHyOj/s/ZPqv5cltOOPbXoKKPFnLC.wQEwiy42feLJAKnAPsq',
    'admin',
    CURRENT_TIMESTAMP
);

-- Farmasi User
INSERT INTO users (
    username, nama_lengkap, password, role, created_at
) VALUES (
    'farmasi1',
    'Fajar Farmasi',
    '$2a$14$9YPWgHyOj/s/ZPqv5cltOOPbXoKKPFnLC.wQEwiy42feLJAKnAPsq',
    'farmasi',
    CURRENT_TIMESTAMP
);

-- +migrate Down

DELETE FROM dokter_spesialisasi
WHERE dokter_id = (SELECT dokter_id FROM dokter WHERE no_izin_praktek = 'SIP-0002')
  AND spesialisasi_id = (SELECT spesialisasi_id FROM spesialisasi WHERE kode_spesialisasi = 'Sp.THT-KL');

DELETE FROM users 
WHERE username IN ('dokter_umum1', 'dokter_tht1', 'frontliner1', 'admin1', 'farmasi1');

DELETE FROM dokter 
WHERE no_izin_praktek IN ('SIP-0001', 'SIP-0002');
