package model

import "time"

type Staff struct {
	ID           int64  `json:"id"`
	Username     string `json:"username"`
	PasswordHash string `json:"-"`
	Hospital     string `json:"hospital"`
}

type Patient struct {
	ID           int64      `json:"id"`
	Hospital     string     `json:"hospital"`
	FirstNameTH  *string    `json:"first_name_th,omitempty"`
	MiddleNameTH *string    `json:"middle_name_th,omitempty"`
	LastNameTH   *string    `json:"last_name_th,omitempty"`
	FirstNameEN  *string    `json:"first_name_en,omitempty"`
	MiddleNameEN *string    `json:"middle_name_en,omitempty"`
	LastNameEN   *string    `json:"last_name_en,omitempty"`
	DateOfBirth  *time.Time `json:"date_of_birth,omitempty"`
	PatientHN    *string    `json:"patient_hn,omitempty"`
	NationalID   *string    `json:"national_id,omitempty"`
	PassportID   *string    `json:"passport_id,omitempty"`
	PhoneNumber  *string    `json:"phone_number,omitempty"`
	Email        *string    `json:"email,omitempty"`
	Gender       *string    `json:"gender,omitempty"`
}

type PatientSearchCriteria struct {
	NationalID  *string `json:"national_id"`
	PassportID  *string `json:"passport_id"`
	FirstName   *string `json:"first_name"`
	MiddleName  *string `json:"middle_name"`
	LastName    *string `json:"last_name"`
	DateOfBirth *string `json:"date_of_birth"`
	PhoneNumber *string `json:"phone_number"`
	Email       *string `json:"email"`
}
