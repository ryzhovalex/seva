<script lang="ts">
	import { onMount } from "svelte";

    let domain = ""

    function submit(event: Event) {
        event.preventDefault()
        console.log(domain)
    }

    onMount(async () => {
        const response = await fetch("localhost:3000/Rpc/Domains/GetDomains")
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
            Domain:
            <select name="Domain" bind:value={domain}>
                <option value="Sports">Sports</option>
            </select>
        </div>

        <button type="submit" class="bg-c1 p-1">SUBMIT</button>
    </form>
</div>
