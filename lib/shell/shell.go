// Organizes interactive shell.
package shell

import (
	"bufio"
	"fmt"
	"os"
	"seva/lib/bone"
	"strings"
)

type Command_Handler func(ctx *Command_Context) int

var commands = map[string]Command_Handler{}
var project_name string

type Command_Context struct {
	Raw_Input    string
	Command_Name string
	Args         []string
}

const (
	OK = iota
	ERROR
)

var hooks = []any{}
var prompted = false
var prompted_callback func(answer bool) int = nil

func (ctx *Command_Context) Has_Arg(arg string) bool {
	for _, a := range ctx.Args {
		if a == arg {
			return true
		}
	}
	return false
}

func (ctx *Command_Context) Has_Arg_Index(arg string) (bool, int) {
	for i, a := range ctx.Args {
		if a == arg {
			return true, i
		}
	}
	return false, -1
}

func Set_Hooks[T any](items []T) {
	Clear_Hooks()
	for _, i := range items {
		hooks = append(hooks, i)
	}
}

func Clear_Hooks() {
	hooks = []any{}
}

func Answer_Prompt(answer bool) {
	if !prompted {
		bone.Log_Error("Inactive prompt")
		return
	}
	prompted = false
	e := prompted_callback(answer)
	if e != OK {
		bone.Log_Error("During prompted callback, an error #%d occured", e)
	}
	prompted_callback = nil
}

func Prompt(text string, callback func(answer bool) int) {
	if prompted {
		bone.Log_Error("Already prompted")
		return
	}
	prompted = true
	prompted_callback = callback
	bone.Log(text + " [Y/N]")
}

func process_input(input string) {
	// Quoted strings are not yet supported - they will be separated as everything else.
	input_parts := strings.Fields(input)
	if len(input_parts) == 0 {
		return
	}

	if prompted {
		var answer bool
		switch input {
		case "y":
			answer = true
		case "n":
			answer = false
		case "Y":
			answer = true
		case "N":
			answer = false
		default:
			bone.Log("Type answer 'Y' or 'N'")
			return
		}
		Answer_Prompt(answer)
		return
	}

	command_name := input_parts[0]

	cmd, ok := commands[command_name]
	if !ok {
		bone.Log_Error("Unrecognized command: " + input)
		return
	}

	args := []string{}
	if len(input_parts) > 1 {
		args = input_parts[1:]
	}
	ctx := Command_Context{
		Raw_Input:    input,
		Command_Name: command_name,
		Args:         args,
	}

	e := cmd(&ctx)
	if e > 0 {
		bone.Log_Error("While calling a command `%s`, an error occured: %s", command_name, bone.Tr_Code(e))
	}
}

func Init(project_name_ string) {
	project_name = project_name_
}

func Set_Command(key string, handler Command_Handler) {
	commands[key] = handler
}

func Run() {
	console_reader := bufio.NewReader(os.Stdin)

	// Main loop is blocking on input, other background tasks are goroutines.
	for {
		var final_sign = ">"
		if prompted {
			final_sign = "?"
		}
		fmt.Printf("\033[33m(%s)\033[0m\033[35m%s\033[0m ", project_name, final_sign)
		input, er := console_reader.ReadString('\n')
		if er != nil {
			if er.Error() != "EOF" {
				bone.Log_Error("Unexpected error occured while reading console: %s", er)
			}
			return
		}
		input = strings.TrimSpace(input)
		if input == "q" {
			return
		}
		process_input(input)
	}
}
