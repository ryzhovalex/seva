const COMMANDS = {
    clear: clear
}

function clear() {

}

export function ExecuteCommand(context: Context): any {
    let fn = COMMANDS[context.Prompt]
    if (fn == undefined) {
        if (context.Prompt != "") {
        }
    }
}

export interface Context {
    Prompt: string
}
