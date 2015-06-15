package git_test

import (
	"bytes"
	"encoding/base64"
	"io"
	"io/ioutil"

	"github.com/lucas-clemente/git-cr/git"
)

// A fixtureBackend for tests
type fixtureBackend struct {
	currentRefs     git.Refs
	packfilesFromTo map[string]map[string][]byte
}

var _ git.Backend = &fixtureBackend{}

func newFixtureBackend() *fixtureBackend {
	return &fixtureBackend{
		packfilesFromTo: map[string]map[string][]byte{"": map[string][]byte{}},
	}
}

func (b *fixtureBackend) addPackfile(from, to, b64 string) {
	m, ok := b.packfilesFromTo[from]
	if !ok {
		b.packfilesFromTo[from] = map[string][]byte{}
		m = b.packfilesFromTo[from]
	}
	pack, err := base64.StdEncoding.DecodeString(b64)
	if err != nil {
		panic("invalid base64 in fixtureBackend.addPackfile")
	}
	m[to] = pack
}

func (b *fixtureBackend) FindDelta(from, to string) (git.Delta, error) {
	m, ok := b.packfilesFromTo[from]
	if !ok {
		return nil, git.ErrorDeltaNotFound
	}
	p, ok := m[to]
	if !ok {
		return nil, git.ErrorDeltaNotFound
	}
	return ioutil.NopCloser(bytes.NewBuffer(p)), nil
}

func (b *fixtureBackend) GetRefs() (git.Refs, error) {
	return b.currentRefs, nil
}

func (*fixtureBackend) ReadPackfile(d git.Delta) (io.ReadCloser, error) {
	return d.(io.ReadCloser), nil
}

func (b *fixtureBackend) UpdateRef(update git.RefUpdate) error {
	if update.NewID == "" {
		delete(b.currentRefs, update.Name)
	} else {
		b.currentRefs[update.Name] = update.NewID
	}
	return nil
}

func (b *fixtureBackend) WritePackfile(from, to string, r io.Reader) error {
	data, err := ioutil.ReadAll(r)
	if err != nil {
		return err
	}
	m, ok := b.packfilesFromTo[from]
	if !ok {
		b.packfilesFromTo[from] = map[string][]byte{}
		m = b.packfilesFromTo[from]
	}
	m[to] = data
	return nil
}

func (b *fixtureBackend) ListAncestors(target string) ([]string, error) {
	var results []string
	for from, toMap := range b.packfilesFromTo {
		for to := range toMap {
			if to == target {
				results = append(results, from)
			}
		}
	}
	return results, nil
}
