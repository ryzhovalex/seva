package sevent

import (
	"encoding/json"
	"io"
	"os"
	"seva/utils"
)

type StateEvent struct {
	Id   string
	Type string
	Time utils.Time
	Body any
}

type Spec []Field

type FieldType int

const (
	String FieldType = iota
	Number
	Boolean
	Array
	Object
	Null
)

type Field struct {
	Name string
	// Array of contained specs. For Type=Array their name is not considered.
	// For Type=Object each subspec name is a key in the nested object.
	ContainedSpecs []Field
}

type TypeToSpec map[string]Spec

func GetDomainDir(domain string) string {
	return "Var/Domains/" + domain
}

func GetTypes(domain string) ([]string, *utils.Error) {
	dir := GetDomainDir(domain)

	// Read the directory
	files, be := os.ReadDir(dir)
	if be != nil {
		return nil, utils.NewDefaultErrorFromBase(be)
	}

	r := []string{}
	// Iterate over the files and directories
	for _, file := range files {
		// Check if the file is a directory
		if file.IsDir() {
			r = append(r, file.Name())
		}
	}
	return r, nil
}

func GetSpecs(domain string) (TypeToSpec, *utils.Error) {
	specs := TypeToSpec{}
	types, e := GetTypes(domain)
	if e != nil {
		return nil, e
	}
	for _, t := range types {
		if _, ok := specs[t]; ok {
			utils.Log("[Warning] Duplicate specification for type: " + t)
			continue
		}
		spec, e := GetSpec(domain, t)
		if e != nil {
			return nil, e
		}
		specs[t] = spec
	}

	return specs, nil
}

func GetSpec(domain string, eventType string) (Spec, *utils.Error) {
	path := GetDomainDir(domain) + "/" + eventType + ".json"
	f, be := os.Open(path)
	if be != nil {
		return nil, utils.NewDefaultErrorFromBase(be)
	}
	defer f.Close()

	data := Spec{}
	b, be := io.ReadAll(f)
	if be != nil {
		return nil, utils.NewDefaultErrorFromBase(be)
	}
	json.Unmarshal(b, &data)

	return data, nil
}

func NewSpec(domain string, eventType string, spec Spec) *utils.Error {
	_, e := GetSpec(domain, eventType)
	if e == nil {
		return utils.NewDefaultError("Event type already exists: " + eventType)
	}
	return nil
}
