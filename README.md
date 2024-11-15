## Getting Started

### Installation

1. Clone the repository

```bash
git clone https://github.com/yourusername/tx-parser.git
cd tx-parser
```

2. Install dependencies

```bash
go mod download
```

3. Generate mocks (required for testing)

```bash
make mock
```

### Running the Application

1. Start the server in development mode:

```bash
make run
```

The server will start on `http://localhost:5005`

## Project structure

### Folder structure

```
tx-parser/
├── cmd/
│   └── server                # Application entry point
├── internal/
│   ├── api                   # HTTP API routes
│   ├── models                # Transaction model definitions
│   ├── parser                # Include Parser interface and Ethereum parser implementation
│   └── storage               # Storage interface
│       └── memory_storage.go # In-memory storage implementation
├── pkg/
│   ├── logger                # Logging utilities
│   ├── notification          # Notification interface to communicate with notification service
│   │   └── console.go        # Console notifier implementation - print any notify to console
│   └── rpc                   # RPC client interface and Ethereum RPC client
├── mocks                     # Generated mock files, ignored by git, need run `make mock`
├── Makefile                  # Build and development commands
├── go.mod                    # Go module file
├── go.sum                    # Go dependencies checksums
└── README.md                 # Project documentation
```

### Key Components

1. Command Layer (cmd/)

- Entry point for the application
- Server initialization and configuration

2. Internal Package (internal/)

- API: HTTP router and handlers
- Models: contains transaction data structure shared by parser and the api
- Parser: Core business logic for parsing Ethereum blocks. If no current block (current block is 0) we'll process from current latest block fetched from the RPC.
- Storage: Data persistence layer to support the Parser. Currently we have in-memory storage, evertime the server re-start data will be wipe-out.

3. Package Layer (pkg/)

- Logger: Logging utilities
- Notification: Notification utilities to communicate with notification service
- RPC: Ethereum JSON-RPC client

4. Build Tools:
   A makefile with commands for:

- Building: `make build`
- Develop: `make run`
- Testing: `make test`
- Mock generation: `make mock`
- Code coverage: `make coverage`

### Diagrams

- Component diagrams

```mermaid
graph TB
    subgraph External
        ETH[Ethereum Node]
    end

    subgraph Txn Parser Service
        API[API Layer]
        Parser[ETH Parser]
        Storage[In-Memory Storage]
        RPC[RPC Client]
        Notifier[Console Notifier]
    end

    Client[API Client] --> |REST API| API
    RPC --> |JSON-RPC| ETH

    API --> |Subscribe/Query| Parser
    Parser --> |Read/Write| Storage
    Parser --> |Process Blocks| RPC
    Parser --> |Notify| Notifier
```

- Parser flow diagrams

```mermaid
sequenceDiagram
    participant BG as Background Process
    participant Parser as ETH Parser
    participant RPC as RPC Client
    participant Storage as Memory Storage
    participant Notifier as Console Notifier

    loop Every 15 seconds
        BG->>Parser: Trigger block processing
        Parser->>RPC: Get latest block number
        RPC-->>Parser: Return current block

        Parser->>Storage: Get last processed block
        Storage-->>Parser: Return block number

        alt New blocks available
            loop For each new block
                Parser->>RPC: Get block by number
                RPC-->>Parser: Return block data

                Parser->>Storage: Get subscribers
                Storage-->>Parser: Return subscriber addresses

                loop For each transaction
                    Parser->>Parser: Match transaction addresses

                    alt Address matches subscriber
                        Parser->>Storage: Save transaction
                        Parser->>Notifier: Send notification
                        Notifier-->>Parser: Notification sent
                    end
                end

                Parser->>Storage: Update current block
                Storage-->>Parser: Confirm update
            end
        end
    end
```

### API documentation

#### Get current processed block

```bash
curl -X GET 'http://localhost:5005/api/v1/block/current'
```

#### Subscribe an address for notification

```bash
curl -X POST 'http://localhost:5005/api/v1/subscribe' \
-H 'Content-Type: application/json' \
-d '{
    "address": "0x3fC91A3afd70395Cd496C647d5a6CC9D4B2b7FAD"
}'
```

#### Get transactions for an address

```bash
curl -X GET 'http://localhost:5005/api/v1/transactions?address=0x3fC91A3afd70395Cd496C647d5a6CC9D4B2b7FAD'
```
