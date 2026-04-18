package constant

const (
	AuthBase     = "/auth"
	AuthRegister = "/register"
	AuthLogin    = "/login"
	AuthRefresh  = "/refresh"
	AuthMe       = "/me"

	AdminBase         = "/admins"
	AdminPending      = "/users/pending"
	AdminUsersActive  = "/users/active"
	AdminAccept       = "/users/:id/accept"
	AdminReject       = "/users/:id/reject"
	AdminProjects     = "/projects"
	AdminApplications = "/applications"
	AdminAffiliations = "/affiliations"
	AdminAffiliation  = "/affiliations/:id"

	AffiliationBase = "/affiliations"

	ProjectBase   = "/projects"
	ProjectMyList = "/my-list"
	ProjectByID   = "/:id"

	HealthCheck = "/health"
	Swagger     = "/swagger/*any"
)
