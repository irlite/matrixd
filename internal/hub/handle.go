package hub

import (
	"bytes"
	"encoding/json"
	"fmt"
	"reflect"
	"strings"

	"github.com/irlite/matrixd/internal/wire"
)

func handle[T any](
	fn func(connName string, content T) (any, error),
) handler {
	return func(req request) wire.Outgoing {
		content, err := decode[T](req.Content)
		if err != nil {
			return fail(req, err.Error())
		}
		result, err := fn(req.ConnName, content)
		if err != nil {
			return fail(req, err.Error())
		}
		return success(req, result)
	}
}

func decode[T any](raw json.RawMessage) (T, error) {
	var content T
	if len(raw) == 0 || bytes.Equal(raw, []byte("null")) {
		raw = json.RawMessage("{}")
	}
	dec := json.NewDecoder(bytes.NewReader(raw))
	dec.DisallowUnknownFields()
	if err := dec.Decode(&content); err != nil {
		return content, err
	}
	if err := checkRequired(raw, reflect.TypeFor[T]()); err != nil {
		return content, err
	}
	return content, nil
}

func checkRequired(raw json.RawMessage, t reflect.Type) error {
	if t.Kind() != reflect.Struct || t.NumField() == 0 {
		return nil
	}
	var present map[string]json.RawMessage
	if err := json.Unmarshal(raw, &present); err != nil {
		return err
	}
	var missing []string
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		if !field.IsExported() || field.Tag.Get("optional") == "true" {
			continue
		}
		name, _, _ := strings.Cut(field.Tag.Get("json"), ",")
		if name == "-" {
			continue
		}
		if name == "" {
			name = field.Name
		}
		value, found := present[name]
		if !found || isEmpty(value, field.Type) {
			missing = append(missing, name)
		}
	}
	if len(missing) > 0 {
		return fmt.Errorf(
			"missing required: %s", strings.Join(missing, ", "),
		)
	}
	return nil
}

func isEmpty(value json.RawMessage, t reflect.Type) bool {
	if bytes.Equal(value, []byte("null")) {
		return true
	}
	return t.Kind() == reflect.String && bytes.Equal(value, []byte(`""`))
}
