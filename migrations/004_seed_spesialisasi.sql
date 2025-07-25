-- +migrate Up
INSERT INTO spesialisasi (kode_spesialisasi, nama_spesialisasi, deskripsi, created_by)
VALUES 
('Sp.A', 'Spesialis Anak', 'Dokter spesialis yang menangani kesehatan anak dan bayi', 'seed'),
('Sp.PD', 'Spesialis Penyakit Dalam', 'Dokter spesialis untuk penyakit organ dalam dewasa', 'seed'),
('Sp.OG', 'Spesialis Obstetri dan Ginekologi', 'Dokter kandungan dan kebidanan', 'seed'),
('Sp.B', 'Spesialis Bedah', 'Dokter ahli dalam tindakan bedah umum', 'seed'),
('Sp.M', 'Spesialis Mata', 'Dokter spesialis yang menangani kesehatan mata', 'seed'),
('Sp.THT-KL', 'Spesialis Telinga Hidung Tenggorokan-Kepala Leher', 'Menangani masalah THT dan kepala-leher', 'seed'),
('Sp.JP', 'Spesialis Jantung dan Pembuluh Darah', 'Ahli jantung dan sistem kardiovaskular', 'seed'),
('Sp.KFR', 'Spesialis Kedokteran Fisik dan Rehabilitasi', 'Fisioterapi medis', 'seed'),
('Sp.Rad', 'Spesialis Radiologi', 'Ahli pencitraan medis (X-ray, CT, MRI)', 'seed'),
('Sp.An', 'Spesialis Anestesi', 'Ahli pembiusan untuk operasi dan perawatan intensif', 'seed'),
('Sp.KJ', 'Spesialis Kedokteran Jiwa', 'Dokter psikiater', 'seed'),
('Sp.PD-KGEH', 'Spesialis Penyakit Dalam Konsultan Gastroentero-Hepatologi', 'Ahli saluran pencernaan dan hati', 'seed');

-- +migrate Down
DELETE FROM spesialisasi 
WHERE kode_spesialisasi IN (
  'Sp.A', 'Sp.PD', 'Sp.OG', 'Sp.B', 'Sp.M', 
  'Sp.THT-KL', 'Sp.JP', 'Sp.KFR', 'Sp.Rad', 
  'Sp.An', 'Sp.KJ', 'Sp.PD-KGEH'
);