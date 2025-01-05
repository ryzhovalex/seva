<script lang="ts">
	import { ExecuteCommand, type Context } from "$lib/Commands";
	import { onMount } from "svelte";
	import HistoryPrompt from "./HistoryPrompt.svelte";

    function handleKeydown(event) {
        if (event.key == "Enter") {
            send()
        }
    }

    let history = []
    let main = null
    onMount(() => {
        main = document.getElementById("Main")
    })

    function send() {
        context = {
            Prompt: prompt
        }
        ExecuteCommand(context)
        history = [...history, HistoryPrompt]
        setTimeout(() => {
            main.scrollTop = main.scrollHeight
        }, 100);
        prompt = ""
    }

    let prompt = ""
    let context: Context = null
</script>

{#each history as item}
    <svelte:component this={item} Context={context}/>
{/each}

<div class="flex flex-row items-center gap-2">
    <div>></div>
    <input type="text" class="w-full" on:keydown={handleKeydown} bind:value={prompt}/>
</div>
