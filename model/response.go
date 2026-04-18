package model

type AuthResponse struct {
	AccessToken  string       `json:"access_token"`
	RefreshToken string       `json:"refresh_token"`
	User         UserResponse `json:"user"`
}

type RegisterResponse struct {
	AccessToken  string       `json:"access_token"`
	RefreshToken string       `json:"refresh_token"`
	User         UserResponse `json:"user"`
	Message      string       `json:"message"`
}

type UserResponse struct {
	ID            uint                 `json:"id"`
	Username      string               `json:"username"`
	FullName      string               `json:"full_name"`
	Role          Role                 `json:"role"`
	IsActive      bool                 `json:"is_active"`
	StudentID     string               `json:"student_id"`
	Email         string               `json:"email"`
	AffiliationID uint                 `json:"affiliation_id"`
	Affiliation   *AffiliationResponse `json:"affiliation,omitempty"`
	AvatarURL     string               `json:"avatar_url"`
	Phone         string               `json:"phone"`
}

type AffiliationResponse struct {
	ID          uint   `json:"id"`
	StdName     string `json:"std_name"`
	Description string `json:"description,omitempty"`
}

type ProjectResponse struct {
	ID              uint                 `json:"id"`
	AffiliationID   uint                 `json:"affiliation_id"`
	Affiliation     *AffiliationResponse `json:"affiliation,omitempty"`
	CommunityUserID uint                 `json:"community_user_id"`
	CommunityUser   *UserResponse        `json:"community_user,omitempty"`
	Name            string               `json:"name"`
	Description     string               `json:"description"`
	NumMax          uint                 `json:"num_max"`
	NumAttending    uint                 `json:"num_attending"`
	ProjectStartDay string               `json:"project_start_day"`
	ProjectEndDay   string               `json:"project_end_day"`
	FormStartDay    string               `json:"form_start_day"`
	FormEndDay      string               `json:"form_end_day"`
	DateApproved    string               `json:"date_approved,omitempty"`
	BannerURL       string               `json:"banner_url,omitempty"`
	CreatedAt       string               `json:"created_at"`
}

type StudentProjectResponse struct {
	ID        uint              `json:"id"`
	UserID    uint              `json:"user_id"`
	User      *UserResponse     `json:"user,omitempty"`
	ProjectID uint              `json:"project_id"`
	Project   *ProjectResponse  `json:"project,omitempty"`
	Status    ApplicationStatus `json:"status"`
	CreatedAt string            `json:"created_at"`
}

func ToUserResponse(u *User) UserResponse {
	resp := UserResponse{
		ID:            u.ID,
		Username:      u.Username,
		FullName:      u.FullName,
		Role:          u.Role,
		IsActive:      u.IsActive,
		StudentID:     u.StudentID,
		Email:         u.Email,
		AffiliationID: u.AffiliationID,
		AvatarURL:     u.AvatarURL,
		Phone:         u.Phone,
	}
	if u.Affiliation.ID != 0 {
		resp.Affiliation = &AffiliationResponse{
			ID:          u.Affiliation.ID,
			StdName:     u.Affiliation.StdName,
			Description: u.Affiliation.Description,
		}
	}
	return resp
}

func ToProjectResponse(p *Project) ProjectResponse {
	resp := ProjectResponse{
		ID:              p.ID,
		AffiliationID:   p.AffiliationID,
		CommunityUserID: p.CommunityUserID,
		Name:            p.Name,
		Description:     p.Description,
		NumMax:          p.NumMax,
		NumAttending:    p.NumAttending,
		ProjectStartDay: p.ProjectStartDay.Format("2006-01-02T15:04:05Z07:00"),
		ProjectEndDay:   p.ProjectEndDay.Format("2006-01-02T15:04:05Z07:00"),
		FormStartDay:    p.FormStartDay.Format("2006-01-02T15:04:05Z07:00"),
		FormEndDay:      p.FormEndDay.Format("2006-01-02T15:04:05Z07:00"),
		BannerURL:       p.BannerURL,
		CreatedAt:       p.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}

	if p.DateApproved != nil {
		resp.DateApproved = p.DateApproved.Format("2006-01-02T15:04:05Z07:00")
	}

	if p.Affiliation.ID != 0 {
		resp.Affiliation = &AffiliationResponse{
			ID:          p.Affiliation.ID,
			StdName:     p.Affiliation.StdName,
			Description: p.Affiliation.Description,
		}
	}
	if p.CommunityUser.ID != 0 {
		u := ToUserResponse(&p.CommunityUser)
		resp.CommunityUser = &u
	}
	return resp
}
