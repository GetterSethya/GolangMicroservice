<script lang="ts">
	import '../app.css';
	import { Toaster } from '$lib/components/ui/sonner';
	let { children } = $props();
</script>

<div class="w-screen h-screen flex">
	<Toaster richColors />
	{@render children()}
</div>
