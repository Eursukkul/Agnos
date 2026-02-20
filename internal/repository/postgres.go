package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"agnos/internal/model"
)

type postgresStaffRepository struct {
	db *sql.DB
}

type postgresPatientRepository struct {
	db *sql.DB
}

func NewPostgresStaffRepository(db *sql.DB) StaffRepository {
	return &postgresStaffRepository{db: db}
}

func NewPostgresPatientRepository(db *sql.DB) PatientRepository {
	return &postgresPatientRepository{db: db}
}

func (r *postgresStaffRepository) Create(username, passwordHash, hospital string) (model.Staff, error) {
	var s model.Staff
	err := r.db.QueryRow(
		`INSERT INTO staffs (username, password_hash, hospital) VALUES ($1, $2, $3)
		 RETURNING id, username, password_hash, hospital`,
		username, passwordHash, hospital,
	).Scan(&s.ID, &s.Username, &s.PasswordHash, &s.Hospital)
	return s, err
}

func (r *postgresStaffRepository) FindByUsernameAndHospital(username, hospital string) (model.Staff, error) {
	var s model.Staff
	err := r.db.QueryRow(
		`SELECT id, username, password_hash, hospital FROM staffs WHERE username = $1 AND hospital = $2`,
		username, hospital,
	).Scan(&s.ID, &s.Username, &s.PasswordHash, &s.Hospital)
	return s, err
}

