Malice Engine API
=================

The Engine API is an HTTP API used by the command-line client to communicate with the daemon. It can also be used by third-party software to control the daemon.

It consists of various components in this repository:

- `api/apiary.apib` A API Blueprint definition of the API.
- `api/types/` Types shared by both the client and server, representing various objects, options, responses, etc.
- `cmd/` The command-line client.
- `client/` The Go client used by the command-line client. It can also be used by third-party Go programs.
- `daemon/` The daemon, which serves the API.