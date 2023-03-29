package test

import (
	"time"

	. "github.com/go-web-kits/paper_trail"
	"github.com/jinzhu/gorm/dialects/postgres"
	"github.com/onsi/gomega"
	"github.com/onsi/gomega/types"
	"gopkg.in/yaml.v2"
)

type h = map[string]interface{}

type Paper struct {
	Title   string         `trail:"true" json:"title" db:"title"`
	Jsonb   postgres.Jsonb `trail:"true" json:"jsonb" db:"jsonb"`
	Nilable *string        `trail:"true" json:"nilable" db:"nilable"`
	// Array *pq.StringArray `trail:"true" json:"array" db:"array"`
	Other string `json:"other" db:"other"`

	ID        uint       `json:"id" db:"id" gorm:"primary_key;index"`
	CreatedAt time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt time.Time  `json:"updated_at" db:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at" db:"deleted_at"` // `sql:"index"`
	EnableTrail
}

var EqualYAML = func(h h) types.GomegaMatcher {
	r, _ := yaml.Marshal(h)
	return gomega.Equal("---\n" + string(r))
}
