// Copyright (c) 2025 Start Codex SAS. All rights reserved.
// SPDX-License-Identifier: BUSL-1.1

import { redirect } from '@sveltejs/kit';
import { instance, auth } from '$lib/api';

export async function load() {
	const { initialized } = await instance.status();

	if (initialized) {
		// Already set up — check if user is logged in
		try {
			const me = await auth.me();
			if (me.authenticated) redirect(302, '/');
		} catch {
			// auth check failed, redirect to login
		}
		redirect(302, '/login');
	}
}
