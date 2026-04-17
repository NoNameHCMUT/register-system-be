package model

type RegisterRequest struct {
	Username      string `json:"username" binding:"required"`
	FullName      string `json:"full_name" binding:"required"`
	Password      string `json:"password" binding:"required,min=6"`
	Email         string `json:"email" binding:"required,email"`
	StudentID     string `json:"student_id"`
	AffiliationID uint   `json:"affiliation_id" binding:"required,gt=0"`
	Role          Role   `json:"role" binding:"required,oneof=admin student school community"`
	Phone         string `json:"phone"`
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
	ProjectStartDay string `json:"project_start_day" binding:"required"`
	ProjectEndDay   string `json:"project_end_day" binding:"required"`
	FormStartDay    string `json:"form_start_day" binding:"required"`
	FormEndDay      string `json:"form_end_day" binding:"required"`
}

type ProjectUpdateRequest struct {
	Name            *string `json:"name"`
	Description     *string `json:"description"`
	NumMax          *uint   `json:"num_max"`
	ProjectStartDay *string `json:"project_start_day"`
	ProjectEndDay   *string `json:"project_end_day"`
	FormStartDay    *string `json:"form_start_day"`
	FormEndDay      *string `json:"form_end_day"`
}

type UserUpdateRequest struct {
	FullName *string `json:"full_name"`
	Phone    *string `json:"phone"`
}

type ApplicationActionRequest struct {
	ApplicationIDs []uint `json:"application_ids" binding:"required,min=1"`
	Action         string `json:"action" binding:"required,oneof=approve reject"`
}

type AffiliationCreateRequest struct {
	StdName     string `json:"std_name" binding:"required"`
	Description string `json:"description"`
}

type AffiliationUpdateRequest struct {
	StdName     *string `json:"std_name"`
	Description *string `json:"description"`
}
