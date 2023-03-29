package paper_trail

import (
	"reflect"

	"github.com/go-web-kits/utils/structx"
	"github.com/jinzhu/gorm"
	"github.com/jinzhu/gorm/dialects/postgres"
	"gopkg.in/yaml.v2"
)

func marshalObject(scope *gorm.Scope) string {
	result, _ := yaml.Marshal(toMap(scope))
	return "---\n" + string(result)
}

func computeChanges(scope *gorm.Scope) string {
	var lastVersion Version
	scope.NewDB().LogMode(false).
		Where("item_type = ? AND item_id = ?", scope.GetModelStruct().ModelType.Name(), scope.PrimaryKeyValue()).
		Order("id DESC", true).
		First(&lastVersion)

	lastObj := lastVersion.ObjectData()
	// existLast := !scope.DB().NewRecord(lastVersion)
	currObj := toMap(scope)
	result := map[string][]interface{}{}

	for k, currV := range currObj {
		// (lastObj != nil && existLast && lastObj[k] == nil) ???
		if reflect.DeepEqual(lastObj[k], currV) {
			continue
		}

		result[k] = []interface{}{lastObj[k], currV}
	}

	r, _ := yaml.Marshal(result)
	return "---\n" + string(r)
}

func toMap(scope *gorm.Scope) map[string]interface{} {
	jMap := structx.ToJsonizeMap(scope.Value)
	m := map[string]interface{}{}
	for _, field := range scope.Fields() {
		if field.Tag.Get("trail") == "true" {
			if field.Field.Type().Name() == "Jsonb" {
				val, _ := field.Field.Interface().(postgres.Jsonb).MarshalJSON()
				if val == nil || string(val) == "null" {
					val = []byte("")
				}
				m[field.DBName] = string(val)
			} else {
				m[field.DBName] = jMap[field.Tag.Get("json")]
			}
		}
	}

	// if len(mapx.Keys(m)) == 0 {
	// 	return map[string]interface{}{"error": "please set the tag `trail:\"true\"` to the fields which you want to trail"}
	return m
}
