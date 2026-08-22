# BUG_REPRO

The following failures were observed while validating the initial project state.
Each section records what failed, how to reproduce it, and the complete command output.
They are preserved intentionally; only failing build gates are omitted from the generated Dockerfile.

## Failure 1: Go test (.)

- Observed problem: `Go test (.)` failed in the initial project state.
- Working directory: `.`
- Command: `cd /app && GOTOOLCHAIN=local GOPROXY=off GOSUMDB=off go test -count=1 ./...`
- Exit status: `1`

```text
?   	coldrisk.local/console/cmd/riskctl	[no test files]
ok  	coldrisk.local/console/internal/domain	0.001s
ok  	coldrisk.local/console/internal/httpapi	0.007s
ok  	coldrisk.local/console/internal/policy	0.001s
ok  	coldrisk.local/console/internal/report	0.001s
--- FAIL: TestBusiness01Regression (0.01s)
    business_regression_test.go:30: regional scope leaked records: []domain.Record{domain.Record{ID:"r-a", StoreID:"store-a", Title:"Authorized repair", Status:"open", Severity:"high", Description:"", PhotoRefs:[]string(nil), Assignee:"", CreatedAt:"t", UpdatedAt:"t"}, domain.Record{ID:"r-b", StoreID:"store-b", Title:"Other repair", Status:"open", Severity:"high", Description:"", PhotoRefs:[]string(nil), Assignee:"", CreatedAt:"t", UpdatedAt:"t"}}
FAIL
FAIL	coldrisk.local/console/internal/service	0.014s
ok  	coldrisk.local/console/internal/store	0.010s
FAIL
```

## Architecture reproduction

### linux/amd64
- Go toolchain version: exit `0`
- Node.js version: exit `0`
- Go build (.): exit `0`
- Go test (.): exit `1`
- Go run smoke (cmd/riskctl): exit `0`
- Frontend build (web): exit `0`
### linux/arm64
- Go toolchain version: exit `0`
- Node.js version: exit `0`
- Go build (.): exit `0`
- Go test (.): exit `1`
- Go run smoke (cmd/riskctl): exit `0`
- Frontend build (web): exit `0`
