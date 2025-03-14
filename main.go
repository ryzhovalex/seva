package main

import (
	"encoding/json"
	"flag"
	"os"
	"path/filepath"
	"seva/lib/bone"
	"seva/lib/shell"
	"strings"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

const (
	OK = iota
	ERROR
)

type Event struct {
	// Time of event injection.
	Created_Sec int `json:"created_sec"`
	// Integer type of an event. Each project has own unsigned set of types,
	// starting from 1.
	Type int `json:"type"`
}

// Domains by their list of events
var state = map[string][]*Event{}

// Datafiles by their domains
var datafiles = map[string]*os.File{}

// Read all files in userdir and unmarshal them to state.
func read_state() int {
	files, er := os.ReadDir(bone.Userdir())
	if er != nil {
		bone.Log_Error("During state reading, cannot read userdir, error: %s", er)
		return ERROR
	}
	for _, file := range files {
		if filepath.Ext(file.Name()) == ".json" {
			path := filepath.Join(bone.Userdir(), file.Name())
			data, er := os.ReadFile(path)
			if er != nil {
				bone.Log_Error("During state reading, cannot read file '%s', error: %s", file.Name(), er)
				return ERROR
			}
			domain, _ := strings.CutSuffix(file.Name(), filepath.Ext(file.Name()))
			events := []*Event{}
			state[domain] = events
			er = json.Unmarshal(data, &events)
			if er != nil {
				bone.Log_Error("Cannot unmarshal file '%s', error: %s", file.Name(), er)
				return ERROR
			}
			bone.Log("Loaded domain '%s'", domain)
		}
	}

	// Add "main" domain if does not exist
	_, ok := state["main"]
	if !ok {
		state["main"] = []*Event{}
		save_state()
	}

	return OK
}

func save_state() {
	datadir := bone.Userdir("data")
	bone.Mkdir(datadir)
	for domain, events := range state {
		data, er := json.MarshalIndent(events, "", "\t")
		if er != nil {
			bone.Log_Error("Error marshalling state to json for domain '%s'", domain)
			return
		}

		f, ok := datafiles[domain]
		if !ok {
			path := filepath.Join(datadir, domain+".json")
			f, er = os.OpenFile(path, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
			if er != nil {
				bone.Log_Error("Cannot open data file for domain '%s' at path '%s'", domain, path)
				continue
			}
			datafiles[domain] = f
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
	for _, f := range datafiles {
		f.Close()
	}
	for k := range datafiles {
		delete(datafiles, k)
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
		shell.Set_Domain(bone.Config.Get_String("main", "domain", "main"))
		shell.Set_Command("setdomain", shell_set_domain)
		shell.Run()
		return
	}

	server := create_server()
	server.Run("0.0.0.0:3000")
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

	state[domain] = []*Event{}
	save_state()

	return shell.OK
}
