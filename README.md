# Person Service

A serverless CRUD API for managing persons, built with AWS CDK, API Gateway, DynamoDB, and Go Lambda functions.

## Architecture

```
API Gateway
  ├── GET  /persons  → getLambda
  └── POST /persons  → postLambda
                          ↓
                     DynamoDB (persons table)
                          ↓ (DynamoDB Stream)
                     streamLambda
                          ↓
                     EventBridge (person-bus)
                          ↓ (rule: PersonCreated)
                     CloudWatch Log Group
```

## Project Structure

```
person-service/
├── lib/
│   └── person-service-stack.ts   # CDK stack — all AWS infrastructure
├── service/                      # GET/POST Lambda (Go)
│   ├── main.go                   # Entry point, dependency wiring
│   ├── domain/
│   │   └── person.go             # Person model
│   ├── ports/
│   │   └── ports.go              # Interfaces (DynamoDBClient, PersonRepository, PersonService)
│   ├── repository/
│   │   └── person_repository.go  # DynamoDB adapter
│   ├── service/
│   │   └── person_service.go     # Business logic
│   └── handler/
│       └── person_handler.go     # Lambda event parsing and routing
└── stream-service/               # DynamoDB Stream Lambda (Go)
    └── main.go                   # Publishes PersonCreated events to EventBridge
```

## Service Layer Architecture

The `service/` Lambda follows a clean, layered architecture (Ports and Adapters):

```
main.go (composition root)
  └── handler  (parses Lambda event, routes to service)
        └── service  (business logic, UUID generation)
              └── repository  (DynamoDB access only)
                    └── domain  (Person struct, no dependencies)
```

Each layer depends only on interfaces defined in `ports/`, making every layer independently testable.

## Live API

A deployed instance is available for testing:

**Base URL:** `https://956hpe9ne4.execute-api.eu-west-1.amazonaws.com/prod`

### List Persons
```bash
curl https://956hpe9ne4.execute-api.eu-west-1.amazonaws.com/prod/persons
```

### Create Person
```bash
curl -X POST https://956hpe9ne4.execute-api.eu-west-1.amazonaws.com/prod/persons \
  -H "Content-Type: application/json" \
  -d '{
    "firstName": "Doruk",
    "lastName": "Demirci",
    "phoneNumber": "01233445",
    "address": "Amsterdam Nederland"
  }'
```

## API

### POST /persons
Creates a new person. ID is auto-generated.

**Request body:**
```json
{
  "firstName": "John",
  "lastName": "Doe",
  "phoneNumber": "1234567890",
  "address": "Amsterdam"
}
```

**Response:**
```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "firstName": "John",
  "lastName": "Doe",
  "phoneNumber": "1234567890",
  "address": "Amsterdam"
}
```

### GET /persons
Returns all persons.

**Response:**
```json
[
  {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "firstName": "John",
    "lastName": "Doe",
    "phoneNumber": "1234567890",
    "address": "Amsterdam"
  }
]
```

## Stream Service

The `stream-service/` Lambda is triggered automatically by DynamoDB Streams whenever a record changes in the `persons` table.

**How it works:**

1. A new person is created via `POST /persons` → written to DynamoDB
2. DynamoDB Stream emits a record with `eventName: INSERT` and the `NewImage` of the item
3. `stream-service` reads the record, ignores non-INSERT events (UPDATE, REMOVE)
4. Publishes a `PersonCreated` event to the custom EventBridge bus
5. An EventBridge rule matches `PersonCreated` and forwards it to CloudWatch Logs

**Event published to EventBridge:**
```json
{
  "source": "person.service",
  "detail-type": "PersonCreated",
  "detail": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "firstName": "John",
    "lastName": "Doe",
    "phoneNumber": "1234567890",
    "address": "Amsterdam"
  }
}
```

**Reliability:**
- Retry attempts: 3
- Failed records after retries are sent to an SQS dead-letter queue (`person-dlq-{env}`) with 14-day retention

## Running Tests

```bash
cd service && go test ./...
```

## Testing Approach

Tests use hand-written mocks with function fields instead of a mock library such as `testify/mock`.

```go
type mockPersonService struct {
    createPerson func(ctx context.Context, person domain.Person) (domain.Person, error)
    listPersons  func(ctx context.Context) ([]domain.Person, error)
}
```

**Why not testify/mock?**

`testify/mock` adds value when you need to verify how many times a method was called or with what exact arguments — useful in complex business logic with side effects. For this service:

- Each layer has only 1-2 methods
- The business logic is straightforward with no conditional branching based on call count
- Function fields already let each test define its own behavior inline without shared state
- Zero external dependencies keeps the test binary lightweight

If the service grows significantly in complexity, adopting `testify/mock` would be a reasonable next step.

## Deployment

The Lambdas use the `provided.al2` runtime (Amazon Linux 2). AWS deprecated the native `go1.x` runtime in 2023, and `provided.al2` is the recommended replacement. It requires you to compile a binary named `bootstrap` and upload it directly — Go is no longer managed by the Lambda runtime itself. This gives you full control over the Go version and produces a smaller, faster cold start compared to `go1.x`.

Build the Go Lambdas for Linux (required for `provided.al2` runtime):

```bash
# Build service Lambda
cd service && GOOS=linux GOARCH=amd64 go build -o bootstrap

# Build stream-service Lambda
cd stream-service && GOOS=linux GOARCH=amd64 go build -o bootstrap
```

Deploy to AWS:

```bash
cdk deploy -c env=dev
```

The `-c env=dev` context value is used to name resources (e.g. `persons-dev`, `person-bus-dev`). Replace `dev` with `staging` or `prod` for other environments.

## CDK Commands

| Command | Description |
|---|---|
| `npm run build` | Compile TypeScript to JS |
| `npm run watch` | Watch and compile |
| `npx cdk deploy` | Deploy stack to AWS |
| `npx cdk diff` | Compare deployed stack with current state |
| `npx cdk synth` | Emit synthesized CloudFormation template |

## Environment Variables

| Variable | Used by | Description |
|---|---|---|
| `TABLE_NAME` | service | DynamoDB table name |
| `EVENT_BUS_NAME` | stream-service | EventBridge bus name |
