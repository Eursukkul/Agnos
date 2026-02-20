package service

import (
	"strings"

	"agnos/internal/his"
	"agnos/internal/model"
	"agnos/internal/repository"
)

type PatientService interface {
	Search(hospital string, c model.PatientSearchCriteria) ([]model.Patient, error)
}

type patientService struct {
	repo      repository.PatientRepository
	hisClient his.HospitalAClient
}

func NewPatientService(repo repository.PatientRepository, hisClient his.HospitalAClient) PatientService {
	return &patientService{repo: repo, hisClient: hisClient}
}

func (s *patientService) Search(hospital string, c model.PatientSearchCriteria) ([]model.Patient, error) {
	hospital = strings.TrimSpace(hospital)
	if hospital == "" {
		return nil, nil
	}

	if (c.NationalID != nil && strings.TrimSpace(*c.NationalID) != "") || (c.PassportID != nil && strings.TrimSpace(*c.PassportID) != "") {
		_, found, err := s.repo.FindByIdentifier(hospital, c.NationalID, c.PassportID)
		if err != nil {
			return nil, err
		}
		if !found {
			id := ""
			if c.NationalID != nil && strings.TrimSpace(*c.NationalID) != "" {
				id = strings.TrimSpace(*c.NationalID)
			} else if c.PassportID != nil {
				id = strings.TrimSpace(*c.PassportID)
			}
			if id != "" {
				if externalPatient, err := s.hisClient.FetchByID(id); err == nil {
					externalPatient.Hospital = hospital
					_, _ = s.repo.UpsertByNationalOrPassport(hospital, externalPatient)
				}
			}
		}
	}

	return s.repo.SearchByHospital(hospital, c)
}
