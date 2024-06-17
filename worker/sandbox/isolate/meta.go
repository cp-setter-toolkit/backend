package isolate

import (
	"bufio"
	"os"
	"strconv"
	"strings"

	"github.com/thepluck/cp-setter-toolkit/helper/errors"
)

// NewMetaFile creates a new meta file in the given directory.
func NewMetaFile(dir string) (string, error) {
	metaFile, err := os.CreateTemp(dir, "isolate-metafile-*")
	if err != nil {
		return "", errors.Wrap(err, "creating meta file")
	}
	return metaFile.Name(), nil
}

// Meta contains the fields of the meta file.
type Meta struct {
	Fields map[string]string
}

// Read reads the meta file and returns the Meta object.
func Read(path string) (*Meta, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	defer f.Close()

	sc := bufio.NewScanner(f)
	fields := make(map[string]string)

	for sc.Scan() {
		lst := strings.SplitN(sc.Text(), ":", 2)
		if len(lst) < 2 {
			continue
		}
		fields[lst[0]] = lst[1]
	}

	return &Meta{fields}, nil
}

// Int returns the integer value of the field.
func (m *Meta) Int(key string) (int, error) {
	if v, ok := m.Fields[key]; ok {
		n, err := strconv.Atoi(v)
		if err != nil {
			return 0, errors.Wrapf(err, "parsing key %s", key)
		}
		return n, nil
	}
	return 0, errors.Errorf("key %s not found", key)
}

// String returns the string value of the field.
func (m *Meta) String(key string) (string, error) {
	if v, ok := m.Fields[key]; ok {
		return v, nil
	}
	return "", errors.Errorf("key %s not found", key)
}

// Float64 returns the float64 value of the field.
func (m *Meta) Float64(key string) (float64, error) {
	if v, ok := m.Fields[key]; ok {
		n, err := strconv.ParseFloat(v, 64)
		if err != nil {
			return 0, errors.Wrapf(err, "parsing key %s", key)
		}
		return n, nil
	}
	return 0, errors.Errorf("key %s not found", key)
}
