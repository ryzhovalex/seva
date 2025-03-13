package main

import (
	"encoding/json"
	"flag"
	"os"
	"path/filepath"
	"seva/lib/bone"
	"seva/lib/shell"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

const (
	OK = iota
	ERROR
)

type Event struct {
	Domain string `json:"domain"`
	// Particular number of event in domain.
	Order int `json:"order"`
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
		bone.Log_Error("During state reading, cannot read userdir")
		return ERROR
	}
	for _, file := range files {
		if filepath.Ext(file.Name()) == ".json" {
			path := filepath.Join()
			data, er := os.ReadFile(path)
			if er != nil {
				bone.Log_Error("During state reading, cannot read file '%s'", file.Name())
				return ERROR
			}
			domain := filepath.Base(file.Name())
			events := []*Event{}
			state[domain] = events
			er = json.Unmarshal(data, &events)
			if er != nil {
				bone.Log_Error("Cannot unmarshal file '%s'", file.Name())
				return ERROR
			}
			bone.Log("Loaded domain '%s'", domain)
		}
	}
	return OK
}

func save_state() {
	for domain, events := range state {
		data, er := json.MarshalIndent(events, "", "\t")
		if er != nil {
			bone.Log_Error("Error marshalling state to json for domain '%s'", domain)
			return
		}

		f, ok := datafiles[domain]
		if !ok {
			f, er = os.OpenFile(bone.Userdir("data", domain+".json"), os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
			if er != nil {
				bone.Log_Error("Cannot open data file for domain '%s'", domain)
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
	shell.Init("seva")

	e := read_state()
	if e != OK {
		bone.Log_Error("During state reading, an error occurred: %d", e)
		os.Exit(e)
		return
	}

	if *shell_enabled {
		shell.Run()
		return
	}

	server := create_server()
	server.Run("0.0.0.0:3000")
}
