# ER Diagram

```mermaid
erDiagram
    STAFFS {
        BIGSERIAL id PK
        VARCHAR username
        TEXT password_hash
        VARCHAR hospital
        TIMESTAMPTZ created_at
    }

    PATIENTS {
        BIGSERIAL id PK
        VARCHAR hospital
        VARCHAR first_name_th
        VARCHAR middle_name_th
        VARCHAR last_name_th
        VARCHAR first_name_en
        VARCHAR middle_name_en
        VARCHAR last_name_en
        DATE date_of_birth
        VARCHAR patient_hn
        VARCHAR national_id
        VARCHAR passport_id
        VARCHAR phone_number
        VARCHAR email
        CHAR gender
        TIMESTAMPTZ created_at
        TIMESTAMPTZ updated_at
    }
```

Notes:
- `staffs` unique key: `(username, hospital)`.
- `patients` unique partial indexes: `(hospital, national_id)` and `(hospital, passport_id)`.
- Access control is enforced by JWT claim `hospital` for patient search.
