# go-psrp - project mgr notes

High-level PSRP/WinRM client library. Sits on top of go-psrpcore (protocol) and ntlmssp (auth).

## Package Layout

client/      — Public API: New(), Connect(), Execute(), Close()
powershell/  — Backend bridge: WSManBackend, HvSocketBackend, WSManTransport
wsman/       — WSMan protocol: Client, SOAP envelopes, EPR
wsman/auth/  — Auth handlers: BasicAuth, NTLMAuth, NegotiateAuth, KerberosProvider
wsman/transport/ — HTTP transport: HTTPTransport
winrs/       — WinRS (cmd.exe) shell, separate from PSRP

## Key Types

**wsman.Client** — WSMan protocol operations:

- `Create(ctx, options, creationXML)` → `*EndpointReference`
- `Command(ctx, epr, commandID, args)` — create pipeline slot
- `Send(ctx, epr, commandID, stream, data)`
- `Receive(ctx, epr, commandID)` → `*ReceiveResult`
- `Signal(ctx, epr, commandID, code)`
- `Delete(ctx, epr)`

**EndpointReference (EPR)** — shell identifier; holds ResourceURI + Selectors (ShellId). Selectors use bare UUID — no `uuid:` prefix.

**WSManBackend** (`powershell/`) — implements `RunspaceBackend`:

- `Init(ctx, pool)` — creates shell with creationXml, then calls `pool.ReceiveHandshakeWSMan(ctx)` (receive-only, no sends)
- `Transport()` → `*WSManTransport` (io.ReadWriter for the pool)
- `Configure(client, epr, commandID)` — sets active shell/command targets

**WSManTransport** (`powershell/transport.go`):

- `Write(p)` — calls `client.Send(ctx, epr, commandID, "stdin pr", p)`
- `Read(p)` — calls `client.Receive()` in loop; continues on TimedOut

## Auth

`wsman/auth/ntlm.go` — `NTLMAuth` wraps ntlmssp http.Client. Encryption is on. PassTheHash available.

`wsman/auth/kerberos.go` — `KerberosProvider` wraps lib/krb5 via spnego.

## PSRP Shell Format (PowerShell URI)

`Create()` always uses PSRP format for `ResourceURIPowerShell`:

- Streams: `stdin pr` / `stdout`
- `<creationXml>` contains base64-encoded SESSION_CAPABILITY + INIT_RUNSPACEPOOL fragments

WinRS format (`stdin`/`stdout stderr`) only used in `default:` branch.

## PSRP Command Flow (ground truth)

CreateShell (with creationXml)  ← fragments embedded here
Receive                          ← SESSION_CAPABILITY + APPLICATION_PRIVATE
Receive                          ← RUNSPACEPOOL_STATE=Opened
Command (with CREATE_PIPELINE)   ← pipeline definition embedded in Command body
Receive                          ← pipeline output

No Send after Command. CREATE_PIPELINE goes in the Command request.
