package constant

const (
	AuthBase     = "/auth"
	AuthRegister = "/register"
	AuthLogin    = "/login"
	AuthRefresh  = "/refresh"
	AuthMe       = "/me"

	AdminBase    = "/admin"
	AdminPending = "/users/pending"
	AdminAccept  = "/users/:id/accept"
	AdminReject  = "/users/:id/reject"

	HealthCheck = "/health"
	Swagger     = "/swagger/*any"
)
