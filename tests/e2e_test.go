package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"register-system-be/business"
	"register-system-be/config"
	"register-system-be/email"
	"register-system-be/handler"
	"register-system-be/model"
	"register-system-be/repo"
	"register-system-be/router"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	db             *gorm.DB
	cfg            *config.Config
	r              *gin.Engine
	authHandler    *handler.AuthHandler
	adminHandler   *handler.AdminHandler
	projectHandler *handler.ProjectHandler
	appHandler     *handler.ApplicationHandler
)

func TestMain(m *testing.M) {
	if os.Getenv("TEST_DB_HOST") != "" {
		os.Setenv("DB_HOST", os.Getenv("TEST_DB_HOST"))
	}
	if os.Getenv("TEST_DB_NAME") != "" {
		os.Setenv("DB_NAME", os.Getenv("TEST_DB_NAME"))
	}

	cfg = config.Load()

	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName,
	)
	var err error
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Println("Cannot connect to test DB, skipping e2e tests:", err)
		os.Exit(0)
	}

	db.Migrator().DropTable(&model.StudentProject{}, &model.Project{}, &model.User{}, &model.Affiliation{})
	repo.Migrate(db)

	userRepo := repo.NewUserRepo(db)
	affRepo := repo.NewAffiliationRepo(db)
	projRepo := repo.NewProjectRepo(db)
	appRepo := repo.NewApplicationRepo(db)

	emailSender := email.NewNoopSender()

	authBiz := business.NewAuthBusiness(userRepo, affRepo, cfg, emailSender)
	adminBiz := business.NewAdminBusiness(userRepo, projRepo, appRepo, affRepo, emailSender)
	affBiz := business.NewAffiliationBusiness(affRepo)
	projBiz := business.NewProjectBusiness(projRepo, affRepo, userRepo)
	appBiz := business.NewApplicationBusiness(appRepo, projRepo, userRepo, emailSender)

	authHandler = handler.NewAuthHandler(authBiz, cfg.UploadDir)
	adminHandler = handler.NewAdminHandler(adminBiz)
	affHandler := handler.NewAffiliationHandler(affBiz)
	projectHandler = handler.NewProjectHandler(projBiz, cfg.UploadDir)
	appHandler = handler.NewApplicationHandler(appBiz)

	gin.SetMode(gin.TestMode)
	r = router.Setup(authHandler, adminHandler, affHandler, projectHandler, appHandler, cfg, userRepo)

	code := m.Run()
	os.Exit(code)
}

func resetDB(t *testing.T) {
	t.Helper()
	db.Exec("TRUNCATE TABLE student_projects, projects, users, affiliations RESTART IDENTITY CASCADE")
}

func doRequest(t *testing.T, method, path string, body interface{}, token string) *httptest.ResponseRecorder {
	t.Helper()
	var bodyReader io.Reader
	if body != nil {
		b, _ := json.Marshal(body)
		bodyReader = bytes.NewReader(b)
	}
	req, _ := http.NewRequest(method, "http://localhost"+path, bodyReader)
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

func doMultipartRequest(t *testing.T, method, path, fieldName, filename string, fileContent []byte, token string) *httptest.ResponseRecorder {
	t.Helper()
	var buf bytes.Buffer
	w := multipart.NewWriter(&buf)
	part, _ := w.CreateFormFile(fieldName, filename)
	part.Write(fileContent)
	w.Close()
	req, _ := http.NewRequest(method, "http://localhost"+path, &buf)
	req.Header.Set("Content-Type", w.FormDataContentType())
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	w2 := httptest.NewRecorder()
	r.ServeHTTP(w2, req)
	return w2
}

func parseResponse(t *testing.T, w *httptest.ResponseRecorder) map[string]interface{} {
	t.Helper()
	var result map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &result)
	return result
}

func getToken(t *testing.T, username, password string) string {
	t.Helper()
	w := doRequest(t, "POST", "/auth/login", model.LoginRequest{Username: username, Password: password}, "")
	resp := parseResponse(t, w)
	data, ok := resp["data"].(map[string]interface{})
	if !ok {
		t.Fatalf("login failed for %s: %v", username, resp)
	}
	return data["access_token"].(string)
}

func createTestAdmin(t *testing.T) (uint, string) {
	t.Helper()
	hash, _ := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)
	admin := &model.User{
		Username: "admin_root", FullName: "Admin", PasswordHash: string(hash),
		Email: "admin@test.com", Role: model.RoleAdmin, IsActive: true, AffiliationID: 1,
	}
	db.Create(admin)
	return admin.ID, getToken(t, "admin_root", "admin123")
}

