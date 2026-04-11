package model

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