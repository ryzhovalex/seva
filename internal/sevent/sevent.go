package sevent

import (
	"encoding/json"
	"io"
	"os"
	"path"
	"seva/internal/domains"
	"seva/lib/utils"
	"sort"
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

func GetTypes(domain string) ([]string, *utils.Error) {
	e := domains.CheckDomainCreated(domain)
	if e != nil {
		return nil, e
	}

	dir := domains.GetDomainDir(domain)

	// Read the directory
	files, be := os.ReadDir(dir)
	if be != nil {
		return nil, utils.CreateDefaultErrorFromBase(be)
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
	e := domains.CheckDomainCreated(domain)
	if e != nil {
		return nil, e
	}

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
	e := domains.CheckDomainCreated(domain)
	if e != nil {
		return nil, e
	}

	path := domains.GetDomainDir(domain) + "/" + eventType + "/SPEC.json"
	f, be := os.Open(path)
	if be != nil {
		return nil, utils.CreateDefaultErrorFromBase(be)
	}
	defer f.Close()

	data := Spec{}
	b, be := io.ReadAll(f)
	if be != nil {
		return nil, utils.CreateDefaultErrorFromBase(be)
	}
	json.Unmarshal(b, &data)

	return data, nil
}

func CreateSpec(domain string, eventType string, spec Spec) (*Spec, *utils.Error) {
	e := CheckSpecNotCreated(domain, eventType)
	if e != nil {
		return nil, e
	}

	_, e = GetSpec(domain, eventType)
	if e == nil {
		return nil, utils.CreateDefaultError("Event type spec already exists: " + eventType)
	}

	return &spec, nil
}

func CheckSpecCreated(domain string, eventType string) *utils.Error {
	e := domains.CheckDomainCreated(domain)
	if e != nil {
		return e
	}
	specPath := domains.GetDomainDir(domain) + "/" + eventType + "/SPEC.json"
	_, be := os.Stat(specPath)
	if be != nil {
		return utils.CreateDefaultError("Spec is not created for event type: " + eventType)
	}
	return nil
}

func CheckSpecNotCreated(domain string, eventType string) *utils.Error {
	e := CheckSpecCreated(domain, eventType)
	if e == nil {
		return utils.CreateDefaultError("Spec is already created for event type: " + eventType)
	}
	return nil
}

func CreateEvent(domain string, eventType string, body any) (*StateEvent, *utils.Error) {
	e := CheckSpecCreated(domain, eventType)
	if e != nil {
		return nil, e
	}

	// For now don't check spec validity, just insert body as it is.
	id := utils.MakeUuid()
	event := StateEvent{
		Id:   id,
		Type: eventType,
		Time: utils.TimeNow(),
		Body: body,
	}
	eventDir := domains.GetDomainDir(domain) + "/" + eventType
	f, be := os.Create(eventDir + "/" + id + ".json")
	if be != nil {
		return nil, utils.CreateDefaultErrorFromBase(be)
	}
	defer f.Close()

	jsonBytes, be := json.Marshal(event)
	if be != nil {
		return nil, utils.CreateDefaultErrorFromBase(be)
	}

	_, be = f.Write(jsonBytes)
	if be != nil {
		return nil, utils.CreateDefaultErrorFromBase(be)
	}

	return &event, nil
}

func GetEvents(domain string) ([]StateEvent, *utils.Error) {
	e := domains.CheckDomainCreated(domain)
	if e != nil {
		return nil, e
	}
	dir := domains.GetDomainDir(domain)

	files, be := os.ReadDir(dir)
	if be != nil {
		return nil, utils.CreateDefaultErrorFromBase(be)
	}

	events := []StateEvent{}
	for _, fileEntry := range files {
		filePath := path.Join(dir, fileEntry.Name())
		f, be := os.Open(filePath)
		if be != nil {
			return nil, utils.CreateDefaultErrorFromBase(be)
		}

		jsonBytes, be := io.ReadAll(f)
		if be != nil {
			return nil, utils.CreateDefaultErrorFromBase(be)
		}
		event := StateEvent{}
		json.Unmarshal(jsonBytes, &event)
		events = append(events, event)
	}

	sort.Slice(events, func(i, j int) bool {
		return events[i].Time < events[j].Time
	})

	return events, nil
}
