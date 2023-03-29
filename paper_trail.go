package paper_trail

import (
	"reflect"

	"github.com/go-web-kits/dbx"
	"github.com/go-web-kits/utils/structx"
	"github.com/jinzhu/gorm"
)

type Trailable interface {
	Trail(*gorm.DB, interface{}, map[string]interface{}, string) error
}

type EnableTrail struct{}

func (EnableTrail) Trail(db *gorm.DB, obj interface{}, changes map[string]interface{}, who string) error {
	version := newVersionOf(db, obj, changes, who)
	return dbx.Create(version).Err
}

func Untrailable(value interface{}) bool {
	_, ok := value.(Trailable)
	return !ok
}

func Register() error {
	db := dbx.Conn()
	err := db.AutoMigrate(&Version{}).Error
	if err != nil {
		return err
	}

	db.Callback().Create().After("gorm:after_create").Register("trail:create", trailCreate)
	db.Callback().Update().After("gorm:after_update").Register("trail:update", trailUpdate)
	db.Callback().Delete().After("gorm:after_delete").Register("trail:delete", trailDestroy)
	return nil
}

func GetVersionsOf(it interface{}) ([]map[string]interface{}, error) {
	var versions []Version
	result := dbx.Where(
		&versions,
		dbx.EQ{
			"item_type": reflect.Indirect(reflect.ValueOf(it)).Type().Name(),
			"item_id":   dbx.IdOf(it),
		},
		dbx.With{Order: "created_at DESC"},
	)
	if result.Err != nil {
		return nil, result.Err
	}

	ret := []map[string]interface{}{}
	for _, v := range versions {
		m := structx.ToJsonizeMap(v)
		m["changes"] = v.ChangeSet()
		ret = append(ret, m)
	}
	return ret, nil
}

// ========

func trailCreate(scope *gorm.Scope) {
	_ = createVersionByScope(scope, "create")
}

func trailDestroy(scope *gorm.Scope) {
	_ = createVersionByScope(scope, "destroy")
}

func trailUpdate(scope *gorm.Scope) {
	_ = createVersionByScope(scope, "update")
}
