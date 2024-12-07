<script lang="ts">
	import * as Form from '$lib/components/ui/form';
	import { Input } from '$lib/components/ui/input';
	import { superForm, type Infer, type SuperValidated } from 'sveltekit-superforms';
	import { loginSchema, type LoginSchema } from '@/zod';
	import { zodClient } from 'sveltekit-superforms/adapters';
	import { cn } from '@/utils';
	import { AuthRepository } from '@/repository/auth';
	import { toast } from 'svelte-sonner';
	import { goto } from '$app/navigation';

	interface Props {
		data: SuperValidated<Infer<LoginSchema>>;
	}

	let { data }: Props = $props();

	const form = superForm(data, {
		validators: zodClient(loginSchema),
		SPA: true,
		onUpdate: async ({ form }) => {
			if (form.valid) {
				const authRepo = new AuthRepository();
				const { res } = await authRepo.login({
					username: form.data.username,
					password: form.data.password
				});

				if (res.data?.accessToken && res.data.refreshToken) {
					localStorage.setItem('accessToken', res.data.accessToken);
					localStorage.setItem('refreshToken', res.data.refreshToken);

					toast.success('Login success');
					goto('/', { invalidateAll: true });
				} else {
					toast.error('Login failed');
				}
			}
		}
	});

	const { form: formData, enhance } = form;
</script>

<form action="/login" method="post" use:enhance>
	<Form.Field {form} name="username">
		<Form.Control>
			{#snippet children({ props })}
				<Form.Label>Username</Form.Label>
				<Input
					class={cn(props['aria-invalid'] ? 'border-destructive' : '')}
					{...props}
					bind:value={$formData.username}
					placeholder="Input your username"
				/>
			{/snippet}
		</Form.Control>
		<Form.FieldErrors class="font-normal" />
	</Form.Field>
	<Form.Field {form} name="password">
		<Form.Control>
			{#snippet children({ props })}
				<Form.Label>Password</Form.Label>
				<Input
					type="password"
					{...props}
					class={cn(props['aria-invalid'] ? 'border-destructive' : '')}
					bind:value={$formData.password}
					placeholder="Input your password"
				/>
			{/snippet}
		</Form.Control>
		<Form.FieldErrors class="font-normal" />
	</Form.Field>
	<Form.Button class="w-full">Login</Form.Button>
</form>
