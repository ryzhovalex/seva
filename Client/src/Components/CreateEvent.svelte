<script lang="ts">
	import { onMount } from "svelte"
    import { Rpc } from "../lib/Rpc"

    let domain = ""
    let domains = []

    function submit(event) {
        event.preventDefault()
        console.log(domain)
    }

    async function onDomainSelected(event) {
        let events = await Rpc("Sevent/GetEvents")
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
    <form class="flex flex-col justify-start items-start mt-2 gap-2" on:submit={submit} on:change={onDomainSelected}>
        <div>
            Domain:
            <select name="Domain" bind:value={domain}>
                {#each domains as domain}
                    <option value="{domain}">{domain}</option>
                {/each}
            </select>
        </div>

        <button type="submit" class="bg-c1 p-1">SUBMIT</button>
    </form>
</div>
