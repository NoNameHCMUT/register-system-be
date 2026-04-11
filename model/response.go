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
}

type AffiliationResponse struct {
	ID      uint   `json:"id"`
	StdName string `json:"std_name"`
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
	}
	if u.Affiliation.ID != 0 {
		resp.Affiliation = &AffiliationResponse{
			ID:      u.Affiliation.ID,
			StdName: u.Affiliation.StdName,
		}
	}
	return resp
}

type ProjectResponse struct {
	ID              uint                 `json:"id"`
	Name            string               `json:"name"`
	Description     string               `json:"description"`
	NumMax          uint                 `json:"num_max"`
	NumAttending    uint                 `json:"num_attending"`
	ProjectStartDay string               `json:"project_start_day"`
	ProjectEndDay   string               `json:"project_end_day"`
	FormStartDay    string               `json:"form_start_day"`
	FormEndDay      string               `json:"form_end_day"`
	AffiliationID   uint                 `json:"affiliation_id"`
	Affiliation     *AffiliationResponse `json:"affiliation,omitempty"`
	CommunityUser   *UserResponse        `json:"community_user,omitempty"`
	CreatedAt       string               `json:"created_at"`
}

func ToProjectResponse(p *Project) ProjectResponse {
	resp := ProjectResponse{
		ID:              p.ID,
		Name:            p.Name,
		Description:     p.Description,
		NumMax:          p.NumMax,
		NumAttending:    p.NumAttending,
		ProjectStartDay: p.ProjectStartDay.Format("2006-01-02T15:04:05Z07:00"),
		ProjectEndDay:   p.ProjectEndDay.Format("2006-01-02T15:04:05Z07:00"),
		FormStartDay:    p.FormStartDay.Format("2006-01-02T15:04:05Z07:00"),
		FormEndDay:      p.FormEndDay.Format("2006-01-02T15:04:05Z07:00"),
		AffiliationID:   p.AffiliationID,
		CreatedAt:       p.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
	if p.Affiliation.ID != 0 {
		resp.Affiliation = &AffiliationResponse{
			ID:      p.Affiliation.ID,
			StdName: p.Affiliation.StdName,
		}
	}
	if p.CommunityUser.ID != 0 {
		resp.CommunityUser = &UserResponse{
			ID:            p.CommunityUser.ID,
			Username:      p.CommunityUser.Username,
			FullName:      p.CommunityUser.FullName,
			Email:         p.CommunityUser.Email,
			Role:          p.CommunityUser.Role,
		}
	}
	return resp
}
