package cmd

import "testing"

func TestExpandEscapedNewlines(t *testing.T) {
	input := `server_address: http://localhost:8181\nrealm: default-realm\r\nclient_id: root`
	got := expandEscapedNewlines(input)
	want := "server_address: http://localhost:8181\nrealm: default-realm\nclient_id: root"

	if got != want {
		t.Fatalf("expandEscapedNewlines() = %q, want %q", got, want)
	}
}
