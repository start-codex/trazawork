<!-- Copyright (c) 2025 Start Codex SAS. All rights reserved. -->
<!-- SPDX-License-Identifier: BUSL-1.1 -->

<script lang="ts">
	import ChangePasswordForm from '$lib/components/change-password-form.svelte';
	import * as m from '$lib/paraglide/messages';
	import { i18n } from '$lib/i18n.svelte';

	const title = $derived.by(() => { i18n.locale; return m.settings_account_title(); });
</script>

<div class="space-y-6">
	<div>
		<h2 class="text-lg font-semibold">{title}</h2>
		<hr class="mt-3 border-border" />
	</div>

	<ChangePasswordForm />
</div>
