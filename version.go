package paper_trail

import (
	"time"

	"github.com/go-web-kits/utils/mapx"
	"github.com/jinzhu/gorm"
	"gopkg.in/yaml.v2"
)

// Version is a main entity, which used to trail changes.
type Version struct {
	ID            uint      `json:"-"     db:"id"            gorm:"primary_key;index"`
	CreatedAt     time.Time `json:"time"  db:"created_at"`
	Event         string    `json:"event" db:"event"`
	ItemID        uint      `json:"-"     db:"item_id"`
	ItemType      string    `json:"-"     db:"item_type"`
	Object        string    `json:"-"     db:"object"`
	ObjectChanges string    `json:"-"     db:"object_changes"`
	Whodunnit     string    `json:"admin" db:"whodunnit"`
}

func (v Version) ChangeSet() map[string][]interface{} {
	result := map[string][]interface{}{}
	_ = yaml.Unmarshal([]byte(v.ObjectChanges), &result)
	return result
}

func (v Version) ObjectData() map[string]interface{} {
	result := map[string]interface{}{}
	_ = yaml.Unmarshal([]byte(v.Object), &result)
	return result
}

func newVersionOf(db *gorm.DB, obj interface{}, changes map[string]interface{}, who string) *Version {
	scope := db.NewScope(obj)
	c, _ := yaml.Marshal(changes)
	if who == "" {
		who = "admin"
	}

	return &Version{
		Event:         "update " + mapx.Keys(changes)[0],
		ItemID:        scope.PrimaryKeyValue().(uint),
		ItemType:      scope.GetModelStruct().ModelType.Name(),
		Object:        marshalObject(scope),
		ObjectChanges: "---\n" + string(c),
		Whodunnit:     who,
	}
}

func newVersionByScope(scope *gorm.Scope, event string) *Version {
	version := Version{
		Event:         event,
		ItemID:        scope.PrimaryKeyValue().(uint),
		ItemType:      scope.GetModelStruct().ModelType.Name(),
		Object:        marshalObject(scope),
		ObjectChanges: computeChanges(scope),
	}

	if who, setted := scope.DB().Get("trail:who"); setted {
		version.Whodunnit = who.(string)
	}
	if version.Whodunnit == "" {
		version.Whodunnit = "admin"
	}

	return &version
}

func createVersionByScope(scope *gorm.Scope, action string) error {
	if Untrailable(scope.Value) {
		return nil
	}

	version := newVersionByScope(scope, action)
	// No Changes
	if version.Event != "destroy" && version.ObjectChanges == "---\n{}\n" {
		return nil
	}

	return scope.DB().LogMode(false).Create(version).Error
}
