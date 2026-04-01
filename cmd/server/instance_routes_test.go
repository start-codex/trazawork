// Copyright (c) 2025 Start Codex SAS. All rights reserved.
// SPDX-License-Identifier: BUSL-1.1
// Use of this software is governed by the Business Source License 1.1
// included in the LICENSE file at the root of this repository.

package main

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/start-codex/tookly/internal/testpg"
)

// resetInstance resets the instance state to uninitialized and cleans up
// any users/sessions created during tests. Must be called in t.Cleanup.
func resetInstance(t *testing.T, db *sqlx.DB) {
	t.Helper()
	ctx := context.Background()
	db.ExecContext(ctx, `DELETE FROM sessions`)
	db.ExecContext(ctx, `DELETE FROM workspace_members`)
	db.ExecContext(ctx, `DELETE FROM app_users`)
	db.ExecContext(ctx, `UPDATE instance_config SET value = 'false', updated_at = NOW() WHERE key = 'initialized'`)
}

func setupFreshInstanceServer(t *testing.T) (*httptest.Server, *sqlx.DB) {
	t.Helper()
	db := testpg.Open(t)
	testpg.EnsureMigrated(t, db)
	resetInstance(t, db)
	t.Cleanup(func() { resetInstance(t, db) })

	handler := newAPIHandler(db)
	srv := httptest.NewServer(handler)
	t.Cleanup(srv.Close)
	return srv, db
}

func doInstancePost(t *testing.T, srv *httptest.Server, path string, body any) (*http.Response, map[string]any) {
	t.Helper()
	var buf bytes.Buffer
	if body != nil {
		json.NewEncoder(&buf).Encode(body)
	}
	req, _ := http.NewRequest("POST", srv.URL+path, &buf)
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("do request: %v", err)
	}
	defer resp.Body.Close()
	var result map[string]any
	json.NewDecoder(resp.Body).Decode(&result)
	return resp, result
}

func doInstanceGet(t *testing.T, srv *httptest.Server, path string) (*http.Response, map[string]any) {
	t.Helper()
	req, _ := http.NewRequest("GET", srv.URL+path, nil)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("do request: %v", err)
	}
	defer resp.Body.Close()
	var result map[string]any
	json.NewDecoder(resp.Body).Decode(&result)
	return resp, result
}

func TestInstanceStatus_FreshDB(t *testing.T) {
	srv, _ := setupFreshInstanceServer(t)

	resp, body := doInstanceGet(t, srv, "/instance/status")
	if resp.StatusCode != 200 {
		t.Fatalf("status = %d, want 200", resp.StatusCode)
	}
	data := body["data"].(map[string]any)
	if init, ok := data["initialized"].(bool); !ok || init {
		t.Fatalf("initialized = %v, want false", data["initialized"])
	}
}

func TestBootstrap_FreshInstance(t *testing.T) {
	srv, db := setupFreshInstanceServer(t)
	suffix := testpg.UniqueSuffix(t, db)

	resp, body := doInstancePost(t, srv, "/instance/bootstrap", map[string]string{
		"email":    suffix + "@test.local",
		"name":     "Admin " + suffix,
		"password": "securepass123",
	})
	if resp.StatusCode != 201 {
		t.Fatalf("bootstrap = %d, want 201, body: %v", resp.StatusCode, body)
	}

	// Check cookie
	found := false
	for _, c := range resp.Cookies() {
		if c.Name == "session_id" && c.Value != "" {
			found = true
		}
	}
	if !found {
		t.Fatal("session_id cookie not set")
	}

	// Check is_instance_admin
	data := body["data"].(map[string]any)
	if admin, ok := data["is_instance_admin"].(bool); !ok || !admin {
		t.Fatalf("is_instance_admin = %v, want true", data["is_instance_admin"])
	}

	// Status should be initialized
	resp2, body2 := doInstanceGet(t, srv, "/instance/status")
	if resp2.StatusCode != 200 {
		t.Fatalf("status = %d", resp2.StatusCode)
	}
	data2 := body2["data"].(map[string]any)
	if init, ok := data2["initialized"].(bool); !ok || !init {
		t.Fatalf("initialized = %v, want true", data2["initialized"])
	}
}

func TestBootstrap_AlreadyInitialized(t *testing.T) {
	srv, db := setupFreshInstanceServer(t)
	suffix := testpg.UniqueSuffix(t, db)

	resp, _ := doInstancePost(t, srv, "/instance/bootstrap", map[string]string{
		"email":    suffix + "a@test.local",
		"name":     "Admin",
		"password": "securepass123",
	})
	if resp.StatusCode != 201 {
		t.Fatalf("first bootstrap = %d, want 201", resp.StatusCode)
	}

	resp2, _ := doInstancePost(t, srv, "/instance/bootstrap", map[string]string{
		"email":    suffix + "b@test.local",
		"name":     "Admin2",
		"password": "securepass456",
	})
	if resp2.StatusCode != 409 {
		t.Fatalf("second bootstrap = %d, want 409", resp2.StatusCode)
	}
}

func TestBootstrap_InvalidEmail(t *testing.T) {
	srv, _ := setupFreshInstanceServer(t)

	resp, _ := doInstancePost(t, srv, "/instance/bootstrap", map[string]string{
		"email":    "not-an-email",
		"name":     "Admin",
		"password": "securepass123",
	})
	if resp.StatusCode != 422 {
		t.Fatalf("bootstrap invalid email = %d, want 422", resp.StatusCode)
	}
}

func TestCreateUser_BeforeBootstrap(t *testing.T) {
	srv, _ := setupFreshInstanceServer(t)

	resp, _ := doInstancePost(t, srv, "/users", map[string]string{
		"email":    "user@test.local",
		"name":     "User",
		"password": "password123",
	})
	if resp.StatusCode != 409 {
		t.Fatalf("POST /users before bootstrap = %d, want 409", resp.StatusCode)
	}
}

func TestCreateUser_AfterBootstrap(t *testing.T) {
	srv, db := setupFreshInstanceServer(t)
	suffix := testpg.UniqueSuffix(t, db)

	resp, _ := doInstancePost(t, srv, "/instance/bootstrap", map[string]string{
		"email":    suffix + "@test.local",
		"name":     "Admin",
		"password": "securepass123",
	})
	if resp.StatusCode != 201 {
		t.Fatalf("bootstrap = %d, want 201", resp.StatusCode)
	}

	resp2, _ := doInstancePost(t, srv, "/users", map[string]string{
		"email":    "user" + suffix + "@test.local",
		"name":     "User",
		"password": "password123",
	})
	if resp2.StatusCode != 201 {
		t.Fatalf("POST /users after bootstrap = %d, want 201", resp2.StatusCode)
	}
}
