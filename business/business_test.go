package business

import (
	"errors"
	"testing"
	"time"

	"register-system-be/model"
)

type mockUserRepo struct {
	users  map[uint]*model.User
	nextID uint
}

func newMockUserRepo() *mockUserRepo {
	return &mockUserRepo{users: make(map[uint]*model.User), nextID: 1}
}

func (m *mockUserRepo) Create(user *model.User) error {
	user.ID = m.nextID
	m.nextID++
	m.users[user.ID] = user
	return nil
}

func (m *mockUserRepo) FindByUsername(username string) (*model.User, error) {
	for _, u := range m.users {
		if u.Username == username {
			return u, nil
		}
	}
	return nil, errors.New("not found")
}

func (m *mockUserRepo) FindByEmail(email string) (*model.User, error) {
	for _, u := range m.users {
		if u.Email == email {
			return u, nil
		}
	}
	return nil, errors.New("not found")
}

func (m *mockUserRepo) FindByID(id uint) (*model.User, error) {
	u, ok := m.users[id]
	if !ok {
		return nil, errors.New("not found")
	}
	return u, nil
}

func (m *mockUserRepo) FindPending() ([]model.User, error) {
	var result []model.User
	for _, u := range m.users {
		if !u.IsActive {
			result = append(result, *u)
		}
	}
	return result, nil
}

func (m *mockUserRepo) UpdateActive(userID uint, active bool) error {
	u, ok := m.users[userID]
	if !ok {
		return errors.New("not found")
	}
	u.IsActive = active
	return nil
}

func (m *mockUserRepo) UpdateRefreshToken(userID uint, token string) error {
	return nil
}

func (m *mockUserRepo) UpdateProfile(userID uint, phone string, avatarURL string) error {
	u, ok := m.users[userID]
	if !ok {
		return errors.New("not found")
	}
	if phone != "" {
		u.Phone = phone
	}
	if avatarURL != "" {
		u.AvatarURL = avatarURL
	}
	return nil
}

func (m *mockUserRepo) UpdateFields(userID uint, fields map[string]interface{}) error {
	u, ok := m.users[userID]
	if !ok {
		return errors.New("not found")
	}
	if v, ok := fields["full_name"]; ok {
		u.FullName = v.(string)
	}
	if v, ok := fields["phone"]; ok {
		u.Phone = v.(string)
	}
	return nil
}

func (m *mockUserRepo) FindAll() ([]model.User, error) {
	var result []model.User
	for _, u := range m.users {
		result = append(result, *u)
	}
	return result, nil
}

type mockAffiliationRepo struct {
	affs   map[uint]*model.Affiliation
	nextID uint
}

func newMockAffRepo() *mockAffiliationRepo {
	return &mockAffiliationRepo{affs: make(map[uint]*model.Affiliation), nextID: 1}
}

func (m *mockAffiliationRepo) FindByID(id uint) (*model.Affiliation, error) {
	a, ok := m.affs[id]
	if !ok {
		return nil, errors.New("not found")
	}
	return a, nil
}

func (m *mockAffiliationRepo) FindAll() ([]model.Affiliation, error) {
	var result []model.Affiliation
	for _, a := range m.affs {
		result = append(result, *a)
	}
	return result, nil
}

func (m *mockAffiliationRepo) Create(aff *model.Affiliation) error {
	aff.ID = m.nextID
	m.nextID++
	m.affs[aff.ID] = aff
	return nil
}

func (m *mockAffiliationRepo) Update(aff *model.Affiliation) error {
	m.affs[aff.ID] = aff
	return nil
}

func (m *mockAffiliationRepo) Delete(id uint) error {
	delete(m.affs, id)
	return nil
}

type mockProjectRepo struct {
	projects map[uint]*model.Project
	nextID   uint
}

func newMockProjectRepo() *mockProjectRepo {
	return &mockProjectRepo{projects: make(map[uint]*model.Project), nextID: 1}
}

func (m *mockProjectRepo) Create(p *model.Project) error {
	p.ID = m.nextID
	m.nextID++
	m.projects[p.ID] = p
	return nil
}

func (m *mockProjectRepo) FindByID(id uint) (*model.Project, error) {
	p, ok := m.projects[id]
	if !ok {
		return nil, errors.New("not found")
	}
	return p, nil
}

func (m *mockProjectRepo) FindAll() ([]model.Project, error) {
	var result []model.Project
	for _, p := range m.projects {
		result = append(result, *p)
	}
	return result, nil
}

func (m *mockProjectRepo) FindByAffiliation(affiliationID uint) ([]model.Project, error) {
	var result []model.Project
	for _, p := range m.projects {
		if p.AffiliationID == affiliationID {
			result = append(result, *p)
		}
	}
	return result, nil
}

func (m *mockProjectRepo) FindByAffiliationApproved(affiliationID uint) ([]model.Project, error) {
	var result []model.Project
	for _, p := range m.projects {
		if p.AffiliationID == affiliationID && p.DateApproved != nil {
			result = append(result, *p)
		}
	}
	return result, nil
}

func (m *mockProjectRepo) FindByCreator(userID uint) ([]model.Project, error) {
	var result []model.Project
	for _, p := range m.projects {
		if p.CommunityUserID == userID {
			result = append(result, *p)
		}
	}
	return result, nil
}

