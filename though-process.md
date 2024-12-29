## Verve Assignment Thought Process

### Problem Statement

Build a high-performance REST service containing a GET endpoint, `/api/verve/accept`, capable of processing 10K requests per second.

#### Functional Requirements
- **GET API**:
  - Accepts `id` (integer) as a mandatory query parameter and `url` (string) as an optional parameter.
  - Returns `ok` if the request is successfully processed and `failed` otherwise.

- **Logging**:
  - Every minute, log the count of unique `id`s received during that minute.

- **POST Request**:
  - If the `url` parameter is provided, send the count of unique `id`s for the current minute via a POST request to the given URL.

- **Extension 1**:
  - Ensure `id` deduplication works seamlessly in a distributed environment (e.g., behind a load balancer).

- **Extension 2**:
  - Instead of logging to a file, send the unique count to a distributed streaming service.

#### Non-Functional Requirements
- **High Throughput**:
  - The service must handle a large volume of requests (10K+ requests per second).

- **Scalability**:
  - The service should scale horizontally to handle increased loads and distributed setups.

- **Maintainability**:
  - The codebase should be modular and easy to extend, debug, or refactor.

### Proposed Solution

#### Technology Choice
- **Programming Language**:
  - The service is implemented in **Go** for its efficient concurrency model, minimal runtime overhead, and strong support for high-performance server-side applications.

- **Router**:
  - Utilized `chi-router`, a lightweight and fast HTTP router, to handle API requests efficiently.

- **Code Design**:
  - Adopted the **repository pattern** with dependency injection to achieve modularity and maintainability.

#### Deduplication

- **Local Deduplication**:
  - We can use `sync.Mutex` to ensure thread-safe deduplication:
    - `mutex.Lock()` and `mutex.Unlock()` provide safe access to shared data structures in concurrent scenarios.

- **Distributed Deduplication**:
  - Leveraged **Redis**, a highly available in-memory data store:
    - **Set Data Structure**:
      - Redis' `SADD` command ensures O(1) time complexity for adding unique `id`s.
      - Redis' `SCARD` command provides O(1) time complexity for retrieving the count of unique `id`s.
    - **TTL Management**:
      - Since Redis sets do not support native TTL for individual keys, we use `DEL` to expire keys efficiently (O(N) for multiple keys, ~O(1) for single keys).

#### Logging
- Integrated **slog**, Go's standard structured logging library, for:
  - Timestamps and log levels.
  - Enhanced debugging and observability.

#### Streaming Service
- **Kafka**:
  - Chosen for its ability to handle high-throughput event streaming and scalability.

#### Background Tasks
- **Periodic Unique Count Logging**:
  - Scheduled using Go's built-in `time.Ticker` to log or send unique counts every minute.
  - Runs as a background task during application initialization.

- **Future Enhancements**:
  - For complex background workflows, cron-like schedulers or specialized task libraries can be explored.

### Architecture Diagram

Below is a conceptual representation of the system architecture:

```
 +-------------+           +--------------+          +-----------------+
 |   Client    |           |    Redis     |          |   Kafka Stream  |
 | (API Call)  + --------->+ (Dedup Cache)+--------->+  (Data Pipeline)|
 +-------------+           +--------------+          +-----------------+                              
        |                       ^   
        v                       |   
 +-------------+           +------------+  
 |  REST API   |           | Ticker     |
 | (Go Service)|           |(Go Service)|
 +-------------+           +------------+ 
        |
        v                   
 +-----------------+
 | Structured Logs |    
 +-----------------+
```

---

### Conclusion
This design utilizes Go's concurrency model, Redis' efficient data structures for deduplication, and Kafka's scalability for streaming. By adhering to best practices like modular architecture and structured logging, the service is prepared to meet the outlined functional and non-functional requirements effectively while remaining flexible for future enhancements.
