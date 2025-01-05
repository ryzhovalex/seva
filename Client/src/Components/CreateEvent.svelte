<script lang="ts">
	import { onMount } from "svelte"
    import { Rpc } from "../lib/Rpc"

    let chosenDomain = ""
    let chosenSpec: {[Key: string]: {Type: string, Fields: any[]}} = null

    let domains = []
    let specs = null
    let eventTypes = []

    function submit(event) {
        event.preventDefault()
    }

    async function onDomainSelected(event) {
        specs = await Rpc("Sevent/GetSpecs", {Domain: chosenDomain})
        eventTypes = Object.keys(specs)
        if (eventTypes.length > 0) {
            chosenSpec = specs[eventTypes[0]]
        }
    }

    async function onEventTypeSelected(event) {
        chosenSpec = specs[event.target.value]
        console.log(chosenSpec)
    }

    onMount(async () => {
        domains = await Rpc("Domains/GetDomains")
    })
</script>

<div class="flex flex-col">
    <div>
        CREATE EVENT
        <br/>
        ------------
    </div>
    <form class="flex flex-col justify-start items-start mt-2 gap-2" on:submit={submit}>
        <div>
            DOMAIN:
            <select name="Domain" bind:value={chosenDomain} on:change={onDomainSelected}>
                {#each domains as domain}
                    <option value="{domain}">{domain}</option>
                {/each}
            </select>
        </div>

        {#if eventTypes.length > 0}
            <div>
                EVENT TYPE:
                <select name="EventType" on:change={onEventTypeSelected}>
                    {#each eventTypes as et}
                        <option value="{et}">{et}</option>
                    {/each}
                </select>
            </div>
        {/if}

        {#if chosenSpec !== null}
            <div>
                _BODY_
            </div>
            <div class="ml-16 flex flex-col gap-2">
                {#each Object.entries(chosenSpec) as [key, field]}
                    <div>
                        {key}:
                        {#if field.Type == "string"}
                            <input type="text" name="{key}"/>
                        {:else if field.Type == "number"}
                            <input type="number" name="{key}" value="0"/>
                        {:else if field.Type == "boolean"}
                            <input type="checkbox" name="{key}" value="false" class="w-6 h-6"/>
                        {:else if field.Type == "array"}
                            ...array here
                        {:else if field.Type == "object"}
                            ...object here
                        {/if}
                    </div>
                {/each}
            </div>
        {/if}

        <button type="submit" class="bg-c0 p-1 hover:bg-c1">SUBMIT</button>
    </form>
</div>
