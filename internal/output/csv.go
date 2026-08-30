package output

import (
	"encoding/csv"
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"
	"sysinfo3-go/internal/collector"
	"time"
)

func RenderCSV(s *collector.SystemSnapshot) {
	w := csv.NewWriter(os.Stdout)
	w.Write([]string{"key", "value"})

	flatten(w, "", reflect.ValueOf(*s))

	w.Flush()
	if err := w.Error(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: CSV encoding failed: %v\n", err)
		os.Exit(1)
	}
}

func flatten(w *csv.Writer, prefix string, v reflect.Value) {
	if v.Kind() == reflect.Ptr {
		if v.IsNil() {
			return
		}
		v = v.Elem()
	}

	t := v.Type()

	switch v.Kind() {
	case reflect.Struct:
		if v.Type() == reflect.TypeOf(collector.SystemSnapshot{}) || v.Type().String() == "collector.SystemSnapshot" {
			for i := 0; i < v.NumField(); i++ {
				field := t.Field(i)
				if field.Name == "Timestamp" {
					if ts, ok := v.Field(i).Interface().(time.Time); ok {
						w.Write([]string{"timestamp", ts.Format(time.RFC3339)})
					} else {
						w.Write([]string{"timestamp", fmt.Sprintf("%v", v.Field(i).Interface())})
					}
					continue
				}
				key := field.Tag.Get("json")
				if key == "" || key == "-" {
					key = strings.ToLower(field.Name)
				}
				keyParts := strings.Split(key, ",")
				key = keyParts[0]

				flatten(w, key, v.Field(i))
			}
			return
		}
		for i := 0; i < v.NumField(); i++ {
			field := t.Field(i)
			key := field.Tag.Get("json")
			if key == "" || key == "-" {
				key = strings.ToLower(field.Name)
			}
			keyParts := strings.Split(key, ",")
			key = keyParts[0]
			if prefix != "" {
				key = prefix + "." + key
			}
			flatten(w, key, v.Field(i))
		}

	case reflect.Slice:
		for i := 0; i < v.Len(); i++ {
			elemPrefix := fmt.Sprintf("%s[%d]", prefix, i)
			flatten(w, elemPrefix, v.Index(i))
		}

	case reflect.String:
		w.Write([]string{prefix, v.String()})
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		w.Write([]string{prefix, strconv.FormatInt(v.Int(), 10)})
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		w.Write([]string{prefix, strconv.FormatUint(v.Uint(), 10)})
	case reflect.Float32, reflect.Float64:
		w.Write([]string{prefix, strconv.FormatFloat(v.Float(), 'f', 2, 64)})
	default:
		w.Write([]string{prefix, fmt.Sprintf("%v", v.Interface())})
	}
}
