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

type ProjectCreateRequest struct {
	AffiliationID   uint   `json:"affiliation_id" binding:"required,gt=0"`
	Name            string `json:"name" binding:"required"`
	Description     string `json:"description"`
	NumMax          uint   `json:"num_max" binding:"required,gt=0"`
	ProjectStartDay string `json:"project_start_day" binding:"required"` // RFC3339
	ProjectEndDay   string `json:"project_end_day" binding:"required"`   // RFC3339
	FormStartDay    string `json:"form_start_day" binding:"required"`    // RFC3339
	FormEndDay      string `json:"form_end_day" binding:"required"`      // RFC3339
}

type ProjectUpdateRequest struct {
	Name            *string `json:"name"`
	Description     *string `json:"description"`
	NumMax          *uint   `json:"num_max"`
	ProjectStartDay *string `json:"project_start_day"` // RFC3339
	ProjectEndDay   *string `json:"project_end_day"`   // RFC3339
	FormStartDay    *string `json:"form_start_day"`    // RFC3339
	FormEndDay      *string `json:"form_end_day"`      // RFC3339
}
