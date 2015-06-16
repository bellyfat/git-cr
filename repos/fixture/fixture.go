package fixture

import (
	"bytes"
	"encoding/base64"
	"io"
	"io/ioutil"

	"github.com/lucas-clemente/git-cr/git"
)

// A FixtureRepo for tests
type FixtureRepo struct {
	CurrentRefs     git.Refs
	PackfilesFromTo map[string]map[string][]byte
}

var _ git.ListingRepo = &FixtureRepo{}

// NewFixtureRepo makes a new fixture repo
func NewFixtureRepo() *FixtureRepo {
	return &FixtureRepo{
		CurrentRefs:     git.Refs{},
		PackfilesFromTo: map[string]map[string][]byte{"": map[string][]byte{}},
	}
}

// AddPackfile adds a base64-encoded packfile to the repo
func (b *FixtureRepo) AddPackfile(from, to, b64 string) {
	m, ok := b.PackfilesFromTo[from]
	if !ok {
		b.PackfilesFromTo[from] = map[string][]byte{}
		m = b.PackfilesFromTo[from]
	}
	pack, err := base64.StdEncoding.DecodeString(b64)
	if err != nil {
		panic("invalid base64 in FixtureRepo.AddPackfile")
	}
	m[to] = pack
}

// FindDelta implements git.Repo
func (b *FixtureRepo) FindDelta(from, to string) (git.Delta, error) {
	m, ok := b.PackfilesFromTo[from]
	if !ok {
		return nil, git.ErrorDeltaNotFound
	}
	p, ok := m[to]
	if !ok {
		return nil, git.ErrorDeltaNotFound
	}
	return ioutil.NopCloser(bytes.NewBuffer(p)), nil
}

// GetRefs implements git.Repo
func (b *FixtureRepo) GetRefs() (git.Refs, error) {
	return b.CurrentRefs, nil
}

// ReadPackfile implements git.Repo
func (*FixtureRepo) ReadPackfile(d git.Delta) (io.ReadCloser, error) {
	return d.(io.ReadCloser), nil
}

// UpdateRef implements git.Repo
func (b *FixtureRepo) UpdateRef(update git.RefUpdate) error {
	if update.NewID == "" {
		delete(b.CurrentRefs, update.Name)
	} else {
		b.CurrentRefs[update.Name] = update.NewID
	}
	return nil
}

// WritePackfile implements git.Repo
func (b *FixtureRepo) WritePackfile(from, to string, r io.Reader) error {
	data, err := ioutil.ReadAll(r)
	if err != nil {
		return err
	}
	m, ok := b.PackfilesFromTo[from]
	if !ok {
		b.PackfilesFromTo[from] = map[string][]byte{}
		m = b.PackfilesFromTo[from]
	}
	m[to] = data
	return nil
}

// ListAncestors implements git.Repo
func (b *FixtureRepo) ListAncestors(target string) ([]string, error) {
	var results []string
	for from, toMap := range b.PackfilesFromTo {
		for to := range toMap {
			if to == target {
				results = append(results, from)
			}
		}
	}
	return results, nil
}