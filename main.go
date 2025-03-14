package main

import (
	"encoding/json"
	"flag"
	"os"
	"path/filepath"
	"seva/lib/bone"
	"seva/lib/shell"
	"strconv"
	"strings"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

const (
	OK = iota
	ERROR
)

type Event_Signature struct {
	// Integer type is index of signature in the state array.
	Type_Name string `json:"type_name"`
	// Values can be:
	//   - int
	//   - str
	//   - float
	//   - array
	//   - dict
	//   - bool
	Fields map[string]string `json:"fields"`
}

type Event struct {
	// Time of event injection.
	Created_Sec int `json:"created_sec"`
	// Integer type of an event. Each project has own unsigned set of types,
	// starting from 1.
	Type   int               `json:"type"`
	Fields map[string]string `json:"fields"`
}

// Domains by their list of events
var events = map[string][]*Event{}

// Domains by their list of event signatures
var signatures = map[string][]*Event_Signature{}

// Event files by their domains
var eventfiles = map[string]*os.File{}

// Signature files by their domains
var sigfiles = map[string]*os.File{}

func read_event_state() int {
	dir := bone.Userdir("events")
	bone.Mkdir(dir)
	files, er := os.ReadDir(dir)
	if er != nil {
		bone.Log_Error("During state reading, cannot read userdir, error: %s", er)
		return ERROR
	}
	for _, file := range files {
		if filepath.Ext(file.Name()) == ".json" {
			path := filepath.Join(dir, file.Name())
			data, er := os.ReadFile(path)
			if er != nil {
				bone.Log_Error("During state reading, cannot read file '%s', error: %s", path, er)
				return ERROR
			}
			domain, _ := strings.CutSuffix(file.Name(), filepath.Ext(file.Name()))
			evs := []*Event{}
			er = json.Unmarshal(data, &evs)
			events[domain] = evs
			if er != nil {
				bone.Log_Error("Cannot unmarshal file '%s', error: %s", file.Name(), er)
				return ERROR
			}
		}
	}
	return OK
}

func read_signature_state() int {
	dir := bone.Userdir("signatures")
	bone.Mkdir(dir)
	files, er := os.ReadDir(dir)
	if er != nil {
		bone.Log_Error("During state reading, cannot read userdir, error: %s", er)
		return ERROR
	}
	for _, file := range files {
		if filepath.Ext(file.Name()) == ".json" {
			path := filepath.Join(dir, file.Name())
			data, er := os.ReadFile(path)
			if er != nil {
				bone.Log_Error("During state reading, cannot read file '%s', error: %s", path, er)
				return ERROR
			}
			domain, _ := strings.CutSuffix(file.Name(), filepath.Ext(file.Name()))
			sigs := []*Event_Signature{}
			er = json.Unmarshal(data, &sigs)
			signatures[domain] = sigs
			if er != nil {
				bone.Log_Error("Cannot unmarshal file '%s', error: %s", file.Name(), er)
				return ERROR
			}
		}
	}
	return OK
}

// Read all files in userdir and unmarshal them to state.
func read_state() int {
	e := read_signature_state()
	if e != OK {
		return e
	}
	e = read_event_state()
	if e != OK {
		return e
	}

	// Add "main" domain if does not exist
	_, ok := signatures["main"]
	if !ok {
		signatures["main"] = []*Event_Signature{}
		events["main"] = []*Event{}
		save_state()
	}

	return OK
}

func save_state() {
	eventdir := bone.Userdir("events")
	bone.Mkdir(eventdir)
	for domain, evs := range events {
		data, er := json.MarshalIndent(evs, "", "\t")
		if er != nil {
			bone.Log_Error("Error marshalling state to json for domain '%s'", domain)
			return
		}

		f, ok := eventfiles[domain]
		if !ok {
			path := filepath.Join(eventdir, domain+".json")
			f, er = os.OpenFile(path, os.O_WRONLY|os.O_CREATE, 0644)
			if er != nil {
				bone.Log_Error("Cannot open data file for domain '%s' at path '%s'", domain, path)
				continue
			}
			eventfiles[domain] = f
		}

		er = f.Truncate(0)
		if er != nil {
			bone.Log_Error("Cannot truncate file for domain '%s'", domain)
			continue
		}

		_, er = f.Seek(0, 0)
		if er != nil {
			bone.Log_Error("Cannot seek file of domain '%s'", domain)
			continue
		}

		_, er = f.Write(data)
		if er != nil {
			bone.Log_Error("Cannot write to a file for domain '%s'", domain)
			continue
		}
	}

	sigdir := bone.Userdir("signatures")
	bone.Mkdir(sigdir)
	for domain, sigs := range signatures {
		data, er := json.MarshalIndent(sigs, "", "\t")
		if er != nil {
			bone.Log_Error("Error marshalling state to json for domain '%s'", domain)
			return
		}

		f, ok := sigfiles[domain]
		if !ok {
			path := filepath.Join(sigdir, domain+".json")
			f, er = os.OpenFile(path, os.O_WRONLY|os.O_CREATE, 0644)
			if er != nil {
				bone.Log_Error("Cannot open data file for domain '%s' at path '%s'", domain, path)
				continue
			}
			sigfiles[domain] = f
		}

		er = f.Truncate(0)
		if er != nil {
			bone.Log_Error("Cannot truncate file for domain '%s'", domain)
			continue
		}

		_, er = f.Seek(0, 0)
		if er != nil {
			bone.Log_Error("Cannot seek file of domain '%s'", domain)
			continue
		}

		_, er = f.Write(data)
		if er != nil {
			bone.Log_Error("Cannot write to a file for domain '%s'", domain)
			continue
		}
	}
}

func create_server() *gin.Engine {
	server := gin.New()
	server.Use(gin.Recovery())
	server.Use(cors.Default())

	return server
}

func deinit() {
	for _, f := range eventfiles {
		f.Close()
	}
	for k := range eventfiles {
		delete(eventfiles, k)
	}
}

func main() {
	defer deinit()

	shell_enabled := flag.Bool("shell", false, "Enables shell mode.")
	bone.Init("seva")

	e := read_state()
	if e != OK {
		bone.Log_Error("During state reading, an error occurred: %d", e)
		os.Exit(e)
		return
	}

	if *shell_enabled {
		shell.Init()

		domain := bone.Config.Get_String("main", "domain", "main")
		shell.Set_Domain(domain)

		shell.Set_Command("setdomain", shell_set_domain)
		shell.Set_Command("addevent", shell_add_event)
		shell.Set_Command("addsig", shell_add_signature)

		save_state()
		shell.Run()
		return
	}

	server := create_server()
	server.Run("0.0.0.0:3000")
}

func shell_add_signature(c *shell.Command_Context) int {
	buffer := c.Arg_String("_", "")
	if buffer == "" {
		bone.Log_Error("Specify at least event type")
		return shell.ERROR
	}
	parts := strings.Split(buffer, " ")
	str_type := strings.ToUpper(parts[0])
	domain := shell.Get_Domain()

	for _, signature := range signatures[domain] {
		if signature.Type_Name == str_type {
			bone.Log_Error("Signature '%s' already exist", str_type)
			return shell.ERROR
		}
	}

	fields := map[string]string{}

	for i, part := range parts {
		if i == 0 {
			continue
		}
		subparts := strings.Split(part, "=")
		if len(subparts) != 2 {
			bone.Log_Error("Invalid part '%s'", part)
			return shell.ERROR
		}
		key := subparts[0]
		value := subparts[1]

		// We store string anyways, but check signature
		switch value {
		case "int":
		case "str":
		case "float":
		case "bool":
		case "arr":
		case "dict":
		default:
			bone.Log_Error("Unrecognized signature value '%s' for event '%s'", value, str_type)
			return shell.ERROR
		}
		fields[key] = value
	}

	signature := &Event_Signature{
		Type_Name: str_type,
		Fields:    fields,
	}
	_, ok := signatures[domain]
	if !ok {
		signatures[domain] = []*Event_Signature{}
	}
	signatures[domain] = append(signatures[domain], signature)

	save_state()
	return shell.OK
}

func shell_add_event(c *shell.Command_Context) int {
	buffer := c.Arg_String("_", "")
	if buffer == "" {
		bone.Log_Error("Specify at least event type")
		return shell.ERROR
	}
	parts := strings.Split(buffer, " ")
	str_type := strings.ToUpper(parts[0])

	domain := shell.Get_Domain()
	sigs, ok := signatures[domain]
	if !ok {
		bone.Log_Error("Cannot find domain '%s'", domain)
		return shell.ERROR
	}
	var target_signature_type int
	var target_signature *Event_Signature = nil
	for i, signature := range sigs {
		if signature.Type_Name == str_type {
			target_signature_type = i + 1
			target_signature = signature
		}
	}
	if target_signature == nil {
		bone.Log_Error("Cannot find signature for type '%s'", str_type)
		return shell.ERROR
	}

	// Parse event fields and compare with signature
	fields := map[string]string{}
	for i, part := range parts {
		if i == 0 {
			continue
		}
		subparts := strings.Split(part, "=")
		if len(subparts) != 2 {
			bone.Log_Error("Invalid part '%s'", part)
			return shell.ERROR
		}
		key := subparts[0]
		value := subparts[1]
		sig_value, ok := target_signature.Fields[key]
		if !ok {
			bone.Log_Error("No field with key '%s' in signature for event '%s'", key, str_type)
			return shell.ERROR
		}

		// We store string anyways, but check signature
		switch sig_value {
		case "int":
			_, er := strconv.Atoi(value)
			if er != nil {
				bone.Log_Error("Cannot convert value '%s' to int for event of type '%s'", sig_value, str_type)
				return shell.ERROR
			}
		case "str":
		case "float":
			_, er := strconv.ParseFloat(value, 64)
			if er != nil {
				bone.Log_Error("Cannot convert value '%s' to float for event of type '%s'", sig_value, str_type)
				return shell.ERROR
			}
		case "bool":
			if value != "1" && value != "0" && value != "true" && value != "false" {
				bone.Log_Error("Cannot convert value '%s' to bool for event of type '%s'", sig_value, str_type)
				return shell.ERROR
			}
		// @Todo implement parsers for arr and dict
		case "arr":
		case "dict":
		default:
			bone.Log_Error("Unrecognized value '%s' for signature of event '%s'", sig_value, str_type)
			return shell.ERROR
		}
		fields[key] = value
	}

	event := &Event{
		Created_Sec: int(bone.Utc()),
		Type:        target_signature_type,
		Fields:      fields,
	}

	evs, ok := events[domain]
	if !ok {
		bone.Log_Error("Cannot find domain '%s'", domain)
		return shell.ERROR
	}
	events[domain] = append(evs, event)

	save_state()
	return shell.OK
}

func shell_set_domain(c *shell.Command_Context) int {
	domain := c.Arg_String("_", "main")

	if shell.Get_Domain() == domain {
		return shell.OK
	}

	e := shell.Set_Domain(domain)
	if e != OK {
		return shell.ERROR
	}

	// Cache domain in config for future logins
	bone.Config.Write_String("main", "domain", domain)

	_, ok := events[domain]
	if !ok {
		events[domain] = []*Event{}
	}
	_, ok = signatures[domain]
	if !ok {
		signatures[domain] = []*Event_Signature{}
	}
	save_state()

	return shell.OK
}
