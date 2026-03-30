// Copyright (c) 2025 Start Codex SAS. All rights reserved.
// SPDX-License-Identifier: BUSL-1.1

import { redirect } from '@sveltejs/kit';
import { workspaces as workspacesApi, projects as projectsApi, ApiError } from '$lib/api';
import type { LayoutLoad } from './$types';

export const load: LayoutLoad = async ({ params, parent }) => {
	await parent();

	let list;
	try {
		list = await workspacesApi.list();
	} catch (err) {
		if (err instanceof ApiError && err.status === 401) redirect(302, '/login');
		throw err;
	}

	const workspace = list.find((w) => w.slug === params.workspace);
	if (!workspace) redirect(302, '/');

	let projectList;
	try {
		projectList = await projectsApi.list(workspace.id);
	} catch (err) {
		if (err instanceof ApiError && err.status === 401) redirect(302, '/login');
		throw err;
	}

	return {
		workspace,
		workspaceProjects: projectList ?? [],
		workspaceList: list
	};
};
