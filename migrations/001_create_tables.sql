-- +migrate Up

-- ENUM types
DO $$ BEGIN IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'gender') THEN CREATE TYPE gender AS ENUM ('L', 'P'); END IF; END $$;
DO $$ BEGIN IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'blood_type') THEN CREATE TYPE blood_type AS ENUM ('A', 'B', 'AB', 'O'); END IF; END $$;
DO $$ BEGIN IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'visit_type') THEN CREATE TYPE visit_type AS ENUM ('Umum', 'Rawat Jalan', 'Kontrol'); END IF; END $$;
DO $$ BEGIN IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'visit_status') THEN CREATE TYPE visit_status AS ENUM ('registrasi', 'pemeriksaan', 'menunggu_resep', 'selesai', 'batal'); END IF; END $$;
DO $$ BEGIN IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'visit_priority') THEN CREATE TYPE visit_priority AS ENUM ('normal', 'prioritas'); END IF; END $$;

-- MASTER TABLES
CREATE TABLE IF NOT EXISTS pasien (
    pasien_id SERIAL PRIMARY KEY,
    nama VARCHAR(255) NOT NULL,
    tanggal_lahir DATE NOT NULL,
    jenis_kelamin gender NOT NULL,
    alamat TEXT NOT NULL,
    no_hp VARCHAR(20) NOT NULL,
    no_telepon VARCHAR(20),
    ktp VARCHAR(20) NOT NULL UNIQUE,
    email VARCHAR(100),
    golongan_darah blood_type,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_by VARCHAR(255),
    modified_at TIMESTAMP,
    modified_by VARCHAR(255)
);

CREATE TABLE IF NOT EXISTS dokter (
    dokter_id SERIAL PRIMARY KEY,
    no_izin_praktek VARCHAR(65) NULL UNIQUE,
    nama VARCHAR NOT NULL,
    tanggal_lahir DATE NOT NULL,
    jenis_kelamin gender,
    alamat TEXT NOT NULL,
    no_hp VARCHAR(20) NOT NULL,
    no_telepon VARCHAR(20),
    ktp VARCHAR(20) NOT NULL UNIQUE,
    email VARCHAR(255) UNIQUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_by VARCHAR,
    modified_at TIMESTAMP,
    modified_by VARCHAR
);

CREATE TABLE IF NOT EXISTS spesialisasi (
    spesialisasi_id SERIAL PRIMARY KEY,
    kode_spesialisasi VARCHAR(25) NOT NULL,
    nama_spesialisasi VARCHAR(100) NOT NULL,
    deskripsi TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_by VARCHAR(255),
    modified_at TIMESTAMP,
    modified_by VARCHAR(255)
);

CREATE TABLE IF NOT EXISTS dokter_spesialisasi (
    dokter_id INT REFERENCES dokter(dokter_id) ON DELETE RESTRICT ON UPDATE CASCADE,
    spesialisasi_id INT REFERENCES spesialisasi(spesialisasi_id) ON DELETE RESTRICT ON UPDATE CASCADE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_by VARCHAR(255),
    modified_at TIMESTAMP,
    modified_by VARCHAR(255),
    PRIMARY KEY(dokter_id, spesialisasi_id)
);

CREATE TABLE IF NOT EXISTS diagnosa (
    diagnosa_id SERIAL PRIMARY KEY,
    kode_icd10 VARCHAR(25) NOT NULL,
    nama_diagnosa VARCHAR(255) NOT NULL,
    deskripsi TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_by VARCHAR(255),
    modified_at TIMESTAMP,
    modified_by VARCHAR(255)
);

CREATE TABLE IF NOT EXISTS tindakan (
    tindakan_id SERIAL PRIMARY KEY,
    kode_icd9 VARCHAR(25) NOT NULL,
    nama_tindakan VARCHAR(255) NOT NULL,
    deskripsi TEXT,
    biaya_dasar NUMERIC(19,2),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_by VARCHAR(255),
    modified_at TIMESTAMP,
    modified_by VARCHAR(255)
);

