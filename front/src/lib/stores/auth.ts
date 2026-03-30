// Copyright (c) 2025 Start Codex SAS. All rights reserved.
// SPDX-License-Identifier: BUSL-1.1

import { writable } from 'svelte/store';
import { goto } from '$app/navigation';
import { auth as authApi, type User } from '$lib/api';

const _store = writable<User | null>(null);

export const currentUser = { subscribe: _store.subscribe };

export function login(user: User): void {
	_store.set(user);
}

export async function signIn(email: string, password: string): Promise<User> {
	const user = await authApi.login({ email, password });
	_store.set(user);
	return user;
}

export async function restore(): Promise<User | null> {
	const res = await authApi.me();
	if (res.authenticated && res.user) {
		_store.set(res.user);
		return res.user;
	}
	_store.set(null);
	return null;
}

export async function logout(): Promise<void> {
	try {
		await authApi.logout();
	} catch {
		// best effort
	}
	_store.set(null);
	goto('/login');
}
