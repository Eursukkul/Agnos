package his

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"agnos/internal/model"
)

type HospitalAClient interface {
	FetchByID(id string) (model.Patient, error)
}

type hospitalAClient struct {
	baseURL string
	client  *http.Client
}

type hospitalAResponse struct {
	FirstNameTH  *string `json:"first_name_th"`
	MiddleNameTH *string `json:"middle_name_th"`
	LastNameTH   *string `json:"last_name_th"`
	FirstNameEN  *string `json:"first_name_en"`
	MiddleNameEN *string `json:"middle_name_en"`
	LastNameEN   *string `json:"last_name_en"`
	DateOfBirth  *string `json:"date_of_birth"`
	PatientHN    *string `json:"patient_hn"`
	NationalID   *string `json:"national_id"`
	PassportID   *string `json:"passport_id"`
	PhoneNumber  *string `json:"phone_number"`
	Email        *string `json:"email"`
	Gender       *string `json:"gender"`
}

func NewHospitalAClient(baseURL string, client *http.Client) HospitalAClient {
	if client == nil {
		client = http.DefaultClient
	}
	return &hospitalAClient{baseURL: strings.TrimRight(baseURL, "/"), client: client}
}

func (c *hospitalAClient) FetchByID(id string) (model.Patient, error) {
	id = strings.TrimSpace(id)
	if id == "" {
		return model.Patient{}, fmt.Errorf("id is required")
	}

	u := c.baseURL + "/patient/search/" + url.PathEscape(id)
	req, err := http.NewRequest(http.MethodGet, u, nil)
	if err != nil {
		return model.Patient{}, err
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return model.Patient{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return model.Patient{}, fmt.Errorf("hospital-a responded with %d", resp.StatusCode)
	}

	var payload hospitalAResponse
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return model.Patient{}, err
	}

	var dob *time.Time
	if payload.DateOfBirth != nil && strings.TrimSpace(*payload.DateOfBirth) != "" {
		if t, err := time.Parse("2006-01-02", strings.TrimSpace(*payload.DateOfBirth)); err == nil {
			dob = &t
		}
	}

	return model.Patient{
		FirstNameTH:  payload.FirstNameTH,
		MiddleNameTH: payload.MiddleNameTH,
		LastNameTH:   payload.LastNameTH,
		FirstNameEN:  payload.FirstNameEN,
		MiddleNameEN: payload.MiddleNameEN,
		LastNameEN:   payload.LastNameEN,
		DateOfBirth:  dob,
		PatientHN:    payload.PatientHN,
		NationalID:   payload.NationalID,
		PassportID:   payload.PassportID,
		PhoneNumber:  payload.PhoneNumber,
		Email:        payload.Email,
		Gender:       payload.Gender,
	}, nil
}
