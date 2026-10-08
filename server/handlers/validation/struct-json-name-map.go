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

const maxRecursionDepth = 10

// cloneFieldMap returns a shallow copy of the map to ensure callers cannot mutate the cached entry.
func cloneFieldMap(m map[string]string) map[string]string {
	cp := make(map[string]string, len(m))
	for k, v := range m {
		cp[k] = v
	}
	return cp
}

// getFieldMapReadOnly returns the cached field map directly without defensive cloning.
// Callers must treat the returned map as read-only.
func getFieldMapReadOnly(t reflect.Type) map[string]string {
	if v, ok := structFieldMapCache.Load(t); ok {
		return v.(map[string]string)
	}
	m := buildJSONFieldMap(t, "", "", 0)
	structFieldMapCache.Store(t, m)
	return m
}

// GetOrBuildFieldMap returns a cached field map or builds and caches it if not present.
// It returns a defensive copy to treat the internal cache as read-only.
func GetOrBuildFieldMap(t reflect.Type, parentKey, parentTag string) map[string]string {
	if parentKey == "" && parentTag == "" {
		return cloneFieldMap(getFieldMapReadOnly(t))
	}
	if v, ok := structFieldMapCache.Load(t); ok {
		return cloneFieldMap(v.(map[string]string))
	}
	m := buildJSONFieldMap(t, parentKey, parentTag, 0)
	structFieldMapCache.Store(t, m)
	return cloneFieldMap(m)
}

// buildJSONFieldMap recursively builds a map from struct namespace to json tag path,
// enforcing a maximum recursion depth to guard against cyclic/self-referential structures.
func buildJSONFieldMap(t reflect.Type, parentKey, parentTag string, depth ...int) map[string]string {
	currentDepth := 0
	if len(depth) > 0 {
		currentDepth = depth[0]
	}
	m := make(map[string]string)
	if currentDepth > maxRecursionDepth {
		return m
	}
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
		if ft.Kind() == reflect.Pointer {
			ft = ft.Elem()
		}
		if ft.Kind() == reflect.Struct && !f.Anonymous && ft.Name() != "Time" {
			for k, v := range buildJSONFieldMap(ft, key, tagPath, currentDepth+1) {
				m[k] = v
			}
		}
	}
	return m
}
