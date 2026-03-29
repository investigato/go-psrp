package auth

import (
	"net/http"

	ntlmSsp "github.com/investigato/ntlmssp"
	ntlmHttp "github.com/investigato/ntlmssp/http"
)

// NTLMAuth implements NTLM authentication with optional Channel Binding Token (CBT) support.
type NTLMAuth struct {
	creds     Credentials
	enableCBT bool
}

// NTLMAuthOption configures NTLM authentication.
type NTLMAuthOption func(*NTLMAuth)

// WithCBT enables Channel Binding Tokens for Extended Protection.
// When enabled, the NTLM authentication will include a CBT derived from
// the TLS server certificate, protecting against NTLM relay attacks.
func WithCBT(enable bool) NTLMAuthOption {
	return func(a *NTLMAuth) {
		a.enableCBT = enable
	}
}

// NewNTLMAuth creates a new NTLM authentication handler.
// By default, CBT is disabled for backwards compatibility.
// Use WithCBT(true) to enable Extended Protection.
func NewNTLMAuth(creds Credentials, opts ...NTLMAuthOption) *NTLMAuth {
	a := &NTLMAuth{creds: creds}
	for _, opt := range opts {
		opt(a)
	}
	return a
}

// Name returns the authentication scheme name.
func (a *NTLMAuth) Name() string {
	return "NTLM"
}

// Transport wraps an http.RoundTripper with NTLM authentication.
// Uses github.com/Azure/go-ntlmssp for the NTLM handshake logic (connection management),
// and a custom injector to add CBT support if enabled.
func (a *NTLMAuth) Transport(base http.RoundTripper) (http.RoundTripper, error) {
	ntlmClient, err := ntlmSsp.NewClient(
		ntlmSsp.SetUserInfo(a.creds.Username, a.creds.Password),
		ntlmSsp.SetDomain(a.creds.Domain),
	)
	if err != nil {
		return nil, err
	}
	ntlmHTTPClient, err := ntlmHttp.NewClient(
		&http.Client{Transport: base},
		ntlmClient,
		ntlmHttp.SendCBT(a.enableCBT),
		ntlmHttp.Encryption(false),
	)
	if err != nil {
		return nil, err
	}
	return ntlmHTTPClient, nil
}
