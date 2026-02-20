package config

import (
	"os"
	"strconv"
	"time"
)

type Config struct {
	DatabaseURL      string
	JWTSecret        string
	TokenTTL         time.Duration
	HospitalABaseURL string
}

func Load() Config {
	ttlHours := 24
	if v := os.Getenv("TOKEN_TTL_HOURS"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 {
			ttlHours = n
		}
	}

	cfg := Config{
		DatabaseURL:      getenv("DATABASE_URL", "postgres://postgres:postgres@db:5432/agnos?sslmode=disable"),
		JWTSecret:        getenv("JWT_SECRET", "ky2>B(#0sB65D9Mj"),
		TokenTTL:         time.Duration(ttlHours) * time.Hour,
		HospitalABaseURL: getenv("HOSPITAL_A_BASE_URL", "https://hospital-a.api.co.th"),
	}
	return cfg
}

func getenv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
