<!-- Copyright (c) 2025 Start Codex SAS. All rights reserved. -->
<!-- SPDX-License-Identifier: BUSL-1.1 -->

<script lang="ts">
	import { goto } from '$app/navigation';
	import { Button } from '$lib/components/ui/button/index.js';
	import * as Card from '$lib/components/ui/card/index.js';
	import {
		FieldGroup,
		Field,
		FieldLabel
	} from '$lib/components/ui/field/index.js';
	import { Input } from '$lib/components/ui/input/index.js';
	import { instance } from '$lib/api';
	import { login } from '$lib/stores/auth';
	import * as m from '$lib/paraglide/messages';
	import { i18n } from '$lib/i18n.svelte';

	const t = $derived.by(() => {
		i18n.locale;
		return {
			title: m.setup_title(),
			description: m.setup_description(),
			name: m.setup_name(),
			email: m.setup_email(),
			password: m.setup_password(),
			confirmPassword: m.setup_confirm_password(),
			submit: m.setup_submit(),
			creating: m.setup_creating(),
			mismatch: m.setup_passwords_mismatch(),
			error: m.setup_error()
		};
	});

	let name = $state('');
	let email = $state('');
	let password = $state('');
	let confirmPassword = $state('');
	let errorMessage = $state('');
	let loading = $state(false);

	async function handleSubmit(e: SubmitEvent) {
		e.preventDefault();
		errorMessage = '';

		if (password !== confirmPassword) {
			errorMessage = t.mismatch;
			return;
		}

		loading = true;
		try {
			const user = await instance.bootstrap({ email, name, password });
			login(user);
			goto('/');
		} catch (err) {
			errorMessage = err instanceof Error ? err.message : t.error;
		} finally {
			loading = false;
		}
	}
</script>

<div class="bg-muted flex min-h-svh flex-col items-center justify-center gap-6 p-6 md:p-10">
	<div class="flex w-full max-w-sm flex-col gap-6">
		<a href="/setup" class="flex items-center gap-2 self-center font-medium">
			<div class="bg-primary text-primary-foreground flex size-6 items-center justify-center rounded-md">
				<svg xmlns="http://www.w3.org/2000/svg" class="size-4" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
					<rect width="8" height="8" x="2" y="2" rx="2" /><rect width="8" height="8" x="14" y="2" rx="2" />
					<rect width="8" height="8" x="2" y="14" rx="2" /><rect width="8" height="8" x="14" y="14" rx="2" />
				</svg>
			</div>
			Tookly
		</a>
		<Card.Root>
			<Card.Header class="text-center">
				<Card.Title class="text-xl">{t.title}</Card.Title>
				<Card.Description>{t.description}</Card.Description>
			</Card.Header>
			<Card.Content>
				<form onsubmit={handleSubmit}>
					<FieldGroup>
						<Field>
							<FieldLabel for="setup-name">{t.name}</FieldLabel>
							<Input id="setup-name" type="text" placeholder="Admin" required bind:value={name} />
						</Field>
						<Field>
							<FieldLabel for="setup-email">{t.email}</FieldLabel>
							<Input id="setup-email" type="email" placeholder="admin@example.com" required bind:value={email} />
						</Field>
						<Field>
							<FieldLabel for="setup-password">{t.password}</FieldLabel>
							<Input id="setup-password" type="password" required bind:value={password} />
						</Field>
						<Field>
							<FieldLabel for="setup-confirm">{t.confirmPassword}</FieldLabel>
							<Input id="setup-confirm" type="password" required bind:value={confirmPassword} />
						</Field>
						{#if errorMessage}
							<p class="text-destructive text-sm">{errorMessage}</p>
						{/if}
						<Field>
							<Button type="submit" class="w-full" disabled={loading}>
								{loading ? t.creating : t.submit}
							</Button>
						</Field>
					</FieldGroup>
				</form>
			</Card.Content>
		</Card.Root>
	</div>
</div>