func (r *postgresPatientRepository) SearchByHospital(hospital string, c model.PatientSearchCriteria) ([]model.Patient, error) {
	base := `SELECT id, hospital, first_name_th, middle_name_th, last_name_th, first_name_en, middle_name_en,
		last_name_en, date_of_birth, patient_hn, national_id, passport_id, phone_number, email, gender
		FROM patients WHERE hospital = $1`
	args := []any{hospital}
	idx := 2

	appendLike := func(cond string, value *string) {
		if value == nil || strings.TrimSpace(*value) == "" {
			return
		}
		base += fmt.Sprintf(" AND %s", fmt.Sprintf(cond, idx))
		args = append(args, "%"+strings.TrimSpace(*value)+"%")
		idx++
	}
	appendEq := func(field string, value *string) {
		if value == nil || strings.TrimSpace(*value) == "" {
			return
		}
		base += fmt.Sprintf(" AND %s = $%d", field, idx)
		args = append(args, strings.TrimSpace(*value))
		idx++
	}

	appendEq("national_id", c.NationalID)
	appendEq("passport_id", c.PassportID)
	appendLike("(first_name_en ILIKE $%d OR first_name_th ILIKE $%d)", c.FirstName)
	appendLike("(middle_name_en ILIKE $%d OR middle_name_th ILIKE $%d)", c.MiddleName)
	appendLike("(last_name_en ILIKE $%d OR last_name_th ILIKE $%d)", c.LastName)
	appendLike("phone_number ILIKE $%d", c.PhoneNumber)
	appendLike("email ILIKE $%d", c.Email)

	if c.DateOfBirth != nil && strings.TrimSpace(*c.DateOfBirth) != "" {
		base += fmt.Sprintf(" AND date_of_birth = $%d", idx)
		args = append(args, strings.TrimSpace(*c.DateOfBirth))
	}

	base += ` ORDER BY id DESC LIMIT 100`

	rows, err := r.db.Query(base, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make([]model.Patient, 0)
	for rows.Next() {
		p, err := scanPatient(rows)
		if err != nil {
			return nil, err
		}
		result = append(result, p)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return result, nil
}

func (r *postgresPatientRepository) FindByIdentifier(hospital string, nationalID, passportID *string) (model.Patient, bool, error) {
	if (nationalID == nil || strings.TrimSpace(*nationalID) == "") && (passportID == nil || strings.TrimSpace(*passportID) == "") {
		return model.Patient{}, false, nil
	}
	query := `SELECT id, hospital, first_name_th, middle_name_th, last_name_th, first_name_en, middle_name_en,
		last_name_en, date_of_birth, patient_hn, national_id, passport_id, phone_number, email, gender
		FROM patients WHERE hospital = $1 AND (`
	args := []any{hospital}
	idx := 2
	conds := make([]string, 0, 2)
	if nationalID != nil && strings.TrimSpace(*nationalID) != "" {
		conds = append(conds, fmt.Sprintf("national_id = $%d", idx))
		args = append(args, strings.TrimSpace(*nationalID))
		idx++
	}
	if passportID != nil && strings.TrimSpace(*passportID) != "" {
		conds = append(conds, fmt.Sprintf("passport_id = $%d", idx))
		args = append(args, strings.TrimSpace(*passportID))
	}
	query += strings.Join(conds, " OR ") + `) LIMIT 1`

	row := r.db.QueryRow(query, args...)
	p, err := scanPatient(row)
	if errors.Is(err, sql.ErrNoRows) {
		return model.Patient{}, false, nil
	}
	if err != nil {
		return model.Patient{}, false, err
	}
	return p, true, nil
}

func (r *postgresPatientRepository) UpsertByNationalOrPassport(hospital string, p model.Patient) (model.Patient, error) {
	var dob any
	if p.DateOfBirth != nil {
		dob = p.DateOfBirth.Format("2006-01-02")
	}

	row := r.db.QueryRow(`
		INSERT INTO patients (
			hospital, first_name_th, middle_name_th, last_name_th,
			first_name_en, middle_name_en, last_name_en, date_of_birth,
			patient_hn, national_id, passport_id, phone_number, email, gender
		)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14)
		ON CONFLICT (hospital, national_id) WHERE national_id IS NOT NULL
		DO UPDATE SET
			first_name_th = EXCLUDED.first_name_th,
			middle_name_th = EXCLUDED.middle_name_th,
			last_name_th = EXCLUDED.last_name_th,
			first_name_en = EXCLUDED.first_name_en,
			middle_name_en = EXCLUDED.middle_name_en,
			last_name_en = EXCLUDED.last_name_en,
			date_of_birth = EXCLUDED.date_of_birth,
			patient_hn = EXCLUDED.patient_hn,
			passport_id = EXCLUDED.passport_id,
			phone_number = EXCLUDED.phone_number,
			email = EXCLUDED.email,
			gender = EXCLUDED.gender,
			updated_at = now()
		RETURNING id, hospital, first_name_th, middle_name_th, last_name_th, first_name_en, middle_name_en,
			last_name_en, date_of_birth, patient_hn, national_id, passport_id, phone_number, email, gender
	`, hospital, p.FirstNameTH, p.MiddleNameTH, p.LastNameTH, p.FirstNameEN, p.MiddleNameEN, p.LastNameEN,
		dob, p.PatientHN, p.NationalID, p.PassportID, p.PhoneNumber, p.Email, p.Gender)

	stored, err := scanPatient(row)
	if err == nil {
		return stored, nil
	}

	row = r.db.QueryRow(`
		INSERT INTO patients (
			hospital, first_name_th, middle_name_th, last_name_th,
			first_name_en, middle_name_en, last_name_en, date_of_birth,
			patient_hn, national_id, passport_id, phone_number, email, gender
		)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14)
		ON CONFLICT (hospital, passport_id) WHERE passport_id IS NOT NULL
		DO UPDATE SET
			first_name_th = EXCLUDED.first_name_th,
			middle_name_th = EXCLUDED.middle_name_th,
			last_name_th = EXCLUDED.last_name_th,
			first_name_en = EXCLUDED.first_name_en,
			middle_name_en = EXCLUDED.middle_name_en,
			last_name_en = EXCLUDED.last_name_en,
			date_of_birth = EXCLUDED.date_of_birth,
			patient_hn = EXCLUDED.patient_hn,
			national_id = EXCLUDED.national_id,
			phone_number = EXCLUDED.phone_number,
			email = EXCLUDED.email,
			gender = EXCLUDED.gender,
			updated_at = now()
		RETURNING id, hospital, first_name_th, middle_name_th, last_name_th, first_name_en, middle_name_en,
			last_name_en, date_of_birth, patient_hn, national_id, passport_id, phone_number, email, gender
	`, hospital, p.FirstNameTH, p.MiddleNameTH, p.LastNameTH, p.FirstNameEN, p.MiddleNameEN, p.LastNameEN,
		dob, p.PatientHN, p.NationalID, p.PassportID, p.PhoneNumber, p.Email, p.Gender)

	return scanPatient(row)
}

type rowScanner interface {
	Scan(dest ...any) error
}

func scanPatient(s rowScanner) (model.Patient, error) {
	var p model.Patient
	var dob sql.NullTime
	var firstTH, middleTH, lastTH, firstEN, middleEN, lastEN sql.NullString
	var hn, nationalID, passportID, phone, email, gender sql.NullString

	err := s.Scan(
		&p.ID,
		&p.Hospital,
		&firstTH,
		&middleTH,
		&lastTH,
		&firstEN,
		&middleEN,
		&lastEN,
		&dob,
		&hn,
		&nationalID,
		&passportID,
		&phone,
		&email,
		&gender,
	)
	if err != nil {
		return model.Patient{}, err
	}
	if firstTH.Valid {
		p.FirstNameTH = &firstTH.String
	}
	if middleTH.Valid {
		p.MiddleNameTH = &middleTH.String
	}
	if lastTH.Valid {
		p.LastNameTH = &lastTH.String
	}
	if firstEN.Valid {
		p.FirstNameEN = &firstEN.String
	}
	if middleEN.Valid {
		p.MiddleNameEN = &middleEN.String
	}
	if lastEN.Valid {
		p.LastNameEN = &lastEN.String
	}
	if dob.Valid {
		t := dob.Time
		t = time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.UTC)
		p.DateOfBirth = &t
	}
	if hn.Valid {
		p.PatientHN = &hn.String
	}
	if nationalID.Valid {
		p.NationalID = &nationalID.String
	}
	if passportID.Valid {
		p.PassportID = &passportID.String
	}
	if phone.Valid {
		p.PhoneNumber = &phone.String
	}
	if email.Valid {
		p.Email = &email.String
	}
	if gender.Valid {
		p.Gender = &gender.String
	}
	return p, nil
}
