import UnknownCommand from "../Components/UnknownCommand.svelte"

const COMMANDS = {
    clear: clear
}

function clear(c: Context) {
    c.ClearHistory()
}

export function ExecuteCommand(c: Context): any {
    let fn = COMMANDS[c.Prompt]
    if (fn == undefined) {
        if (c.Prompt != "") {
            return UnknownCommand
        }
        return null
    }
    return fn(c)
}

export interface Context {
    Prompt: string
    ClearHistory: () => void
}
