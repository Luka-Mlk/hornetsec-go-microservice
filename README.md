## Quick Start

The fastest way to get the service running is using **Docker Compose**:

```bash
docker-compose up --build
```

* **REST API:** localhost:8080
* **gRPC Server:** localhost:50051

## Prerequisites

* **Docker** (Recommended for running in container)
* **Go 1.25+** (Local development and testing)

## API Usage

### REST Endpoints
1. `POST /api/v1/documents` - Create a new document
2. `GET /api/v1/documents` - List all documents
3. `GET /api/v1/documents/{id}` - Retrieve a specific document by UUID
4. `DELETE /api/v1/documents/{id}` - Remove a document by UUID

### gRPC Service Methods
The service definition can be found in the `.proto` files.
* `document.DocumentService/Create`
* `document.DocumentService/List`
* `document.DocumentService/Get`
* `document.DocumentService/Delete`

### Usage Examples

**Create a Document**
```bash
curl -X POST http://localhost:8080/api/v1/documents \
     -H "Content-Type: application/json" \
     -d '{"name": "Architecture Plan", "description": "Draft for the new microservice layout"}'
```

**Retrieve by ID**
```bash
curl -X GET http://localhost:8080/api/v1/documents/{id}
```

**Delete by ID**
```bash
curl -X DELETE http://localhost:8080/api/v1/documents/{id}
```

## Testing

**Run all tests:**
```bash
go test ./...
```

## Configuration

Services are configured via environment variables (defined in `docker-compose.yaml` or a `.env` file):

| Variable    | Description                             |
| :---------- | :-------------------------------------- |
| `HTTP_PORT` | Port for the REST API (e.g., :8080)     |
| `GRPC_PORT` | Port for the gRPC server (e.g., :50051) |
