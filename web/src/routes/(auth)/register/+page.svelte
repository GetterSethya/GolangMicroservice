<script lang="ts">
	import * as Card from '$lib/components/ui/card/index.js';
	import type { PageData } from './$types';
	import RegisterForm from './register-form.svelte';

	interface Props {
		data: PageData;
	}

	let { data }: Props = $props();
</script>

<Card.Root>
	<Card.Header>
		<Card.Title>Register</Card.Title>
		<Card.Description>Create your account by filling the form bellow</Card.Description>
	</Card.Header>
	<Card.Content>
		<RegisterForm data={data.form} />
	</Card.Content>
	<Card.Footer>
		<span
			>Already have an account?
			<a class="underline" href="/login">Login</a>
		</span>
	</Card.Footer>
</Card.Root>
