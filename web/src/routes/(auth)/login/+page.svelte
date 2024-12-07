<script lang="ts">
	import * as Card from '$lib/components/ui/card/index.js';
	import type { PageData } from './$types';
	import LoginForm from './login-form.svelte';

	interface Props {
		data: PageData;
	}

	let { data }: Props = $props();
</script>

<Card.Root>
	<Card.Header>
		<Card.Title>Login</Card.Title>
		<Card.Description>Insert your account information bellow to continue</Card.Description>
	</Card.Header>
	<Card.Content>
		<LoginForm data={data.form} />
	</Card.Content>
	<Card.Footer>
		<span
			>Dont have an account yet?
			<a class="underline" href="/register">Register</a>
		</span>
	</Card.Footer>
</Card.Root>
