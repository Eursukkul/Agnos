package http

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"agnos/internal/model"
	"agnos/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type fakeStaffService struct {
	createFn func(username, password, hospital string) (model.Staff, error)
	loginFn  func(username, password, hospital string) (string, error)
}

func (f *fakeStaffService) Create(username, password, hospital string) (model.Staff, error) {
	return f.createFn(username, password, hospital)
}

func (f *fakeStaffService) Login(username, password, hospital string) (string, error) {
	return f.loginFn(username, password, hospital)
}

type fakePatientService struct {
	searchFn func(hospital string, c model.PatientSearchCriteria) ([]model.Patient, error)
}

func (f *fakePatientService) Search(hospital string, c model.PatientSearchCriteria) ([]model.Patient, error) {
	return f.searchFn(hospital, c)
}

func setupRouter(staff service.StaffService, patient service.PatientService) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	RegisterRoutes(r, staff, patient, "secret")
	return r
}

func TestStaffCreateSuccess(t *testing.T) {
	r := setupRouter(&fakeStaffService{
		createFn: func(username, password, hospital string) (model.Staff, error) {
			return model.Staff{ID: 1, Username: username, Hospital: hospital}, nil
		},
		loginFn: func(username, password, hospital string) (string, error) {
			return "", nil
		},
	}, &fakePatientService{searchFn: func(hospital string, c model.PatientSearchCriteria) ([]model.Patient, error) { return nil, nil }})

	body := map[string]string{"username": "john", "password": "secret", "hospital": "A"}
	b, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/staff/create", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201 got %d", w.Code)
	}
}

func TestStaffCreateBadRequest(t *testing.T) {
	r := setupRouter(&fakeStaffService{
		createFn: func(username, password, hospital string) (model.Staff, error) {
			return model.Staff{}, errors.New("invalid")
		},
		loginFn: func(username, password, hospital string) (string, error) {
			return "", nil
		},
	}, &fakePatientService{searchFn: func(hospital string, c model.PatientSearchCriteria) ([]model.Patient, error) { return nil, nil }})

	req := httptest.NewRequest(http.MethodPost, "/staff/create", bytes.NewReader([]byte(`{}`)))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 got %d", w.Code)
	}
}

func TestStaffLoginSuccess(t *testing.T) {
	r := setupRouter(&fakeStaffService{
		createFn: func(username, password, hospital string) (model.Staff, error) {
			return model.Staff{}, nil
		},
		loginFn: func(username, password, hospital string) (string, error) {
			return "token", nil
		},
	}, &fakePatientService{searchFn: func(hospital string, c model.PatientSearchCriteria) ([]model.Patient, error) { return nil, nil }})

	req := httptest.NewRequest(http.MethodPost, "/staff/login", bytes.NewReader([]byte(`{"username":"john","password":"secret","hospital":"A"}`)))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200 got %d", w.Code)
	}
}

func TestStaffLoginUnauthorized(t *testing.T) {
	r := setupRouter(&fakeStaffService{
		createFn: func(username, password, hospital string) (model.Staff, error) {
			return model.Staff{}, nil
		},
		loginFn: func(username, password, hospital string) (string, error) {
			return "", service.ErrInvalidCredentials
		},
	}, &fakePatientService{searchFn: func(hospital string, c model.PatientSearchCriteria) ([]model.Patient, error) { return nil, nil }})

	req := httptest.NewRequest(http.MethodPost, "/staff/login", bytes.NewReader([]byte(`{"username":"john","password":"wrong","hospital":"A"}`)))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401 got %d", w.Code)
	}
}

func TestPatientSearchSuccess(t *testing.T) {
	patient := model.Patient{ID: 1, Hospital: "A"}
	dob, _ := time.Parse("2006-01-02", "1990-01-01")
	patient.DateOfBirth = &dob

	r := setupRouter(&fakeStaffService{
		createFn: func(username, password, hospital string) (model.Staff, error) { return model.Staff{}, nil },
		loginFn:  func(username, password, hospital string) (string, error) { return "", nil },
	}, &fakePatientService{searchFn: func(hospital string, c model.PatientSearchCriteria) ([]model.Patient, error) {
		if hospital != "A" {
			t.Fatalf("expected hospital A, got %s", hospital)
		}
		return []model.Patient{patient}, nil
	}})

	tokenObj := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"hospital": "A",
		"iat":      time.Now().Unix(),
		"exp":      time.Now().Add(1 * time.Hour).Unix(),
	})
	token, err := tokenObj.SignedString([]byte("secret"))
	if err != nil {
		t.Fatalf("sign token: %v", err)
	}
	req := httptest.NewRequest(http.MethodPost, "/patient/search", bytes.NewReader([]byte(`{"first_name":"Jo"}`)))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200 got %d, body=%s", w.Code, w.Body.String())
	}
}

func TestPatientSearchUnauthorized(t *testing.T) {
	r := setupRouter(&fakeStaffService{
		createFn: func(username, password, hospital string) (model.Staff, error) { return model.Staff{}, nil },
		loginFn:  func(username, password, hospital string) (string, error) { return "", nil },
	}, &fakePatientService{searchFn: func(hospital string, c model.PatientSearchCriteria) ([]model.Patient, error) {
		return nil, nil
	}})

	req := httptest.NewRequest(http.MethodPost, "/patient/search", bytes.NewReader([]byte(`{}`)))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401 got %d", w.Code)
	}
}

func TestPatientSearchBadRequestInvalidBody(t *testing.T) {
	r := setupRouter(&fakeStaffService{
		createFn: func(username, password, hospital string) (model.Staff, error) { return model.Staff{}, nil },
		loginFn:  func(username, password, hospital string) (string, error) { return "", nil },
	}, &fakePatientService{searchFn: func(hospital string, c model.PatientSearchCriteria) ([]model.Patient, error) {
		t.Fatalf("search service should not be called on invalid request body")
		return nil, nil
	}})

	tokenObj := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"hospital": "A",
		"iat":      time.Now().Unix(),
		"exp":      time.Now().Add(1 * time.Hour).Unix(),
	})
	token, err := tokenObj.SignedString([]byte("secret"))
	if err != nil {
		t.Fatalf("sign token: %v", err)
	}

	req := httptest.NewRequest(http.MethodPost, "/patient/search", bytes.NewReader([]byte(`{"date_of_birth":123}`)))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 got %d, body=%s", w.Code, w.Body.String())
	}
}

func TestPatientSearchInternalServerError(t *testing.T) {
	r := setupRouter(&fakeStaffService{
		createFn: func(username, password, hospital string) (model.Staff, error) { return model.Staff{}, nil },
		loginFn:  func(username, password, hospital string) (string, error) { return "", nil },
	}, &fakePatientService{searchFn: func(hospital string, c model.PatientSearchCriteria) ([]model.Patient, error) {
		return nil, errors.New("db error")
	}})

	tokenObj := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"hospital": "A",
		"iat":      time.Now().Unix(),
		"exp":      time.Now().Add(1 * time.Hour).Unix(),
	})
	token, err := tokenObj.SignedString([]byte("secret"))
	if err != nil {
		t.Fatalf("sign token: %v", err)
	}

	req := httptest.NewRequest(http.MethodPost, "/patient/search", bytes.NewReader([]byte(`{"first_name":"Jo"}`)))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500 got %d, body=%s", w.Code, w.Body.String())
	}
}
