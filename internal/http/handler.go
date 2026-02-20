package http

import (
	"database/sql"
	"errors"
	"net/http"

	"agnos/internal/middleware"
	"agnos/internal/model"
	"agnos/internal/service"

	"github.com/gin-gonic/gin"
)

type handler struct {
	staffService   service.StaffService
	patientService service.PatientService
}

func RegisterRoutes(r *gin.Engine, staffService service.StaffService, patientService service.PatientService, jwtSecret string) {
	h := &handler{staffService: staffService, patientService: patientService}

	r.GET("/healthz", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	r.POST("/staff/create", h.staffCreate)
	r.POST("/staff/login", h.staffLogin)

	authed := r.Group("", middleware.JWTAuth(jwtSecret))
	authed.POST("/patient/search", h.patientSearch)
}

type staffCreateRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Hospital string `json:"hospital"`
}

func (h *handler) staffCreate(c *gin.Context) {
	var req staffCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	staff, err := h.staffService.Create(req.Username, req.Password, req.Hospital)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, staff)
}

type staffLoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Hospital string `json:"hospital"`
}

func (h *handler) staffLogin(c *gin.Context) {
	var req staffLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}
	jwt, err := h.staffService.Login(req.Username, req.Password, req.Hospital)
	if err != nil {
		if errors.Is(err, service.ErrInvalidCredentials) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"token": jwt})
}

func (h *handler) patientSearch(c *gin.Context) {
	var criteria model.PatientSearchCriteria
	if err := c.ShouldBindJSON(&criteria); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}
	hospital := middleware.HospitalFromContext(c)
	result, err := h.patientService.Search(hospital, criteria)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "search failed"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"patients": result})
}
