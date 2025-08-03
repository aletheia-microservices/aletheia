package services

import (
	"encoding/json"
	"sort"

	"analyzer/pkg/app/backends"
)

type Service struct {
	name string // service name

	impl        string // impl name
	pkg         string // simple package name
	path        string // format: <pkgpath>.<name>
	constructor string // method name

	deps []*Service
	dbs  []*backends.Database

	// TODO: create a struct with more info for methods
	methods     []string
	initializer bool // true if it has Run method
}

func NewService(name string, impl string, pkg string, path string, constructor string) *Service {
	return &Service{
		name:        name,
		impl:        impl,
		pkg:         pkg,
		path:        path,
		constructor: constructor,
	}
}

func (service *Service) GetMethods() []string {
	return service.methods
}

func (service *Service) SetMethods(methods ...string) {
	for _, method := range methods {
		if method == "Run" {
			service.initializer = true
		}
		service.methods = append(service.methods, method)
	}
}

func (service *Service) HasInitializerMethod() bool {
	return service.initializer
}

func (service *Service) HasMethod(name string) bool {
	for _, method := range service.methods {
		if method == name {
			return true
		}
	}
	return false
}

func (service *Service) AddDependency(dep *Service) {
	service.deps = append(service.deps, dep)
}

func (service *Service) AddDatabase(db *backends.Database) {
	service.dbs = append(service.dbs, db)
}

func (service *Service) GetName() string {
	return service.name
}

func (service *Service) GetImpl() string {
	return service.impl
}

func (service *Service) GetPackage() string {
	return service.pkg
}

func (service *Service) GetPath() string {
	return service.path
}

func (service *Service) String() string {
	var str string
	str += "\n\tpath: " + service.path
	str += "\n\tpkg: " + service.pkg
	str += "\n\timpl: " + service.impl
	str += "\n\tservices: {"
	for i, dep := range service.deps {
		str += dep.GetName()
		if i < len(service.deps)-1 {
			str += ", "
		}
	}
	str += "}\n\tdatabases: {"
	for i, dep := range service.dbs {
		str += dep.GetName()
		if i < len(service.dbs)-1 {
			str += ", "
		}
	}
	str += "}"
	return service.name + ": " + str
}

func (service *Service) MarshalJSON() ([]byte, error) {
	depNames := make([]string, len(service.deps))
	for i, dep := range service.deps {
		depNames[i] = dep.GetName()
	}

	dbNames := make([]string, len(service.dbs))
	for i, db := range service.dbs {
		dbNames[i] = db.GetName()
	}

	sort.Strings(depNames)
	sort.Strings(dbNames)

	return json.Marshal(&struct {
		Name      string   `json:"name"`
		Path      string   `json:"path"`
		Pkg       string   `json:"pkg"`
		Impl      string   `json:"impl"`
		Methods   []string `json:"methods"`
		Services  []string `json:"services"`
		Databases []string `json:"databases"`
	}{
		Name:      service.name,
		Path:      service.path,
		Pkg:       service.pkg,
		Impl:      service.impl,
		Methods:   service.methods,
		Services:  depNames,
		Databases: dbNames,
	})
}