func TestE2EFullFlow(t *testing.T) {
	resetDB(t)

	// 1. Admin pre-created with affiliation
	aff := &model.Affiliation{StdName: "Admin Org"}
	db.Create(aff)
	_, adminToken := createTestAdmin(t)

	// 2. Admin adds 2 affiliations (1 school, 1 community)
	w := doRequest(t, "POST", "/admin/affiliations", model.AffiliationCreateRequest{StdName: "BK University", Description: "School"}, adminToken)
	if w.Code != 201 {
		t.Fatalf("create school affiliation: %d %s", w.Code, w.Body.String())
	}
	w = doRequest(t, "POST", "/admin/affiliations", model.AffiliationCreateRequest{StdName: "Binh Phuoc Community", Description: "Community"}, adminToken)
	if w.Code != 201 {
		t.Fatalf("create community affiliation: %d %s", w.Code, w.Body.String())
	}

	// 3. School, community, student register
	schoolAffID := uint(2)
	communityAffID := uint(3)

	w = doRequest(t, "POST", "/auth/register", model.RegisterRequest{
		Username: "school_user", FullName: "School User", Password: "school456",
		Email: "school@test.com", AffiliationID: schoolAffID, Role: model.RoleSchool,
	}, "")
	if w.Code != 201 {
		t.Fatalf("register school: %d %s", w.Code, w.Body.String())
	}

	w = doRequest(t, "POST", "/auth/register", model.RegisterRequest{
		Username: "community_user", FullName: "Community User", Password: "community789",
		Email: "community@test.com", AffiliationID: communityAffID, Role: model.RoleCommunity,
	}, "")
	if w.Code != 201 {
		t.Fatalf("register community: %d %s", w.Code, w.Body.String())
	}

	w = doRequest(t, "POST", "/auth/register", model.RegisterRequest{
		Username: "student_user", FullName: "Student User", Password: "student123",
		Email: "student@test.com", AffiliationID: schoolAffID, Role: model.RoleStudent,
	}, "")
	if w.Code != 201 {
		t.Fatalf("register student: %d %s", w.Code, w.Body.String())
	}

	// Inactive users can't login
	w = doRequest(t, "POST", "/auth/login", model.LoginRequest{Username: "school_user", Password: "school456"}, "")
	if w.Code != 401 {
		t.Fatalf("inactive user should not login: %d", w.Code)
	}

	// 4. Admin approves accounts
	var pendingUsers []model.User
	db.Where("is_active = ?", false).Find(&pendingUsers)
	for _, u := range pendingUsers {
		doRequest(t, "POST", fmt.Sprintf("/admin/users/%d/accept", u.ID), nil, adminToken)
	}

	schoolToken := getToken(t, "school_user", "school456")
	communityToken := getToken(t, "community_user", "community789")
	studentToken := getToken(t, "student_user", "student123")

	// 5. Community creates a project under schoolAffID so students can see it
	now := time.Now()
	formStart := now.Add(-1 * time.Hour).Format(time.RFC3339)
	formEnd := now.Add(72 * time.Hour).Format(time.RFC3339)
	projStart := now.Add(168 * time.Hour).Format(time.RFC3339)
	projEnd := now.Add(336 * time.Hour).Format(time.RFC3339)

	w = doRequest(t, "POST", "/projects", model.ProjectCreateRequest{
		AffiliationID: schoolAffID, Name: "School Project", Description: "For students",
		NumMax: 10, ProjectStartDay: projStart, ProjectEndDay: projEnd,
		FormStartDay: formStart, FormEndDay: formEnd,
	}, communityToken)
	if w.Code != 201 {
		t.Fatalf("create project: %d %s", w.Code, w.Body.String())
	}

	// 6. Community views project -> sees it hasn't been approved
	w = doRequest(t, "GET", "/community/projects", nil, communityToken)
	resp := parseResponse(t, w)
	data := resp["data"].([]interface{})
	proj := data[0].(map[string]interface{})
	if proj["date_approved"] != nil {
		t.Fatal("project should not be approved yet")
	}

	// 7. Student views their affiliation's projects -> sees nothing (not approved)
	w = doRequest(t, "GET", "/students/projects", nil, studentToken)
	resp = parseResponse(t, w)
	data = resp["data"].([]interface{})
	if len(data) != 0 {
		t.Fatal("student should see no projects before school approval")
	}

	// 8. School approves the project
	w = doRequest(t, "POST", "/schools/projects/1/approve", nil, schoolToken)
	if w.Code != 200 {
		t.Fatalf("school approve project: %d %s", w.Code, w.Body.String())
	}

	// 9. Community views project -> sees it has been approved
	w = doRequest(t, "GET", "/community/projects", nil, communityToken)
	resp = parseResponse(t, w)
	data = resp["data"].([]interface{})
	proj = data[0].(map[string]interface{})
	if proj["date_approved"] == nil {
		t.Fatal("project should be approved now")
	}

	// 10. Student views their affiliation's projects -> sees the project
	w = doRequest(t, "GET", "/students/projects", nil, studentToken)
	resp = parseResponse(t, w)
	data = resp["data"].([]interface{})
	if len(data) == 0 {
		t.Fatal("student should see approved projects from their affiliation")
	}

	// 11. School opens the project -> view its registration list, see no applicant
	w = doRequest(t, "GET", "/schools/projects/1/applicants", nil, schoolToken)
	resp = parseResponse(t, w)
	data = resp["data"].([]interface{})
	if len(data) != 0 {
		t.Fatal("should have no applicants yet")
	}

	// 12. Student registers to the project
	w = doRequest(t, "POST", "/students/projects/1/apply", nil, studentToken)
	if w.Code != 201 {
		t.Fatalf("student apply: %d %s", w.Code, w.Body.String())
	}

	// 13. School opens the project -> sees the application with SCHOOL_PENDING
	w = doRequest(t, "GET", "/schools/projects/1/applicants", nil, schoolToken)
	resp = parseResponse(t, w)
	data = resp["data"].([]interface{})
	if len(data) == 0 {
		t.Fatal("school should see applicant")
	}
	applicant := data[0].(map[string]interface{})
	if applicant["status"] != "SCHOOL_PENDING" {
		t.Fatalf("expected SCHOOL_PENDING, got %v", applicant["status"])
	}

	// 14. Community opens the project -> sees no applicant (SCHOOL_PENDING not visible)
	w = doRequest(t, "GET", "/community/projects/1/applicants", nil, communityToken)
	resp = parseResponse(t, w)
	data = resp["data"].([]interface{})
	if len(data) != 0 {
		t.Fatal("community should not see SCHOOL_PENDING applicants")
	}

	// 15. School approves the student applicant
	w = doRequest(t, "POST", "/schools/applicants/action", model.ApplicationActionRequest{
		ApplicationIDs: []uint{1}, Action: "approve",
	}, schoolToken)
	if w.Code != 200 {
		t.Fatalf("school approve applicant: %d %s", w.Code, w.Body.String())
	}

	// 16. Student views their application list -> COMMUNITY_PENDING
	w = doRequest(t, "GET", "/students/applications", nil, studentToken)
	resp = parseResponse(t, w)
	data = resp["data"].([]interface{})
	if len(data) == 0 {
		t.Fatal("student should see applications")
	}
	studentApp := data[0].(map[string]interface{})
	if studentApp["status"] != "COMMUNITY_PENDING" {
		t.Fatalf("expected COMMUNITY_PENDING, got %v", studentApp["status"])
	}

	// 17. Community opens the project -> sees applicant
	w = doRequest(t, "GET", "/community/projects/1/applicants", nil, communityToken)
	resp = parseResponse(t, w)
	data = resp["data"].([]interface{})
	if len(data) == 0 {
		t.Fatal("community should see COMMUNITY_PENDING applicant")
	}

	// 18. Community accepts the application
	w = doRequest(t, "POST", "/community/applicants/action", model.ApplicationActionRequest{
		ApplicationIDs: []uint{1}, Action: "approve",
	}, communityToken)
	if w.Code != 200 {
		t.Fatalf("community approve applicant: %d %s", w.Code, w.Body.String())
	}

	// 19. Student views their application list -> APPROVED
	w = doRequest(t, "GET", "/students/applications", nil, studentToken)
	resp = parseResponse(t, w)
	data = resp["data"].([]interface{})
	studentApp = data[0].(map[string]interface{})
	if studentApp["status"] != "APPROVED" {
		t.Fatalf("expected APPROVED, got %v", studentApp["status"])
	}

	// 20. School views their application list -> APPROVED
	w = doRequest(t, "GET", "/schools/projects/1/applicants", nil, schoolToken)
	resp = parseResponse(t, w)
	data = resp["data"].([]interface{})
	schoolApp := data[0].(map[string]interface{})
	if schoolApp["status"] != "APPROVED" {
		t.Fatalf("expected APPROVED, got %v", schoolApp["status"])
	}
}
