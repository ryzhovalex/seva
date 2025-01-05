<script lang="ts">
	import type { Context } from "$lib/Commands";
	import { onMount } from "svelte";

    export let C: Context
</script>

<div>{C.Extra.Get("Text", "")}</div>
