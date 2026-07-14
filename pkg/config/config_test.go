package config

import "testing"

func TestParseConfigValuesJSON(t *testing.T) {
	values, err := ParseConfigValues([]byte(`{
		"polaris_server_address": "http://localhost:8181",
		"realm": "default-realm",
		"root_client_id": "root",
		"root_client_secret": "secret"
	}`))
	if err != nil {
		t.Fatalf("ParseConfigValues returned error: %v", err)
	}

	cfg := &Config{}
	if err := ApplyConfigValues(cfg, values); err != nil {
		t.Fatalf("ApplyConfigValues returned error: %v", err)
	}

	if cfg.Host != "http://localhost:8181" {
		t.Fatalf("Host = %q, want http://localhost:8181", cfg.Host)
	}
	if cfg.Realm != "default-realm" {
		t.Fatalf("Realm = %q, want default-realm", cfg.Realm)
	}
	if cfg.RootClientID != "root" {
		t.Fatalf("RootClientID = %q, want root", cfg.RootClientID)
	}
	if cfg.RootClientSecret != "secret" {
		t.Fatalf("RootClientSecret = %q, want secret", cfg.RootClientSecret)
	}
}

func TestParseConfigValuesYAML(t *testing.T) {
	values, err := ParseConfigValues([]byte(`
server_address: http://localhost:8181
realm: default-realm
client_id: root
client_secret: "secret:value"
`))
	if err != nil {
		t.Fatalf("ParseConfigValues returned error: %v", err)
	}

	cfg := &Config{}
	if err := ApplyConfigValues(cfg, values); err != nil {
		t.Fatalf("ApplyConfigValues returned error: %v", err)
	}

	if cfg.Host != "http://localhost:8181" {
		t.Fatalf("Host = %q, want http://localhost:8181", cfg.Host)
	}
	if cfg.Realm != "default-realm" {
		t.Fatalf("Realm = %q, want default-realm", cfg.Realm)
	}
	if cfg.RootClientID != "root" {
		t.Fatalf("RootClientID = %q, want root", cfg.RootClientID)
	}
	if cfg.RootClientSecret != "secret:value" {
		t.Fatalf("RootClientSecret = %q, want secret:value", cfg.RootClientSecret)
	}
}
