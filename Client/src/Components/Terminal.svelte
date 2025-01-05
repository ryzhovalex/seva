<script lang="ts">
	import { ExecuteCommand, type Context } from "$lib/Commands";
	import { onMount } from "svelte";
	import HistoryPrompt from "./HistoryPrompt.svelte";

    function handleKeydown(event) {
        if (event.key == "Enter") {
            send()
        }
    }

    let historyComponents = []
    let historyPrompts = []
    let main = null
    onMount(() => {
        main = document.getElementById("Main")
    })

    function send() {
        let context = {
            Prompt: prompt,
            Send: receiveComponent,
            ShowPrompt: () => {},
            HidePrompt: () => {},
            ClearHistory: () => {
                historyComponents = []
                historyPrompts = []
            }
        }
        historyComponents = [
            ...historyComponents, {Context: context, Component: HistoryPrompt}]
        historyPrompts = [...historyPrompts, prompt]

        setTimeout(() => {
            main.scrollTop = main.scrollHeight
        }, 100);
        prompt = ""

        ExecuteCommand(context)
    }

    function receiveComponent(component: any) {
        if (component != null) {
            historyComponents = [...historyComponents, {Context: this, Component: component}]
        }
    }

    let prompt = ""
</script>

<div class="flex flex-col gap-4 mb-4">
    {#each historyComponents as item}
        <svelte:component this={item.Component} C={item.Context}/>
    {/each}
</div>

<div class="flex flex-row items-center gap-2">
    <div>></div>
    <input type="text" class="w-full" on:keydown={handleKeydown} bind:value={prompt}/>
</div>
