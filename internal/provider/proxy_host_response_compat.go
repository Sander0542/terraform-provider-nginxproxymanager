// SPDX-License-Identifier: MIT

package provider

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"strings"
)

const proxyHostsAPIPath = "/api/nginx/proxy-hosts"

type proxyHostAccessListTransport struct {
	base http.RoundTripper
}

func (t proxyHostAccessListTransport) RoundTrip(request *http.Request) (*http.Response, error) {
	base := t.base
	if base == nil {
		base = http.DefaultTransport
	}

	response, err := base.RoundTrip(request)
	if err != nil || response == nil || response.Body == nil || !isProxyHostsAPIPath(request.URL.Path) {
		return response, err
	}

	body, err := io.ReadAll(response.Body)
	response.Body.Close()
	if err != nil {
		return nil, err
	}

	normalized, changed := normalizeProxyHostAccessListResponse(body)
	if !changed {
		response.Body = io.NopCloser(bytes.NewReader(body))
		return response, nil
	}

	response.Body = io.NopCloser(bytes.NewReader(normalized))
	response.ContentLength = int64(len(normalized))
	response.Header.Set("Content-Length", strconv.Itoa(len(normalized)))
	return response, nil
}

func isProxyHostsAPIPath(path string) bool {
	path = strings.TrimRight(path, "/")
	return path == proxyHostsAPIPath || strings.HasPrefix(path, proxyHostsAPIPath+"/")
}

func normalizeProxyHostAccessListResponse(body []byte) ([]byte, bool) {
	trimmed := bytes.TrimSpace(body)
	if len(trimmed) == 0 {
		return body, false
	}

	if trimmed[0] != '[' {
		return normalizeProxyHostObject(trimmed)
	}

	var hosts []json.RawMessage
	if err := json.Unmarshal(trimmed, &hosts); err != nil {
		return body, false
	}

	changed := false
	for index, host := range hosts {
		if normalized, hostChanged := normalizeProxyHostObject(host); hostChanged {
			hosts[index] = normalized
			changed = true
		}
	}
	if !changed {
		return body, false
	}

	normalized, err := json.Marshal(hosts)
	if err != nil {
		return body, false
	}
	return normalized, true
}

func normalizeProxyHostObject(body []byte) ([]byte, bool) {
	var host map[string]json.RawMessage
	if err := json.Unmarshal(body, &host); err != nil {
		return body, false
	}

	accessList, hasAccessList := host["access_list"]
	_, hasAccessListID := host["access_list_id"]
	_, hasDomainNames := host["domain_names"]
	accessList = bytes.TrimSpace(accessList)
	if !hasAccessList || !hasAccessListID || !hasDomainNames || bytes.Equal(accessList, []byte("null")) {
		return body, false
	}
	if !bytes.Equal(accessList, []byte("false")) && (len(accessList) == 0 || accessList[0] != '{') {
		return body, false
	}

	host["access_list"] = json.RawMessage("null")
	normalized, err := json.Marshal(host)
	if err != nil {
		return body, false
	}
	return normalized, true
}