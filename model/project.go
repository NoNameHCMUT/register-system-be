package model

import "time"

type ApplicationStatus string

const (
	StatusSchoolPending    ApplicationStatus = "SCHOOL_PENDING"
	StatusSchoolReject     ApplicationStatus = "SCHOOL_REJECT"
	StatusCommunityPending ApplicationStatus = "COMMUNITY_PENDING"
	StatusCommunityReject  ApplicationStatus = "COMMUNITY_REJECT"
	StatusApproved         ApplicationStatus = "APPROVED"
)

func ValidApplicationStatuses() []ApplicationStatus {
	return []ApplicationStatus{StatusSchoolPending, StatusSchoolReject, StatusCommunityPending, StatusCommunityReject, StatusApproved}
}

type Project struct {
	ID              uint        `gorm:"primaryKey" json:"id"`
	AffiliationID   uint        `gorm:"not null;index" json:"affiliation_id"`
	Affiliation     Affiliation `gorm:"foreignKey:AffiliationID" json:"affiliation,omitempty"`
	CommunityUserID uint        `gorm:"not null;index" json:"community_user_id"`
	CommunityUser   User        `gorm:"foreignKey:CommunityUserID" json:"community_user,omitempty"`
	Name            string      `gorm:"not null" json:"name"`
	Description     string      `gorm:"type:text" json:"description"`
	NumMax          uint        `gorm:"not null;default:0" json:"num_max"`
	ProjectStartDay time.Time   `gorm:"not null" json:"project_start_day"`
	ProjectEndDay   time.Time   `gorm:"not null" json:"project_end_day"`
	FormStartDay    time.Time   `gorm:"not null" json:"form_start_day"`
	FormEndDay      time.Time   `gorm:"not null" json:"form_end_day"`
	DateApproved    *time.Time  `json:"date_approved,omitempty"`
	BannerURL       string      `json:"banner_url,omitempty"`

	NumAttending uint      `gorm:"-" json:"num_attending"`
	CreatedAt    time.Time `json:"created_at"`
}

type StudentProject struct {
	ID        uint              `gorm:"primaryKey" json:"id"`
	UserID    uint              `gorm:"not null;index;uniqueIndex:idx_student_project" json:"user_id"`
	User      User              `gorm:"foreignKey:UserID" json:"user,omitempty"`
	ProjectID uint              `gorm:"not null;index;uniqueIndex:idx_student_project" json:"project_id"`
	Project   Project           `gorm:"foreignKey:ProjectID" json:"project,omitempty"`
	Status    ApplicationStatus `gorm:"type:varchar(20);not null;default:'SCHOOL_PENDING'" json:"status"`
	CreatedAt time.Time         `json:"created_at"`
}
