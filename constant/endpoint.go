package constant

const (
	AuthBase     = "/auth"
	AuthRegister = "/register"
	AuthLogin    = "/login"
	AuthRefresh  = "/refresh"
	AuthMe       = "/me"

	AdminBase         = "/admins"
	AdminPending      = "/users/pending"
	AdminAccept       = "/users/:id/accept"
	AdminReject       = "/users/:id/reject"
	AdminProjects     = "/projects"
	AdminApplications = "/applications"
	AdminAffiliations = "/affiliations"
	AdminAffiliation  = "/affiliations/:id"

	SchoolBase            = "/schools"
	SchoolProjects        = "/projects"
	SchoolProjectsPending = "/projects/pending"
	SchoolProjectApprove  = "/projects/:id/approve"
	SchoolApplicants      = "/projects/:id/applicants"
	SchoolApplicantAction = "/applicants/action"

	CommunityBase              = "/communities"
	CommunityProjects          = "/projects"
	CommunityProjectApplicants = "/projects/:id/applicants"
	CommunityApplicantAction   = "/applicants/action"

	StudentBase         = "/students"
	StudentProjects     = "/projects"
	StudentApply        = "/projects/:id/apply"
	StudentApplications = "/applications"

	UsersBase     = "/users"
	UsersMe       = "/me"
	UsersMeAvatar = "/me/avatar"

	ProjectBase   = "/projects"
	ProjectByID   = "/:id"
	ProjectBanner = "/:id/banner"

	AffiliationBase = "/affiliations"

	HealthCheck = "/health"
	Swagger     = "/swagger/*any"

	Uploads = "/uploads"
)
