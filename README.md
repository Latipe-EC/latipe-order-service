# Order Service
> - Programming Language: Go 1.20
> - Web Framework: Fiber v2
> - ORM: Gorm
> - Dependency Injection: Wire
> - Config: Viper
> - Environment: Docker (Linux)
> - Tools: RabbitMQ, Prometheus, Grafana

Handling the order process in e-commerce using microservices architecture.
The Order service includes two components: order-rest-api and order-worker.
- **Order RestAPI:** Handles HTTP requests, authentication, authorization, and data validation.
- **Order Worker:** Listens for order events, commits transactions to other services, and interacts with databases for tasks such as creating and canceling orders.

The order creation process involves two phases:
- **Phase 1:** Processes HTTP POST requests, retrieves data by making HTTP requests to other services, and sends messages into the event queue.
- **Phase 2:** Receives messages, commits transactions, and saves data into the database.

![order-processing](/doc/order_proceess.png)

<hr></hr>

### Setup:

You can build the Docker image to run the project or use the provided Makefile.
- **Docker:** Navigate to the sub-folder and run `docker build -t <service name> .`
- **Makefile:** Use `make run-api` to run the REST API and `make run-worker` to run the WORKER
<hr></hr>

### Application Developed by Tran Tien Dat
