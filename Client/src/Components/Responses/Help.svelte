<script lang="ts">
	import type { Context } from "$lib/Commands";

    export let C: Context
</script>

<div class="flex flex-col gap-2">
    AVAILABLE COMMANDS
    <div class="flex flex-col ml-8 gap-2">
        {#each C.Extra.Commands as command}
            <div>
                * {command}
            </div>
        {/each}
    </div>
</div>
