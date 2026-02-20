CREATE TABLE IF NOT EXISTS staffs (
    id BIGSERIAL PRIMARY KEY,
    username VARCHAR(100) NOT NULL,
    password_hash TEXT NOT NULL,
    hospital VARCHAR(100) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE (username, hospital)
);

CREATE TABLE IF NOT EXISTS patients (
    id BIGSERIAL PRIMARY KEY,
    hospital VARCHAR(100) NOT NULL,
    first_name_th VARCHAR(255),
    middle_name_th VARCHAR(255),
    last_name_th VARCHAR(255),
    first_name_en VARCHAR(255),
    middle_name_en VARCHAR(255),
    last_name_en VARCHAR(255),
    date_of_birth DATE,
    patient_hn VARCHAR(100),
    national_id VARCHAR(50),
    passport_id VARCHAR(50),
    phone_number VARCHAR(50),
    email VARCHAR(255),
    gender CHAR(1),
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    CONSTRAINT chk_gender CHECK (gender IN ('M', 'F') OR gender IS NULL)
);

CREATE UNIQUE INDEX IF NOT EXISTS ux_patients_hospital_national_id
    ON patients (hospital, national_id)
    WHERE national_id IS NOT NULL;

CREATE UNIQUE INDEX IF NOT EXISTS ux_patients_hospital_passport_id
    ON patients (hospital, passport_id)
    WHERE passport_id IS NOT NULL;

CREATE INDEX IF NOT EXISTS idx_patients_hospital ON patients (hospital);
CREATE INDEX IF NOT EXISTS idx_patients_name_en ON patients (first_name_en, middle_name_en, last_name_en);
CREATE INDEX IF NOT EXISTS idx_patients_name_th ON patients (first_name_th, middle_name_th, last_name_th);