-- TRANSACTION TABLES
CREATE TABLE IF NOT EXISTS kunjungan (
    kunjungan_id SERIAL PRIMARY KEY,
    pasien_id INT NOT NULL REFERENCES pasien(pasien_id) ON DELETE RESTRICT ON UPDATE CASCADE,
    dokter_id INT NOT NULL REFERENCES dokter(dokter_id) ON DELETE RESTRICT ON UPDATE CASCADE,
    tanggal_kunjungan TIMESTAMP,
    keluhan TEXT NOT NULL,
    tinggi_badan NUMERIC(3,1),
    berat_badan NUMERIC(3,2),
    tekanan_darah VARCHAR(50),
    suhu_tubuh NUMERIC(3,1),
    jenis_kunjungan visit_type,
    status visit_status NOT NULL,
    prioritas visit_priority DEFAULT 'normal'
);

CREATE TABLE IF NOT EXISTS rekam_medis (
    rekam_medis_id SERIAL PRIMARY KEY,
    kunjungan_id INT UNIQUE REFERENCES kunjungan(kunjungan_id) ON DELETE RESTRICT ON UPDATE CASCADE,
    dokter_id INT REFERENCES dokter(dokter_id) ON DELETE RESTRICT ON UPDATE CASCADE,
    anamnesa TEXT NOT NULL,
    pemeriksaan_fisik TEXT NOT NULL,
    catatan TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_by VARCHAR(255),
    modified_at TIMESTAMP,
    modified_by VARCHAR(255)
);

CREATE TABLE IF NOT EXISTS rekam_medis_diagnosa (
    rekam_medis_id INT NOT NULL REFERENCES rekam_medis(rekam_medis_id) ON DELETE RESTRICT ON UPDATE CASCADE,
    diagnosa_id INT NOT NULL REFERENCES diagnosa(diagnosa_id) ON DELETE RESTRICT ON UPDATE CASCADE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_by VARCHAR(255),
    modified_at TIMESTAMP,
    modified_by VARCHAR(255)
);

CREATE TABLE IF NOT EXISTS rekam_medis_tindakan (
    rekam_medis_id INT NOT NULL REFERENCES rekam_medis(rekam_medis_id) ON DELETE RESTRICT ON UPDATE CASCADE,
    tindakan_id INT NOT NULL REFERENCES tindakan(tindakan_id) ON DELETE RESTRICT ON UPDATE CASCADE,
    jumlah INT DEFAULT 1 NOT NULL,
    biaya_aktual NUMERIC(19,2) NOT NULL,
    catatan TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_by VARCHAR(255),
    modified_at TIMESTAMP,
    modified_by VARCHAR(255)
);

-- USERS TABLE (LOGIN)
CREATE TABLE IF NOT EXISTS users (
    user_id SERIAL PRIMARY KEY,
    username VARCHAR(50) NOT NULL UNIQUE,
    password TEXT NOT NULL,
    role VARCHAR(20) NOT NULL CHECK (role IN ('admin', 'frontliner', 'dokter', 'farmasi')),
    dokter_id INT REFERENCES dokter(dokter_id) ON DELETE SET NULL,
    refresh_token TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- +migrate Down

DROP TABLE IF EXISTS users;
DROP TABLE IF EXISTS rekam_medis_tindakan;
DROP TABLE IF EXISTS rekam_medis_diagnosa;
DROP TABLE IF EXISTS rekam_medis;
DROP TABLE IF EXISTS kunjungan;
DROP TABLE IF EXISTS tindakan;
DROP TABLE IF EXISTS diagnosa;
DROP TABLE IF EXISTS dokter_spesialisasi;
DROP TABLE IF EXISTS spesialisasi;
DROP TABLE IF EXISTS dokter;
DROP TABLE IF EXISTS pasien;

DROP TYPE IF EXISTS gender;
DROP TYPE IF EXISTS blood_type;
DROP TYPE IF EXISTS visit_type;
DROP TYPE IF EXISTS visit_status;
DROP TYPE IF EXISTS visit_priority;