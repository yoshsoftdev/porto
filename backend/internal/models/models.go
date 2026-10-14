package models

import "time"

type MetaItem struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type Profile struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	Name       string    `json:"name"`
	Title      string    `json:"title"`
	Summary    string    `gorm:"type:text" json:"summary"`
	TargetRole string    `json:"target_role"`
	Email      string    `json:"email"`
	Github     string    `json:"github"`
	Linkedin   string    `json:"linkedin"`
	Strengths  []string  `gorm:"serializer:json;type:text" json:"strengths"`
	PhotoID    *uint     `json:"photo_id"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type SkillGroup struct {
	ID        uint     `gorm:"primaryKey" json:"id"`
	Category  string   `json:"category"`
	Items     []string `gorm:"serializer:json;type:text" json:"items"`
	SortOrder int      `json:"sort_order"`
}

type Project struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Slug      string         `gorm:"uniqueIndex;size:160" json:"slug"`
	Name      string         `json:"name"`
	Subtitle  string         `json:"subtitle"`
	Period    string         `json:"period"`
	Client    string         `json:"client"`
	Status    string         `json:"status"`
	Role      string         `json:"role"`
	Team      string         `json:"team"`
	Lead      string         `json:"lead"`
	Extra     []MetaItem     `gorm:"serializer:json;type:text" json:"extra"`
	Summary   string         `gorm:"type:text" json:"summary"`
	Tech      []string       `gorm:"serializer:json;type:text" json:"tech"`
	Content   string         `gorm:"type:text" json:"content"`
	SortOrder int            `json:"sort_order"`
	Hidden    bool           `json:"hidden"`
	Images    []ProjectImage `gorm:"foreignKey:ProjectID" json:"images"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
}

type ProjectImage struct {
	ID        uint   `gorm:"primaryKey" json:"id"`
	ProjectID uint   `gorm:"index" json:"project_id"`
	ImageID   uint   `json:"image_id"`
	Caption   string `json:"caption"`
	SortOrder int    `json:"sort_order"`
}

type Image struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	ContentType string    `json:"content_type"`
	Size        int       `json:"size"`
	Data        []byte    `json:"-"`
	CreatedAt   time.Time `json:"created_at"`
}

type Section struct {
	ID        uint   `gorm:"primaryKey" json:"id"`
	Slug      string `json:"slug"`
	Title     string `json:"title"`
	Content   string `gorm:"type:text" json:"content"`
	SortOrder int    `json:"sort_order"`
	Hidden    bool   `json:"hidden"`
}

type Setting struct {
	Key   string `gorm:"primaryKey;size:120"`
	Value string `gorm:"type:text"`
}

func All() []any {
	return []any{&Profile{}, &SkillGroup{}, &Project{}, &ProjectImage{}, &Image{}, &Section{}, &Setting{}}
}
