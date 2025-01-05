import CreateEvent from "../Components/CreateEvent.svelte"
import UnknownCommand from "../Components/UnknownCommand.svelte"

const COMMANDS = {
    "clear": clear,
    "event.create": event_create
}

function clear(c: Context) {
    c.ClearHistory()
}

function event_create(c: Context) {
    c.ClearHistory()
    c.HidePrompt()
}

export function ExecuteCommand(c: Context): any {
    let fn = COMMANDS[c.Prompt]
    if (fn == undefined) {
        if (c.Prompt != "") {
            c.Send(UnknownCommand)
        }
        c.Send(null)
    }
    c.Send(fn(c))
}

export interface Context {
    Prompt: string
    Send: (component) => void
    ShowPrompt: () => void
    HidePrompt: () => void
    ClearHistory: () => void
}
