// Copyright (c) 2025 Start Codex SAS. All rights reserved.
// SPDX-License-Identifier: BUSL-1.1

import { browser } from '$app/environment';
import { redirect } from '@sveltejs/kit';
import { restore } from '$lib/stores/auth';

export const ssr = false;

export async function load() {
	if (!browser) return;

	const user = await restore();
	if (!user) redirect(302, '/login');
}
