package repository

import "agnos/internal/model"

type StaffRepository interface {
	Create(username, passwordHash, hospital string) (model.Staff, error)
	FindByUsernameAndHospital(username, hospital string) (model.Staff, error)
}

type PatientRepository interface {
	SearchByHospital(hospital string, c model.PatientSearchCriteria) ([]model.Patient, error)
	FindByIdentifier(hospital string, nationalID, passportID *string) (model.Patient, bool, error)
	UpsertByNationalOrPassport(hospital string, p model.Patient) (model.Patient, error)
}
