package model

import "time"

type Project struct {
	ID              uint        `gorm:"primaryKey" json:"id"`
	AffiliationID   uint        `gorm:"not null;index" json:"affiliation_id"`
	Affiliation     Affiliation `gorm:"foreignKey:AffiliationID" json:"affiliation,omitempty"`
	CommunityUserID uint        `gorm:"not null;index" json:"community_user_id"`
	CommunityUser   User        `gorm:"foreignKey:CommunityUserID" json:"community_user,omitempty"`
	Name            string      `gorm:"not null" json:"name"`
	StartDay        time.Time   `gorm:"not null" json:"start_day"`
	EndDay          time.Time   `gorm:"not null" json:"end_day"`
	CreatedAt       time.Time   `json:"created_at"`
}

type StudentProject struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    uint      `gorm:"not null;index;uniqueIndex:idx_student_project" json:"user_id"`
	User      User      `gorm:"foreignKey:UserID" json:"user,omitempty"`
	ProjectID uint      `gorm:"not null;index;uniqueIndex:idx_student_project" json:"project_id"`
	Project   Project   `gorm:"foreignKey:ProjectID" json:"project,omitempty"`
	Status    string    `gorm:"type:varchar(20);not null;default:'pending'" json:"status"`
	CreatedAt time.Time `json:"created_at"`
}
