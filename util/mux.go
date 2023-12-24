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
	FileName     string     `json:"file"`
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
	mux.mux.HandleFunc("/info", mux.pkgInfo)
	mux.mux.HandleFunc("/get", mux.loadPkg)
	return mux
}

func (mux *RemoteMux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	mux.mux.ServeHTTP(w, r)
}

func (mux *RemoteMux) pkgList(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	name := "*"
	if query.Has("name") {
		name = query.Get("name")
	}

	pkgs, err := filepath.Glob(filepath.Join(mux.dir, fmt.Sprintf("%s.json", filepath.Base(name))))
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

func (mux *RemoteMux) pkgInfo(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	name := "*"
	if query.Has("name") {
		name = query.Get("name")
	}

	files, err := filepath.Glob(filepath.Join(mux.dir, fmt.Sprintf("%s.json", filepath.Base(name))))
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

func (mux *RemoteMux) loadPkg(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	if !query.Has("name") {
		w.WriteHeader(http.StatusBadRequest)

		fmt.Fprintln(w, "key \"name\" is required")
		return
	}

	pkgname := filepath.Base(query.Get("name"))

	file, err := os.Open(filepath.Join(mux.dir, pkgname+".json"))
	if os.IsNotExist(err) {
		w.WriteHeader(http.StatusBadRequest)

		fmt.Fprintf(w, "no package %q found\n", pkgname)
		return
	}

	if err != nil {
		mux.ErrorHandler(w, r, err)
		return
	}
	defer file.Close()

	var data Package
	err = json.NewDecoder(file).Decode(&data)

	if err != nil {
		mux.ErrorHandler(w, r, err)
		return
	}

	http.ServeFile(w, r, filepath.Join(mux.dir, "pkgs", data.FileName+".ipkg"))
}
