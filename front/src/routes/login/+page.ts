// Copyright (c) 2025 Start Codex SAS. All rights reserved.
// SPDX-License-Identifier: BUSL-1.1

import { redirect } from '@sveltejs/kit';
import { instance } from '$lib/api';

export async function load() {
	const { initialized } = await instance.status();
	if (!initialized) redirect(302, '/setup');
}
