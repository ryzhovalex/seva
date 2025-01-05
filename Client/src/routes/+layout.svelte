<script lang="ts">
	import '../app.css';
</script>

<main class="text-c4 text-2xl font-bold bg-c5 h-screen w-screen bg-cover bg-center bg-no-repeat overflow-y-scroll" id="Main">
	<!-- Use py-4 here for a scroll to correctly auto-scroll to the bottom. -->
	<div class="container mx-auto flex flex-col py-4">
		<div class="text-4xl self-center text-center mb-4">
			S.E.V.A.
		</div>
		<slot />
	</div>
</main>
