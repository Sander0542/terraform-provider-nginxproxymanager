// SPDX-License-Identifier: MIT

package provider

import (
	"encoding/json"
	"errors"
	"strings"
	"testing"

	"github.com/sander0542/nginxproxymanager-go"
)

func TestNormalizeProxyHostAccessListResponse(t *testing.T) {
	response := []byte(`{"id":1,"created_on":"2026-01-01 00:00:00","modified_on":"2026-01-01 00:00:00","owner_user_id":1,"domain_names":["evcc.example"],"forward_host":"127.0.0.1","forward_port":80,"access_list_id":0,"certificate_id":null,"ssl_forced":false,"caching_enabled":false,"block_exploits":false,"advanced_config":"","meta":{},"allow_websocket_upgrade":false,"http2_support":false,"forward_scheme":"http","enabled":true,"locations":[],"hsts_enabled":false,"hsts_subdomains":false,"access_list":false}`)

	if err := json.Unmarshal(response, new(nginxproxymanager.GetProxyHosts200ResponseInner)); err == nil {
		t.Fatal("expected SDK decoding to reject access_list=false")
	}

	normalized, changed := normalizeProxyHostAccessListResponse(response)
	if !changed {
		t.Fatal("expected false access_list value to be normalized")
	}

	var host nginxproxymanager.GetProxyHosts200ResponseInner
	if err := json.Unmarshal(normalized, &host); err != nil {
		t.Fatalf("SDK could not decode normalized proxy-host response: %v", err)
	}

	expandedAccessList := `{"id":10,"created_on":"2026-01-01 00:00:00","modified_on":"2026-01-01 00:00:00","owner_user_id":1,"is_deleted":false,"name":"RFC1918 with Auth","meta":{},"satisfy_any":true,"pass_auth":false,"clients":[{"id":28,"created_on":"2026-01-01 00:00:00","modified_on":"2026-01-01 00:00:00","access_list_id":10,"address":"10.0.0.0/8","directive":"allow","meta":{}}],"items":[{"id":10,"created_on":"2026-01-01 00:00:00","modified_on":"2026-01-01 00:00:00","access_list_id":10,"username":"username","password":"password","meta":{}}]}`
	expandedResponse := []byte(strings.Replace(string(response), `"access_list":false`, `"access_list":`+expandedAccessList, 1))
	if err := json.Unmarshal(expandedResponse, new(nginxproxymanager.GetProxyHosts200ResponseInner)); err == nil {
		t.Fatal("expected SDK decoding to reject expanded access_list object")
	}

	normalized, changed = normalizeProxyHostAccessListResponse(expandedResponse)
	if !changed {
		t.Fatal("expected expanded access_list object to be normalized")
	}
	if err := json.Unmarshal(normalized, &host); err != nil {
		t.Fatalf("SDK could not decode normalized proxy-host response with expanded access_list: %v", err)
	}

	unchanged := []byte(`{"access_list":false}`)
	got, changed := normalizeProxyHostAccessListResponse(unchanged)
	if changed || string(got) != string(unchanged) {
		t.Fatalf("non-proxy-host object was modified: %s", got)
	}
}

type testAPIError struct {
	body []byte
}

func (e testAPIError) Error() string {
	return "400 Bad Request"
}

func (e testAPIError) Body() []byte {
	return e.body
}

func TestErrorWithResponseBody(t *testing.T) {
	err := testAPIError{body: []byte(`{"error":"invalid forward host"}`)}
	got := errorWithResponseBody(err)
	want := `400 Bad Request: {"error":"invalid forward host"}`
	if got != want {
		t.Fatalf("errorWithResponseBody() = %q, want %q", got, want)
	}

	ordinaryError := errors.New("connection refused")
	if got := errorWithResponseBody(ordinaryError); got != ordinaryError.Error() {
		t.Fatalf("ordinary error changed: got %q", got)
	}
}