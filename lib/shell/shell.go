// Organizes interactive shell.
package shell

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"seva/lib/bone"
	"strconv"
	"strings"
)

type Command_Handler func(ctx *Command_Context) int

var commands = map[string]Command_Handler{}
var domain string = ""

const (
	OK = iota
	ERROR
)

var hooks = []any{}
var prompted = false
var prompted_callback func(answer bool) int = nil

type Command_Context struct {
	Raw_Input    string
	Command_Name string
	args         map[string]string
}

func (c *Command_Context) parse(raw_args []string) {
	c.args = map[string]string{}

	prev_arg := ""
	// Collect buffer until the next argument
	buffer := ""
	for _, a := range raw_args {
		if strings.HasPrefix(a, "-") {
			if prev_arg != "" {
				// Assign buffer even if it's empty (e.g. for flag arguments)
				c.args[prev_arg] = buffer
			} else if buffer != "" {
				// If previous argument were not defined, it means we're collecting
				// input to the main command - in this case we save it under `_` key.
				// Same operation should be done both before next argument or before end
				// of the arguments.
				c.args["_"] = buffer
			}
			prev_arg = a
			buffer = ""
			continue
		}
		buffer += a
	}

	// End of input, assign even empty buffer
	if prev_arg != "" {
		c.args[prev_arg] = buffer
	} else {
		c.args["_"] = buffer
	}
}

// IF BUFFER IS EMPTY RETURN DEFAULT
func (c *Command_Context) Arg_String(key string, default_ string) string {
	if !strings.HasPrefix(key, "-") && key != "_" {
		bone.Log_Error("Unable to search argument via non-flag key '%s'", key)
		return default_
	}

	buffer, ok := c.args[key]
	if !ok {
		return default_
	}
	if buffer == "" {
		return default_
	}
	return buffer
}

// Returns true if key exists.
func (c *Command_Context) Arg_Bool(key string, default_ bool) bool {
	if !strings.HasPrefix(key, "-") && key != "_" {
		bone.Log_Error("Unable to search argument via non-flag key '%s'", key)
		return default_
	}

	_, ok := c.args[key]
	if !ok {
		return default_
	}
	return true
}

func (c *Command_Context) Arg_Int(key string, default_ int) int {
	if !strings.HasPrefix(key, "-") && key != "_" {
		bone.Log_Error("Unable to search argument via non-flag key '%s'", key)
		return default_
	}

	buffer, ok := c.args[key]
	if !ok {
		return default_
	}
	r, er := strconv.Atoi(buffer)
	if er != nil {
		bone.Log_Error("Unable to convert argument '%s' value '%s' to integer", key, buffer)
		return default_
	}
	return r
}

func (c *Command_Context) Arg_Float(key string, default_ float64) float64 {
	if !strings.HasPrefix(key, "-") && key != "_" {
		bone.Log_Error("Unable to search argument via non-flag key '%s'", key)
		return default_
	}

	buffer, ok := c.args[key]
	if !ok {
		return default_
	}
	r, er := strconv.ParseFloat(buffer, 64)
	if er != nil {
		bone.Log_Error("Unable to convert argument '%s' value '%s' to float", key, buffer)
		return default_
	}
	return r
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

	raw_args := []string{}
	if len(input_parts) > 1 {
		raw_args = input_parts[1:]
	}
	ctx := Command_Context{
		Raw_Input:    input,
		Command_Name: command_name,
	}
	ctx.parse(raw_args)

	e := cmd(&ctx)
	if e > 0 {
		bone.Log_Error("While calling a command `%s`, an error occured: %s", command_name, bone.Tr_Code(e))
	}
}

func Init() {
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
		fmt.Printf("\033[33m(%s)\033[0m\033[35m%s\033[0m ", domain, final_sign)
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

var domain_regex = regexp.MustCompile("^[a-z0-9_]*$")

func Get_Domain() string {
	return domain
}

func Set_Domain(d string) int {
	if !domain_regex.MatchString(d) {
		bone.Log_Error("Incorrect domain '%s'", d)
		return ERROR
	}
	domain = d
	return OK
}
