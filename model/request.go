package model

import "time"

type RegisterRequest struct {
	Username      string `json:"username" binding:"required"`
	FullName      string `json:"full_name" binding:"required"`
	Password      string `json:"password" binding:"required,min=6"`
	Email         string `json:"email" binding:"required,email"`
	StudentID     string `json:"student_id"`
	AffiliationID uint   `json:"affiliation_id" binding:"required,gt=0"`
	Role          Role   `json:"role" binding:"required,oneof=admin student school community"`
}

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type RefreshRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type CreateProjectRequest struct {
	AffiliationID   uint      `json:"affiliation_id" binding:"required"`
	Name            string    `json:"name" binding:"required"`
	Description     string    `json:"description"`
	NumMax          uint      `json:"num_max"`
	ProjectStartDay time.Time `json:"project_start_day" binding:"required"`
	ProjectEndDay   time.Time `json:"project_end_day" binding:"required"`
	FormStartDay    time.Time `json:"form_start_day" binding:"required"`
	FormEndDay      time.Time `json:"form_end_day" binding:"required"`
}