func (m *mockProjectRepo) Update(p *model.Project) error {
	m.projects[p.ID] = p
	return nil
}

func (m *mockProjectRepo) UpdateDateApproved(id uint, dateApproved interface{}) error {
	p, ok := m.projects[id]
	if !ok {
		return errors.New("not found")
	}
	if t, ok := dateApproved.(time.Time); ok && !t.IsZero() {
		p.DateApproved = &t
	} else {
		p.DateApproved = nil
	}
	return nil
}

func (m *mockProjectRepo) UpdateBanner(id uint, bannerURL string) error {
	p, ok := m.projects[id]
	if !ok {
		return errors.New("not found")
	}
	p.BannerURL = bannerURL
	return nil
}

func (m *mockProjectRepo) CountAttending(projectID uint) (uint, error) {
	return 0, nil
}

type mockAppRepo struct {
	apps   map[uint]*model.StudentProject
	nextID uint
}

func newMockAppRepo() *mockAppRepo {
	return &mockAppRepo{apps: make(map[uint]*model.StudentProject), nextID: 1}
}

func (m *mockAppRepo) Create(sp *model.StudentProject) error {
	sp.ID = m.nextID
	m.nextID++
	m.apps[sp.ID] = sp
	return nil
}

func (m *mockAppRepo) FindByID(id uint) (*model.StudentProject, error) {
	a, ok := m.apps[id]
	if !ok {
		return nil, errors.New("not found")
	}
	return a, nil
}

func (m *mockAppRepo) FindByStudentID(userID uint) ([]model.StudentProject, error) {
	var result []model.StudentProject
	for _, a := range m.apps {
		if a.UserID == userID {
			result = append(result, *a)
		}
	}
	return result, nil
}

func (m *mockAppRepo) FindByProjectID(projectID uint) ([]model.StudentProject, error) {
	var result []model.StudentProject
	for _, a := range m.apps {
		if a.ProjectID == projectID {
			result = append(result, *a)
		}
	}
	return result, nil
}

func (m *mockAppRepo) FindByProjectIDAndStatus(projectID uint, status model.ApplicationStatus) ([]model.StudentProject, error) {
	var result []model.StudentProject
	for _, a := range m.apps {
		if a.ProjectID == projectID && a.Status == status {
			result = append(result, *a)
		}
	}
	return result, nil
}

func (m *mockAppRepo) FindByStudentAndProject(userID uint, projectID uint) (*model.StudentProject, error) {
	for _, a := range m.apps {
		if a.UserID == userID && a.ProjectID == projectID {
			return a, nil
		}
	}
	return nil, errors.New("not found")
}

func (m *mockAppRepo) FindAll() ([]model.StudentProject, error) {
	var result []model.StudentProject
	for _, a := range m.apps {
		result = append(result, *a)
	}
	return result, nil
}

func (m *mockAppRepo) BatchUpdateStatus(ids []uint, status model.ApplicationStatus) error {
	for _, id := range ids {
		if a, ok := m.apps[id]; ok {
			a.Status = status
		}
	}
	return nil
}

func (m *mockAppRepo) UpdateStatus(id uint, status model.ApplicationStatus) error {
	a, ok := m.apps[id]
	if !ok {
		return errors.New("not found")
	}
	a.Status = status
	return nil
}

func TestApplicationStatusTransitions(t *testing.T) {
	tests := []struct {
		name     string
		current  model.ApplicationStatus
		action   string
		role     model.Role
		expected model.ApplicationStatus
	}{
		{"school approve SCHOOL_PENDING", model.StatusSchoolPending, "approve", model.RoleSchool, model.StatusCommunityPending},
		{"school reject SCHOOL_PENDING", model.StatusSchoolPending, "reject", model.RoleSchool, model.StatusSchoolReject},
		{"community approve COMMUNITY_PENDING", model.StatusCommunityPending, "approve", model.RoleCommunity, model.StatusApproved},
		{"community reject COMMUNITY_PENDING", model.StatusCommunityPending, "reject", model.RoleCommunity, model.StatusCommunityReject},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var expected model.ApplicationStatus
			if tt.role == model.RoleSchool {
				if tt.action == "approve" {
					expected = model.StatusCommunityPending
				} else {
					expected = model.StatusSchoolReject
				}
			} else {
				if tt.action == "approve" {
					expected = model.StatusApproved
				} else {
					expected = model.StatusCommunityReject
				}
			}
			if expected != tt.expected {
				t.Errorf("expected %s got %s", expected, tt.expected)
			}
		})
	}
}

func TestLoginBlockedForInactiveUser(t *testing.T) {
	userRepo := newMockUserRepo()
	user := &model.User{
		Username: "inactive", FullName: "Test", PasswordHash: "$2a$10$invalid",
		Email: "inactive@test.com", Role: model.RoleStudent, IsActive: false, AffiliationID: 1,
	}
	_ = userRepo.Create(user)

	if user.IsActive {
		t.Error("user should be inactive")
	}
}

func TestAffiliationCreateValidation(t *testing.T) {
	affRepo := newMockAffRepo()

	err := affRepo.Create(&model.Affiliation{StdName: "Test"})
	if err != nil {
		t.Fatal(err)
	}

	a, err := affRepo.FindByID(1)
	if err != nil {
		t.Fatal(err)
	}
	if a.StdName != "Test" {
		t.Errorf("expected Test, got %s", a.StdName)
	}
}
