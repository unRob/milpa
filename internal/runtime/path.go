// Copyright © 2021 Roberto Hidalgo <milpa@un.rob.mx>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package runtime

import (
	"os"
	"path/filepath"
	"sort"
	"sync"

	"github.com/sirupsen/logrus"
)

func IsDir(path string, warn bool) bool {
	if fi, err := os.Stat(path); err == nil {
		if fi.Mode().IsDir() {
			return true
		}
	}

	if warn {
		logrus.Warnf("Discarding non-directory <%s> from MILPA_PATH", path)
	}

	return false
}

type pathLayer map[string]bool

func (pl pathLayer) add(path string) {
	if _, inMap := pl[path]; !inMap {
		pl[path] = true
	}
}

type lookupFunc func() []string

type PathBuilder struct {
	layers   map[int]*pathLayer
	unique   map[string]bool
	lookups  []lookupFunc
	resolved bool
	mutex    sync.Mutex
}

func NewPathBuilder() *PathBuilder {
	return &PathBuilder{
		layers:  map[int]*pathLayer{},
		unique:  map[string]bool{},
		lookups: []lookupFunc{},
	}
}

func (pb *PathBuilder) LookupLen() int {
	return len(pb.lookups)
}

// AddLookup adds a lookup function if envVar is unset or falseish.
func (pb *PathBuilder) AddLookup(envVar string, fn lookupFunc) {
	if !isTrueIsh(os.Getenv(envVar)) {
		pb.lookups = append(pb.lookups, fn)
	}
}

func (pb *PathBuilder) resolve() {
	if pb.resolved {
		return
	}
	var wg sync.WaitGroup
	for idx, lookup := range pb.lookups {
		wg.Add(1)
		lookup := lookup
		layerID := idx + 10
		go func() {
			defer wg.Done()
			for _, f := range lookup() {
				pb.Add(layerID, f)
			}
		}()
	}

	wg.Wait()
	pb.resolved = true
}

// Add appends unique symlink-resolved paths at the given layer.
func (pb *PathBuilder) Add(layerID int, path string) {
	pb.mutex.Lock()
	defer pb.mutex.Unlock()
	if pb.unique == nil {
		pb.unique = map[string]bool{}
	}

	if pb.layers == nil {
		pb.layers = map[int]*pathLayer{}
	}

	// check for uniqueness on unresolved path, as they may not be symlinks
	if _, exists := pb.unique[path]; exists {
		return
	}

	// Resolve symlinks before checking if unique
	if pathR, err := os.Readlink(path); err == nil {
		// Output of os.Readlink is OS-dependent...
		if !filepath.IsAbs(pathR) {
			pathR = filepath.Join(filepath.Dir(path), pathR)
		}
		path = pathR
	}

	if _, exists := pb.unique[path]; exists {
		return
	}

	pb.unique[path] = true

	if _, exists := pb.layers[layerID]; !exists {
		pb.layers[layerID] = &pathLayer{}
	}

	pb.layers[layerID].add(path)
}

func (pb *PathBuilder) Ordered() []string {
	pb.resolve()
	res := []string{}
	keys := []int{}
	for key := range pb.layers {
		keys = append(keys, key)
	}
	sort.Ints(keys)

	for _, key := range keys {
		layer := pb.layers[key]
		group := []string{}
		for path := range *layer {
			group = append(group, path)
		}
		sort.Strings(group)
		res = append(res, group...)
	}

	return res
}