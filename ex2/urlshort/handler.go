package urlshort

import (
	"encoding/json"
	"net/http"

	"gopkg.in/yaml.v2"
)

// MapHandler will return an http.HandlerFunc (which also
// implements http.Handler) that will attempt to map any
// paths (keys in the map) to their corresponding URL (values
// that each key in the map points to, in string format).
// If the path is not provided in the map, then the fallback
// http.Handler will be called instead.
func MapHandler(pathsToUrls map[string]string, fallback http.Handler) (http.HandlerFunc, error) {
	//	TODO: Implement this...

	handler := func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		url, ok := pathsToUrls[path]
		switch ok {
		case true:
			http.Redirect(w, r, url, http.StatusSeeOther)
		default:
			fallback.ServeHTTP(w, r)
		}
	}

	return handler, nil
}

// YAMLHandler will parse the provided YAML and then return
// an http.HandlerFunc (which also implements http.Handler)
// that will attempt to map any paths to their corresponding
// URL. If the path is not provided in the YAML, then the
// fallback http.Handler will be called instead.
//
// YAML is expected to be in the format:
//
//     - path: /some-path
//       url: https://www.some-url.com/demo
//
// The only errors that can be returned all related to having
// invalid YAML data.
//
// See MapHandler to create a similar http.HandlerFunc via
// a mapping of paths to urls.
func YAMLHandler(yml []byte, fallback http.Handler) (http.HandlerFunc, error) {

	// parse yaml
	var pathUrls []pathURL
	err := yaml.Unmarshal(yml, &pathUrls)
	if err != nil {
		return nil, err
	}

	// convert yaml to map
	pathMap, err := buildMap(pathUrls)
	if err != nil {
		return nil, err
	}

	handler, err := MapHandler(pathMap, fallback)

	if err != nil {
		return nil, err
	}

	return handler, nil
}

// JSONHandler is a handler similar to YAMLHandler but parses JSON files
func JSONHandler(jsonData []byte, fallback http.Handler) (http.HandlerFunc, error) {

	// parse json
	var pathMap map[string]string
	err := json.Unmarshal(jsonData, &pathMap)
	if err != nil {
		return nil, err
	}

	handler, err := MapHandler(pathMap, fallback)

	if err != nil {
		return nil, err
	}

	return handler, nil
}

type pathURL struct {
	Path string
	URL  string
}

func buildMap(pathUrls []pathURL) (map[string]string, error) {
	pathMap := make(map[string]string)
	for _, pu := range pathUrls {
		pathMap[pu.Path] = pu.URL
	}

	return pathMap, nil
}
