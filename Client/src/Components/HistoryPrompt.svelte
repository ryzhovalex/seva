<script lang="ts">
	import type { Context } from "$lib/Commands";

    export let Args: Context;
</script>

<div class="flex flex-row items-center gap-2 mb-4">> {Args.Prompt}</div>
