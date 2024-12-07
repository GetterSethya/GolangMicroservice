<script lang="ts">
	import * as Form from '$lib/components/ui/form';
	import { Input } from '$lib/components/ui/input';
	import { superForm, type Infer, type SuperValidated } from 'sveltekit-superforms';
	import { registerSchema, type RegisterSchema } from '@/zod';
	import { zodClient } from 'sveltekit-superforms/adapters';
	import { cn } from '@/utils';
	import { AuthRepository } from '@/repository/auth';
	import { toast } from 'svelte-sonner';
	import { goto } from '$app/navigation';

	interface Props {
		data: SuperValidated<Infer<RegisterSchema>>;
	}

	let { data }: Props = $props();

	const form = superForm(data, {
		validators: zodClient(registerSchema),
		SPA: true,

		onUpdate: async ({ form }) => {
			if (form.valid) {
				const authRepo = new AuthRepository();
				const { status } = await authRepo.register({
					username: form.data.username,
					name: form.data.name,
					password: form.data.password,
					confirmPassword: form.data.confirmPassword
				});

				if (status === 201) {
					toast.success('Register success');
					goto('/login');
				} else {
					toast.error('Register failed');
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
	<Form.Field {form} name="name">
		<Form.Control>
			{#snippet children({ props })}
				<Form.Label>Name</Form.Label>
				<Input
					class={cn(props['aria-invalid'] ? 'border-destructive' : '')}
					{...props}
					bind:value={$formData.name}
					placeholder="Input your name"
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
	<Form.Field {form} name="confirmPassword">
		<Form.Control>
			{#snippet children({ props })}
				<Form.Label>Confirm Password</Form.Label>
				<Input
					type="password"
					{...props}
					class={cn(props['aria-invalid'] ? 'border-destructive' : '')}
					bind:value={$formData.confirmPassword}
					placeholder="Confirm your password"
				/>
			{/snippet}
		</Form.Control>
		<Form.FieldErrors class="font-normal" />
	</Form.Field>
	<Form.Button class="w-full">Register</Form.Button>
</form>
