<script lang="ts">
	import { ContextExtra, ExecuteCommand, type Context } from "$lib/Commands";
	import { onMount } from "svelte";
	import HistoryPrompt from "./Responses/HistoryPrompt.svelte";

    let prompt = $state(null)
    let promptInputFocusTasked = false

    function handleKeydown(event) {
        if (event.key == "Enter") {
            send()
            return
        }
        if (event.key == "Control") {
            reset()
            return
        }
    }

    let historyComponents = $state([])
    let historyPrompts = $state([])
    let main = null
    let promptShown = $state(true)

    onMount(() => {
        main = document.getElementById("Main")
        document.addEventListener("keydown", handleKeydown)
    })

    $effect(() => {
        if (prompt != null && promptInputFocusTasked) {
            promptInputFocusTasked = false
            prompt.focus()
        }
    })


    function showPrompt() {
        promptShown = true
    }

    function focusPrompt() {
        promptInputFocusTasked = true
    }

    function hidePrompt() {
        promptShown = false
    }

    function reset() {
        historyComponents = []
        historyPrompts = []
        showPrompt()
        focusPrompt()
    }

    function send() {
        if (prompt == null) {
            return
        }

        let context: Context = {
            Prompt: prompt.value,
            Extra: new ContextExtra(),
            Send: receiveComponent,
            ShowPrompt: showPrompt,
            HidePrompt: hidePrompt,
            FocusPrompt: focusPrompt,
            Reset: reset
        }
        historyComponents = [
            ...historyComponents, {Context: context, Component: HistoryPrompt}]
        historyPrompts = [...historyPrompts, prompt]

        setTimeout(() => {
            main.scrollTop = main.scrollHeight
        }, 100);
        prompt.value = ""

        ExecuteCommand(context)
    }

    function receiveComponent(component: any) {
        if (component != null) {
            historyComponents = [...historyComponents, {Context: this, Component: component}]
        }
    }
</script>

<div class="flex flex-col gap-4 mb-4">
    {#each historyComponents as item}
        <svelte:component this={item.Component} C={item.Context}/>
    {/each}
</div>

{#if promptShown}
    <div class="flex flex-row items-center gap-2">
        <div>></div>
        <input type="text" class="w-full" bind:this={prompt}/>
    </div>
{/if}
