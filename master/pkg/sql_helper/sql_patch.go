package sql_helper

import (
	"master/pkg/str"
	"math"
	"reflect"

	"github.com/lib/pq"
)

const tagname = "sql"

type SQLPatch struct {
	Fields []string
	Args   map[string]interface{}
}

func SQLPatches(resource interface{}) SQLPatch {
	// log.Info("SQLPatches")
	var sqlPatch SQLPatch
	// resourcePrint, _ := json.Marshal(resource)
	// log.Info("### resourcePrint ###")
	// log.Info(string(resourcePrint))
	// log.Info("### End Of resourcePrint ###")

	rType := reflect.TypeOf(resource)
	rVal := reflect.ValueOf(resource)
	n := rType.NumField()
	// log.Info("n:", n)

	sqlPatch.Fields = make([]string, 0, n)
	sqlPatch.Args = make(map[string]interface{}, 0)

	for i := 0; i < n; i++ {
		fType := rType.Field(i)
		fVal := rVal.Field(i)
		tag := fType.Tag.Get(tagname)

		// skip nil properties (not going to be patched), skip unexported fields, skip fields to be skipped for SQL
		// log.Info("fType:", fType)
		// log.Info("tag:", tag)
		// log.Info("fVal:", &fVal)
		if fType.PkgPath != "" || tag == "-" {
			continue
		}

		// Check if the value is nil (only for nullable types)
		switch fVal.Kind() {
		case reflect.Chan, reflect.Func, reflect.Map, reflect.Ptr, reflect.UnsafePointer, reflect.Interface, reflect.Slice:
			if fVal.IsNil() {
				continue
			}
		}

		// if no tag is set, use the field name
		if tag == "" {
			tag = fType.Name
		}
		// log.Info("if no tag is set, use the field name, tag: ", tag)

		// and make the tag lowercase in the end
		tag = str.Underscore(tag)
		// log.Info("make the tag lowercase in the end, tag: ", tag)

		sqlPatch.Fields = append(sqlPatch.Fields, tag+" = :"+tag)
		// log.Info(tag)
		var val reflect.Value
		if fVal.Kind() == reflect.Ptr {
			val = fVal.Elem()
		} else {
			val = fVal
		}
		// log.Infof("i: %d, tag: %v, val: %v", i, tag, val)
		switch val.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			sqlPatch.Args[tag] = val.Int()
		case reflect.Float32, reflect.Float64:
			sqlPatch.Args[tag] = val.Float()
		case reflect.String:
			if val.String() != "" {
				sqlPatch.Args[tag] = val.String()
			} else {
				sqlPatch.Args[tag] = nil
			}
		case reflect.Bool:
			if val.Bool() {
				sqlPatch.Args[tag] = true
			} else {
				sqlPatch.Args[tag] = false
			}
		case reflect.Struct:
			sqlPatch.Args[tag] = val.Interface()
		case reflect.Slice:
			// log.Info("tag:", tag)
			// log.Info("fVal:", fVal)
			// log.Info("val.Pointer():", val.Pointer())
			sqlPatch.Args[tag] = pq.Array(val.Interface())
		}

		// mapValues[tag] = sqlPatch.Args[i-1]
		// log.Infof("sqlPatch.Fields[%v]: %v", tag, val)

	}

	return sqlPatch
}

// CalculateLastPage calculates the last page number for pagination
func CalculateLastPage(total, limit int) int {
	if limit <= 0 {
		return 1
	}
	return int(math.Ceil(float64(total) / float64(limit)))
}
