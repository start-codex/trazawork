<!-- Copyright (c) 2025 Start Codex SAS. All rights reserved. -->
<!-- SPDX-License-Identifier: BUSL-1.1 -->

<script lang="ts">
	import { auth, ApiError } from '$lib/api';
	import * as Card from '$lib/components/ui/card/index.js';
	import { Button } from '$lib/components/ui/button/index.js';
	import { Input } from '$lib/components/ui/input/index.js';
	import * as m from '$lib/paraglide/messages';
	import { i18n } from '$lib/i18n.svelte';

	const t = $derived.by(() => {
		i18n.locale;
		return {
			title: m.account_change_password(),
			current: m.account_current_password(),
			newPw: m.account_new_password(),
			confirm: m.account_confirm_password(),
			changed: m.account_password_changed(),
			mismatch: m.account_passwords_mismatch(),
			invalid: m.account_current_password_invalid(),
			tooShort: m.account_password_too_short(),
			save: m.account_save(),
			saving: m.account_saving()
		};
	});

	let currentPassword = $state('');
	let newPassword = $state('');
	let confirmPassword = $state('');
	let saving = $state(false);
	let success = $state(false);
	let error = $state('');

	function clearForm() {
		currentPassword = '';
		newPassword = '';
		confirmPassword = '';
	}

	async function handleSubmit(e: SubmitEvent) {
		e.preventDefault();
		error = '';
		success = false;

		if (newPassword !== confirmPassword) {
			error = t.mismatch;
			return;
		}
		if (newPassword.length < 8) {
			error = t.tooShort;
			return;
		}

		saving = true;
		try {
			await auth.changePassword({
				current_password: currentPassword,
				new_password: newPassword
			});
			success = true;
			clearForm();
			setTimeout(() => { success = false; }, 5000);
		} catch (err) {
			if (err instanceof ApiError && err.status === 401) {
				error = t.invalid;
			} else {
				error = err instanceof Error ? err.message : 'Failed to change password';
			}
		} finally {
			saving = false;
		}
	}
</script>

<Card.Root>
	<Card.Header>
		<Card.Title>{t.title}</Card.Title>
	</Card.Header>
	<Card.Content>
		<form onsubmit={handleSubmit} class="space-y-4">
			<div class="space-y-1.5">
				<label for="current-pw" class="text-sm font-medium">{t.current}</label>
				<Input id="current-pw" type="password" bind:value={currentPassword} required />
			</div>
			<div class="space-y-1.5">
				<label for="new-pw" class="text-sm font-medium">{t.newPw}</label>
				<Input id="new-pw" type="password" bind:value={newPassword} required />
			</div>
			<div class="space-y-1.5">
				<label for="confirm-pw" class="text-sm font-medium">{t.confirm}</label>
				<Input id="confirm-pw" type="password" bind:value={confirmPassword} required />
			</div>

			{#if error}
				<p class="text-sm text-destructive">{error}</p>
			{/if}
			{#if success}
				<p class="text-sm text-green-600">{t.changed}</p>
			{/if}

			<Button type="submit" disabled={saving || !currentPassword || !newPassword || !confirmPassword}>
				{saving ? t.saving : t.save}
			</Button>
		</form>
	</Card.Content>
</Card.Root>
