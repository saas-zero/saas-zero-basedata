package middleware

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// TestOperationLog_MaskSensitive 验证登录/改密等敏感字段被脱敏。
func TestOperationLog_MaskSensitive(t *testing.T) {
	in := `{"username":"admin","password":"secret123","nickname":"管理员","mobile":"13800138000"}`
	out := maskSensitive(in)
	// 敏感字段必须脱敏
	for _, k := range []string{"password"} {
		if contains(out, `"`+k+`":"***"`) {
			continue
		}
		t.Fatalf("expected %s to be masked in output: %s", k, out)
	}
	if contains(out, "secret123") {
		t.Fatalf("password value must not leak: %s", out)
	}
	// 普通字段保留，手机号等敏感字段必须脱敏
	if !contains(out, "admin") {
		t.Fatalf("non-sensitive fields should be preserved: %s", out)
	}
	if contains(out, "13800138000") {
		t.Fatalf("mobile value must not leak: %s", out)
	}
}

func TestOperationLog_MaskSensitive_NonJSON(t *testing.T) {
	raw := "not-json-at-all"
	if got := maskSensitive(raw); got != "[non-json payload omitted]" {
		t.Fatalf("non-JSON input must be omitted, got %s", got)
	}
}

func TestOperationLog_MaskSensitive_Deep(t *testing.T) {
	in := `{"user":{"newPassword":"abc123"},"captchaVal":"1234","data":{"token":"jwt.abc"}}`
	out := maskSensitive(in)
	for _, leaked := range []string{"abc123", "1234", "jwt.abc"} {
		if contains(out, leaked) {
			t.Fatalf("sensitive value %q must be masked: %s", leaked, out)
		}
	}
}

func TestOperationLog_ExtractErrorMsg(t *testing.T) {
	body := `{"code":400,"msg":"参数错误"}`
	if got := extractErrorMsg(body); got != "参数错误" {
		t.Fatalf("expected 参数错误, got %s", got)
	}
	if got := extractErrorMsg(`{"code":200,"msg":"success"}`); got != "" {
		t.Fatalf("successful response must not have an error message, got %s", got)
	}
	if got := extractErrorMsg(""); got != "" {
		t.Fatalf("expected empty for empty body, got %s", got)
	}
	if got := extractErrorMsg("not-json"); got != "" {
		t.Fatalf("expected empty for non-json body, got %s", got)
	}
}

func TestOperationLog_ResponseStatus(t *testing.T) {
	cases := []struct {
		httpStatus int
		body       string
		want       string
	}{
		{http.StatusOK, `{"code":200,"msg":"success"}`, "success"},
		{http.StatusOK, `{"code":400,"msg":"参数错误"}`, "failure"},
		{http.StatusBadRequest, `{"code":400,"msg":"参数错误"}`, "failure"},
		{http.StatusOK, `not-json`, "success"},
	}
	for _, tc := range cases {
		if got := responseStatus(tc.httpStatus, tc.body); got != tc.want {
			t.Fatalf("responseStatus(%d, %q) = %q, want %q", tc.httpStatus, tc.body, got, tc.want)
		}
	}
}

func TestOperationLog_LargePayloadNeverKeepsRawSensitiveValue(t *testing.T) {
	payload := `{"password":"secret"}` + strings.Repeat("x", 5000)
	got := safeLogPayload(payload, 4000)
	if strings.Contains(got, "secret") {
		t.Fatalf("password leaked from oversized payload: %s", got)
	}
}

func TestOperationLog_ParseModuleOperation(t *testing.T) {
	cases := []struct {
		path, module, op string
	}{
		{"/system/user/create", "user", "create"},
		{"/system/dict/data/list", "dict", "list"},
		{"/oauth/login", "oauth", "login"},
		{"", "", ""},
	}
	for _, c := range cases {
		m, o := parseModuleOperation(c.path)
		if m != c.module || o != c.op {
			t.Fatalf("parseModuleOperation(%q) = (%q, %q), want (%q, %q)", c.path, m, o, c.module, c.op)
		}
	}
}

// TestCasbinAuth_FailClosed_NilEnforcer 验证 enforcer 为 nil 且未显式禁用时 fail-closed（503）。
func TestCasbinAuth_FailClosed_NilEnforcer(t *testing.T) {
	mw := CasbinAuth(nil, false)
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Fatal("next must not be called when enforcer is nil and not disabled")
	})
	h := mw(next)
	req := httptest.NewRequest(http.MethodGet, "/system/user/list", nil)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)
	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("expected 503/500 for nil enforcer fail-closed, got %d", rec.Code)
	}
}

// TestCasbinAuth_Disabled_PassThrough 验证显式 casbinDisabled 时放行（本地开发）。
func TestCasbinAuth_Disabled_PassThrough(t *testing.T) {
	mw := CasbinAuth(nil, true)
	called := false
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
	})
	h := mw(next)
	req := httptest.NewRequest(http.MethodGet, "/system/user/list", nil)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)
	if !called {
		t.Fatal("expected next to be called when casbin disabled")
	}
}

func contains(s, sub string) bool {
	for i := 0; i+len(sub) <= len(s); i++ {
		if s[i:i+len(sub)] == sub {
			return true
		}
	}
	return false
}
