# gstat

> Talks directly to the Linux kernel.

`gstat` reads your system's vitals straight from `/proc` — the virtual filesystem the kernel exposes to userspace — and renders them in your terminal. CPU usage calculated from raw scheduler ticks. Memory pulled from `MemInfo`. Processes ranked by CPU time accumulated since boot.

No dependencies. Just Go and the kernel.

## Install

```bash
git clone https://github.com/saadhtiwana/gstat
cd gstat
go build -o gstat .
```

## Usage

```bash
# one-time report
./gstat snapshot

# live dashboard, refreshes every 3s
./gstat monitor

# custom refresh rate
./gstat --interval 5 monitor

# dump everything to JSON
./gstat export
```

## How it works

Most system monitor tools shell out to `top` or use heavyweight libraries like `psutil`. `gstat` doesn't. It reads the kernel's own data directly:

| What | Where |
|------|-------|
| CPU usage | `/proc/stat` — sampled twice, delta calculated |
| CPU model | `/proc/cpuinfo` |
| Memory | `/proc/meminfo` |
| Disk | `df -B1` |
| Network I/O | `/proc/net/dev` |
| Processes | `/proc/<pid>/stat` + `/proc/<pid>/status` |

CPU usage is calculated the same way `htop` does it — read idle and total jiffies, sleep 250ms, read again, compute the delta.

## Project structure

```
gstat/
├── main.go          — entry point, flag parsing, subcommand routing
└── internal/
    ├── cpu/         — usage sampling + model name
    ├── mem/         — total, used, available
    ├── disk/        — per-mount usage
    ├── netstat/     — per-interface bytes sent/received
    ├── process/     — top 10 by CPU ticks
    └── render/      — terminal output, color coding, JSON export
```

Each package exposes exactly one function: `Get()`. The `render` package is the only one that knows how to print anything. Clean separation, no circular imports.

## Tests

```bash
go test ./...
```

## Built with

Go 1.22 · Linux · `/proc`
