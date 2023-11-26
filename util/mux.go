package util

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

type Package struct {
	Name         string     `json:"name"`
	Version      string     `json:"version"`
	Dependencies []*Package `json:"deps"`
	DependentBy  []*Package `json:"needed"`
}

type result struct {
	Result any    `json:"result,omitempty"`
	Error  string `json:"error"`
}

type RemoteMux struct {
	mux          *http.ServeMux
	dir          string
	ErrorHandler func(w http.ResponseWriter, r *http.Request, err error)
}

func NewRemoteMux(dir string) *RemoteMux {
	mux := &RemoteMux{dir: dir, ErrorHandler: func(w http.ResponseWriter, r *http.Request, err error) {
		w.WriteHeader(500)
		fmt.Fprintf(w, "Error: %v", err)
	}, mux: http.NewServeMux()}
	mux.mux.HandleFunc("/list", mux.pkgList)
	mux.mux.HandleFunc("/get", mux.pkgSearch)
	return mux
}

func (mux *RemoteMux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	mux.mux.ServeHTTP(w, r)
}

func (mux *RemoteMux) pkgList(w http.ResponseWriter, r *http.Request) {
	pkgs, err := filepath.Glob(filepath.Join(mux.dir, "*.json"))
	if err != nil {
		mux.ErrorHandler(w, r, err)
		return
	}

	for i, pkg := range pkgs {
		pkgs[i] = filepath.Base(strings.TrimSuffix(pkg, ".json"))
	}
	res := result{pkgs, ""}
	json.NewEncoder(w).Encode(res)
}

func (mux *RemoteMux) pkgSearch(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	if !query.Has("name") {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(result{nil, "Required field name"})
		return
	}

	files, err := filepath.Glob(fmt.Sprintf("%s/%s.json", mux.dir, filepath.Base(query.Get("name"))))
	if err != nil {
		mux.ErrorHandler(w, r, err)
		return
	}
	var pkgs []Package
	for _, file := range files {
		var pkg Package

		data, err := os.ReadFile(file)
		if err != nil {
			mux.ErrorHandler(w, r, err)
			return
		}

		err = json.Unmarshal(data, &pkg)
		if err != nil {
			mux.ErrorHandler(w, r, err)
			return
		}

		pkgs = append(pkgs, pkg)
	}
	json.NewEncoder(w).Encode(result{pkgs, ""})
}
