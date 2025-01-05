import CreateEvent from "../Components/Responses/CreateEvent.svelte"
import Help from "../Components/Responses/Help.svelte"
import UnknownCommand from "../Components/Responses/UnknownCommand.svelte"

const COMMANDS = {
    "clear": clear,
    "event.create": event_create,
    "help": help,
}

function help(c: Context) {
    c.Extra.Set("Commands", Object.keys(COMMANDS))
    c.Send(Help)
}

function clear(c: Context) {
    c.Reset()
}

function event_create(c: Context) {
    c.Reset()
    c.HidePrompt()
    c.Send(CreateEvent)
}

export function ExecuteCommand(c: Context): any {
    let fn = COMMANDS[c.Prompt]
    if (fn == undefined) {
        if (c.Prompt != "") {
            c.Send(UnknownCommand)
            return
        }
        c.Send(null)
        return
    }
    c.Send(fn(c))
}

export class ContextExtra {
    private data = {}

    public Get(key: string, default_: any): any {
        let v = this.data[key]
        if (v == undefined) {
            return default_
        }
        return v
    }

    public Set(key: string, value: any) {
        this.data[key] = value
    }
}

export interface Context {
    Prompt: string
    Extra: ContextExtra
    Send: (component) => void
    ShowPrompt: () => void
    HidePrompt: () => void
    FocusPrompt: () => void
    Reset: () => void
}
