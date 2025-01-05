export async function Rpc(path: string, data: any = {}) {
    try {
        const response = await fetch(
            "http://localhost:3000/Rpc/" + path,
            {
                method: "POST",
                headers: {
                    "Content-Type": "application/json"
                },
                body: JSON.stringify(data)
            }
        )
        return (await response.json()).Body
    } catch (error) {
        return console.error(error)
    }
}
