package services

import (
	"encoding/json"
	"log"
	"sort"

	"analyzer/pkg/app/backends"
)

type Service struct {
	name string // service name

	impl        string // impl name
	pkg         string // simple package name
	pkgpath     string
	path        string // format: <pkgpath>.<name>
	constructor string // method name

	deps []*Service
	dbs  []*backends.Database

	fields []*Field

	methods     []string
	initializer bool     // true if it has Run method
	wiringNames []string // IDs for arguments passed in blueprint wiring
}

func NewService(name string, impl string, pkg string, pkgpath string, path string, constructor string, constructorArgs []string) *Service {
	return &Service{
		name:        name,
		impl:        impl,
		pkg:         pkg,
		path:        path,
		pkgpath:     pkgpath,
		constructor: constructor,
		wiringNames: constructorArgs,
	}
}

func (service *Service) GetAllWiringNames() []string {
	return service.wiringNames
}

func (service *Service) GetWiringNameAt(idx int) string {
	if idx < len(service.wiringNames) {
		return service.wiringNames[idx]
	}
	log.Panicf("[SERVICE] index (%d) out of bounds for constructor args: %v", idx, service.wiringNames)
	return ""
}

func (service *Service) AddField(field *Field) {
	service.fields = append(service.fields, field)
}

func (service *Service) GetFieldAt(idx int) *Field {
	if idx < len(service.fields) {
		return service.fields[idx]
	}
	log.Panicf("[SERVICE] index (%d) out of bounds for service fields: %v", idx, service.fields)
	return nil
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

func (service *Service) GetConstructor() string {
	return service.constructor
}

func (service *Service) GetPackage() string {
	return service.pkg
}

func (service *Service) GetPath() string {
	return service.path
}

func (service *Service) GetPackagePath() string {
	return service.pkgpath
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

	// do not sort these because they are already sorted by idx in the service struct
	fieldsStrLst := make([]string, len(service.fields))
	for idx, field := range service.fields {
		fieldsStrLst[idx] = field.String()
	}

	return json.Marshal(&struct {
		Name      string   `json:"name"`
		Path      string   `json:"path"`
		Pkg       string   `json:"pkg"`
		Impl      string   `json:"impl"`
		Fields    []string `json:"fields"`
		Methods   []string `json:"methods"`
		Services  []string `json:"services"`
		Databases []string `json:"databases"`
	}{
		Name:      service.name,
		Path:      service.path,
		Pkg:       service.pkg,
		Impl:      service.impl,
		Fields:    fieldsStrLst,
		Methods:   service.methods,
		Services:  depNames,
		Databases: dbNames,
	})
}
