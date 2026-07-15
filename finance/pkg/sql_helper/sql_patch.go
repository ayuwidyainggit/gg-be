package sql_helper

import (
	"finance/pkg/str"
	"reflect"
)

const tagname = "sql"

type SQLPatch struct {
	Fields []string
	Args   map[string]interface{}
}

func SQLPatches(resource interface{}) SQLPatch {
	var sqlPatch SQLPatch

	// resourcePrint, _ := json.Marshal(resource)
	// log.Println("### resourcePrint ###")
	// log.Println(string(resourcePrint))
	// log.Println("### End Of resourcePrint ###")

	rType := reflect.TypeOf(resource)
	rVal := reflect.ValueOf(resource)
	n := rType.NumField()

	// log.Println("n:", n)

	sqlPatch.Fields = make([]string, 0, n)
	sqlPatch.Args = make(map[string]interface{}, 0)

	for i := 0; i < n; i++ {
		fType := rType.Field(i)
		fVal := rVal.Field(i)
		tag := fType.Tag.Get(tagname)

		// skip nil properties (not going to be patched), skip unexported fields, skip fields to be skipped for SQL
		// log.Println("fType.PkgPath:", fType.PkgPath)
		// log.Println("tag:", tag)
		// log.Println("fVal:", fVal)
		if fType.PkgPath != "" || tag == "-" || fVal.IsNil() {
			continue
		}

		// if no tag is set, use the field name
		if tag == "" {
			tag = fType.Name
		}
		// and make the tag lowercase in the end
		tag = str.Underscore(tag)

		sqlPatch.Fields = append(sqlPatch.Fields, tag+" = :"+tag)

		var val reflect.Value
		if fVal.Kind() == reflect.Ptr {
			val = fVal.Elem()
		} else {
			val = fVal
		}

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
			sqlPatch.Args[tag] = val.Pointer()
		}

		// mapValues[tag] = sqlPatch.Args[i-1]
		// log.Println("sqlPatch.Fields[i]:", tag.)

	}

	// mapValuesPrint, _ := json.Marshal(mapValues)
	// log.Println("### mapValuesPrint ###")
	// log.Println(string(mapValuesPrint))
	// log.Println("### End Of mapValuesPrint ###")

	return sqlPatch
}
