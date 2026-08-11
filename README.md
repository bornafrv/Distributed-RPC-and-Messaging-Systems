# Distributed RPC and Messaging Systems

A multi-machine distributed-system study implemented in Go. It separates authentication and file storage into independent services, communicates through JSON-RPC over TCP, and extends the architecture with publish/subscribe monitoring for high-memory alerts.

## Architecture

```text
┌──────────────┐      JSON-RPC      ┌──────────────┐
│ Web service  │ ─────────────────▶ │ Auth service │
│    (VM1)     │                    │    (VM2)     │
│              │      JSON-RPC      └──────────────┘
│              │ ─────────────────▶ ┌──────────────┐
└──────┬───────┘                    │ File service │
       │ memory alerts              │    (VM3)     │
       ▼                            └──────────────┘
┌──────────────┐
│  Subscriber  │
└──────────────┘
```

The repository contains the backend components supplied for the distributed deployment. The web-facing VM is treated as a separate consumer of these services.

## Authentication service

`auth-vm/` isolates credential storage and authentication from the web application.

- Runs as an independent JSON-RPC server.
- Uses Go's `net/rpc` package with the `jsonrpc` codec.
- Accepts remote authentication requests.
- Validates users stored locally in `users.json`.
- Prevents the web tier from reading credential files directly.

The service README contains the exact request/response types and network configuration.

## File service

`file-vm/` stores and serves files from a separate machine.

- Lists remotely available files.
- Returns requested file content over RPC.
- Keeps storage physically separate from the web service.
- Includes small text and image samples for testing.

This design demonstrates location transparency: the web service consumes remote resources without relying on local file paths.

## Publish/subscribe monitoring

`auth-vm/subscriber/` demonstrates event-driven communication in addition to synchronous RPC.

- The web service periodically monitors memory usage.
- A high-memory threshold triggers an alert event.
- A subscriber listens for and displays published alerts.
- The extension improves observability and illustrates a simple Pub/Sub architecture.

## Repository structure

```text
.
├── auth-vm/
│   ├── main.go
│   ├── users.json
│   ├── subscriber/
│   └── README.md
├── file-vm/
│   ├── main.go
│   ├── files/
│   ├── images/
│   └── README.md
├── rpc-study/
│   └── rpc_summary.pdf
├── assignment.pdf
├── .gitignore
└── README.md
```

## Requirements

- Go
- Three Linux machines or virtual machines for the intended deployment
- Reachable TCP ports between participating machines

The services rely primarily on Go's standard networking and RPC packages.

## Running the components

1. Configure each service address for the target VM network.
2. Start the authentication service on VM2.
3. Start the file service on VM3.
4. Start the subscriber before testing high-memory alerts.
5. Run the web-tier client on VM1 and verify remote authentication and file access.

See the component README files for commands, ports, RPC method descriptions, and testing examples.

## Security note

`users.json` contains only small coursework test records. A production design should use hashed passwords, encrypted transport, authenticated RPC, secret management, and a persistent database.

## Technologies

`Go` · `JSON-RPC` · `TCP` · `Virtual Machines` · `Pub/Sub` · `Distributed Services` · `Monitoring`

