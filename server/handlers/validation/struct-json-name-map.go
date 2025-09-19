package validation

import (
	"reflect"
	"strings"

	"github.com/stfsy/go-api-kit/utils"
)

// Per-entry estimate
// - Each map entry: key (string) + value (string) + map overhead.
// - Go string header: 16 bytes.
// - Assume average key length: 20 bytes, value length: 20 bytes.
// - Each entry: (16+20) + (16+20) = 72 bytes.
// - Map overhead per entry: ~8 bytes.
// - Total per entry: ~80 bytes.
//
// Total entries
// - 500 structs × 20 fields = 10,000 entries.
//
// Total memory usage
// - 10,000 × 80 bytes = 800,000 bytes ≈ 781 KB.
//
// Add some overhead for the maps and sync.Map
// - Realistically, expect total usage to be under 1 MB.
var structFieldMapCache = utils.NewLimitedCache(500)

// GetOrBuildFieldMap returns a cached field map or builds and caches it if not present
func GetOrBuildFieldMap(t reflect.Type, parentKey, parentTag string) map[string]string {
	if v, ok := structFieldMapCache.Load(t); ok {
		return v.(map[string]string)
	}
	m := buildJSONFieldMap(t, parentKey, parentTag)
	structFieldMapCache.Store(t, m)
	return m
}

// buildJSONFieldMap recursively builds a map from struct namespace to json tag path
func buildJSONFieldMap(t reflect.Type, parentKey, parentTag string) map[string]string {
	m := make(map[string]string)
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		jsonTag := strings.Split(f.Tag.Get("json"), ",")[0]
		if jsonTag == "" || jsonTag == "-" {
			jsonTag = strings.ToLower(f.Name)
		}
		key := f.Name
		tagPath := jsonTag
		if parentKey != "" {
			key = parentKey + "." + f.Name
			tagPath = parentTag + "." + jsonTag
		}
		m[key] = tagPath
		ft := f.Type
		if ft.Kind() == reflect.Ptr {
			ft = ft.Elem()
		}
		if ft.Kind() == reflect.Struct && !f.Anonymous && ft.Name() != "Time" {
			for k, v := range buildJSONFieldMap(ft, key, tagPath) {
				m[k] = v
			}
		}
	}
	return m
}
