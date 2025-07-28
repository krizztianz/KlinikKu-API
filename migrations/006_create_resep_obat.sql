-- +migrate Up

CREATE TABLE obat (
    obat_id SERIAL PRIMARY KEY,
    nama_obat VARCHAR(100) NOT NULL,
    deskripsi TEXT,
    stok INT DEFAULT 0,
    harga NUMERIC(19,2) DEFAULT 0,
    satuan VARCHAR(50),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_by VARCHAR(255),
    modified_at TIMESTAMP,
    modified_by VARCHAR(255)
);

CREATE TABLE resep (
    resep_id SERIAL PRIMARY KEY,
    rekam_medis_id INT UNIQUE REFERENCES rekam_medis(rekam_medis_id) ON DELETE CASCADE,
    status resep_status NOT NULL DEFAULT 'menunggu',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_by VARCHAR(255),
    modified_at TIMESTAMP,
    modified_by VARCHAR(255)
);

CREATE TABLE resep_detail (
    resep_detail_id SERIAL PRIMARY KEY,
    resep_id INT REFERENCES resep(resep_id) ON DELETE CASCADE,
    obat_id INT REFERENCES obat(obat_id) ON DELETE RESTRICT,
    jumlah INT NOT NULL,
    dosis VARCHAR(100),
    keterangan TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_by VARCHAR(255),
    modified_at TIMESTAMP,
    modified_by VARCHAR(255)
);

-- +migrate Down

DROP TABLE IF EXISTS resep_detail;
DROP TABLE IF EXISTS resep;
DROP TABLE IF EXISTS obat;