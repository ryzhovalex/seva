package sevent

import (
	"encoding/json"
	"io"
	"os"
	"path"
	"seva/internal/domains"
	"seva/lib/rpc"
	"seva/lib/utils"
	"sort"
	"strings"

	"github.com/gin-gonic/gin"
)

type StateEvent struct {
	Id   string
	Type string
	Time utils.Time
	Body map[string]any
}

// Fields by their names.
type Spec map[string]Field

type FieldType string

const (
	String  FieldType = "string"
	Number  FieldType = "number"
	Boolean FieldType = "boolean"
	Array   FieldType = "array"
	Object  FieldType = "object"
)

type Field struct {
	Type FieldType
	// Array of contained specs.
	// For array: their name is not considered.
	// For array: only first field is considered, since we allow only for
	// 			  same-type arrays.
	// For object: each subspec name is a key in the nested object.
	Fields []Field
}

type TypeToSpec map[string]Spec

func GetEventTypes(domain string) ([]string, *utils.Error) {
	e := domains.CheckDomainCreated(domain)
	if e != nil {
		return nil, e
	}

	dir := domains.GetDomainDir(domain)

	// Event types can be effectively fetched from Var/Domains/DOMAIN/Specs/...
	files, be := os.ReadDir(dir + "/Specs")
	if be != nil {
		return nil, utils.BE(be)
	}

	r := []string{}
	for _, file := range files {
		if strings.HasSuffix(file.Name(), ".json") {
			r = append(r, strings.Replace(file.Name(), ".json", "", 1))
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
	types, e := GetEventTypes(domain)
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

	path := domains.GetDomainDir(domain) + "/Specs/" + eventType + ".json"
	f, be := os.Open(path)
	if be != nil {
		return nil, utils.BE(be)
	}
	defer f.Close()

	data := Spec{}
	b, be := io.ReadAll(f)
	if be != nil {
		return nil, utils.BE(be)
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
		return nil, utils.DE("Event type spec already exists: " + eventType)
	}

	return &spec, nil
}

func CheckSpecCreated(domain string, eventType string) *utils.Error {
	e := domains.CheckDomainCreated(domain)
	if e != nil {
		return e
	}
	specPath := domains.GetDomainDir(domain) + "/Specs/" + eventType + ".json"
	_, be := os.Stat(specPath)
	if be != nil {
		return utils.DE("Spec is not created for event type: " + eventType)
	}
	return nil
}

func CheckSpecNotCreated(domain string, eventType string) *utils.Error {
	e := CheckSpecCreated(domain, eventType)
	if e == nil {
		return utils.DE("Spec is already created for event type: " + eventType)
	}
	return nil
}

func CreateEvent(domain string, eventType string, body map[string]any) (*StateEvent, *utils.Error) {
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
	domainDir := domains.GetDomainDir(domain)
	f, be := os.Create(domainDir + "/" + id + ".json")
	if be != nil {
		return nil, utils.BE(be)
	}
	defer f.Close()

	jsonBytes, be := json.Marshal(event)
	if be != nil {
		return nil, utils.BE(be)
	}

	_, be = f.Write(jsonBytes)
	if be != nil {
		return nil, utils.BE(be)
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
		return nil, utils.BE(be)
	}

	events := []StateEvent{}
	for _, fileEntry := range files {
		filePath := path.Join(dir, fileEntry.Name())
		f, be := os.Open(filePath)
		if be != nil {
			return nil, utils.BE(be)
		}

		jsonBytes, be := io.ReadAll(f)
		if be != nil {
			return nil, utils.BE(be)
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

type CreateEventData struct {
	Domain    string
	EventType string
	Body      map[string]any
}

func RpcCreateEvent(c *gin.Context) {
	var form CreateEventData
	be := c.Bind(&form)
	if be != nil {
		rpc.Error(c, utils.BE(be))
		return
	}

	// var body map[string]any
	// be = json.Unmarshal([]byte(form.Body), &body)
	// if be != nil {
	// 	rpc.Error(c, utils.BE(be))
	// 	return
	// }

	event, e := CreateEvent(
		form.Domain, form.EventType, form.Body,
	)
	if e != nil {
		rpc.Error(c, e)
		return
	}

	rpc.Ok(c, event)
}

type DomainData struct {
	Domain string
}

func RpcGetEvents(c *gin.Context) {
	var data DomainData
	be := c.Bind(&data)
	if be != nil {
		rpc.Error(c, utils.BE(be))
		return
	}

	r, e := GetEvents(data.Domain)
	if e != nil {
		rpc.Error(c, e)
		return
	}
	rpc.Ok(c, r)
}

func RpcGetSpecs(c *gin.Context) {
	var data DomainData
	be := c.Bind(&data)
	if be != nil {
		rpc.Error(c, utils.BE(be))
		return
	}

	r, e := GetSpecs(data.Domain)
	if e != nil {
		rpc.Error(c, e)
		return
	}
	rpc.Ok(c, r)
}
