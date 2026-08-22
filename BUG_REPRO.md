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
--- FAIL: TestBusinessChain06 (0.02s)
    business_chain_test.go:35: file A52-F-6 status disagrees with archived case
FAIL
FAIL	lawdrive	0.035s
?   	lawdrive/cmd/casedrive	[no test files]
?   	lawdrive/internal/domain	[no test files]
ok  	lawdrive/internal/audit	0.016s
ok  	lawdrive/internal/files	0.017s
ok  	lawdrive/internal/httpapi	0.013s
ok  	lawdrive/internal/permissions	0.013s
ok  	lawdrive/internal/query	0.014s
ok  	lawdrive/internal/share	0.013s
ok  	lawdrive/internal/store	0.019s
ok  	lawdrive/internal/workflow	0.022s
FAIL
```

## Architecture reproduction

### linux/amd64
- Go toolchain version: exit `0`
- Go build (.): exit `0`
- Go test (.): exit `1`
- Go run smoke (cmd/casedrive): exit `0`
### linux/arm64
- Go toolchain version: exit `0`
- Go build (.): exit `0`
- Go test (.): exit `1`
- Go run smoke (cmd/casedrive): exit `0`
