# Food Delivery Microservices Architecture

## Microservices Architecture Overview

We've created a food delivery application (similar to Bolt Food or Glovo) using a microservices architecture. Here are the main components:

### Microservices

1. **User Service** (Spring Boot, port 8080)
   - Manages user information
   - Provides API for creating, reading, updating, and deleting users
   - Implemented in Java with Spring Boot

2. **Restaurant Service** (Go, port 8081)
   - Manages restaurants and their menus
   - Provides API for managing restaurants and menu items
   - Implemented in Go with Gorilla Mux

3. **Order Service** (Go, port 8082)
   - Manages user orders
   - Processes order creation and status updates
   - Communicates with other services to notify about changes
   - Implemented in Go with Gorilla Mux

4. **Payment Service** (Go, port 8083)
   - Manages payments for orders
   - Processes payments (simulation)
   - Communicates with Order Service to update order status
   - Implemented in Go with Gorilla Mux

5. **Delivery Service** (Go, port 8084)
   - Manages deliveries and couriers
   - Assigns couriers to orders
   - Updates delivery status and notifies Order Service
   - Implemented in Go with Gorilla Mux

6. **Notification Service** (Go, port 8085)
   - Manages user notifications
   - Receives events from other services and creates notifications
   - Provides API for marking notifications as read
   - Implemented in Go with Gorilla Mux

7. **Frontend** (HTML, CSS, JavaScript)
   - Simple user interface
   - Allows searching for restaurants, placing orders, and viewing status
   - Connects to all microservices to get and send data

### Communication Between Microservices

The microservices communicate with each other through a synchronous communication model using HTTP REST. Each service exposes a REST API that other services can consume. For example:

- When a user places an order, the **Order Service** notifies the **Payment Service**
- When a payment is processed, the **Payment Service** notifies the **Order Service** to update the order status
- When an order is paid, the **Order Service** notifies the **Delivery Service** to assign a courier
- The **Notification Service** is notified by other services to send notifications to users

### Containerization

Each microservice has its own Dockerfile to facilitate building and running in containers. This allows deployment in Kubernetes and integration with Istio Service Mesh later.

### Integration with Kubernetes and Istio

To integrate this architecture with Kubernetes and Istio, you should follow these steps:

1. **Build Docker images for each microservice**
   ```bash
   # For each service
   cd user-service/
   docker build -t user-service:latest .
   
   cd ../restaurant-service/
   docker build -t restaurant-service:latest .
   
   # And so on for each service
   ```

2. **Create a Kubernetes cluster and install Istio**
   ```bash
   # Install Istio
   istioctl install --set profile=demo -y
   kubectl label namespace default istio-injection=enabled
   ```

3. **Create Kubernetes manifests for each service**

   You need to create the following resources for each microservice:
   - Deployment
   - Service
   - ConfigMap/Secret (if necessary)
   - Possibly a ServiceEntry for external access

4. **Use Istio for traffic management**
   
   You can use the following Istio resources:
   - VirtualService - for traffic routing
   - DestinationRule - for connection policies
   - Gateway - for exposing services externally
   - ServiceEntry - for access to external services
   - Sidecar - for configuring the Envoy proxy

5. **Implement Istio features**
   - Distributed tracing (with Jaeger)
   - Monitoring (with Prometheus and Grafana)
   - Traffic visualization (with Kiali)
   - Canary deployments and A/B testing
   - Circuit breaking

### Examples of Istio Features for Testing

1. **Implementing a canary deployment strategy**:
   - Release a new version of a service to a subset of users
   - Monitor performance and errors
   - Gradually increase traffic to the new version

2. **Configuring circuit breakers**:
   - Protect services against cascading failures
   - Limit the number of connections and concurrent requests
   - Implement timeout policies

3. **Securing communication between services**:
   - Implement mutual TLS authentication (mTLS)
   - Define authorization policies for access control

4. **Implementing resiliency patterns**:
   - Retry logic for failed requests
   - Timeouts to prevent blocking
   - Fault injection for testing resilience

## Conclusions

The implemented architecture follows the principles of microservices and is ready to be deployed in Kubernetes and integrated with Istio Service Mesh. The code is minimalistic but functional and can be extended according to your requirements.

To test and integrate this system with Kubernetes and Istio, follow the steps mentioned above and make sure you have correctly installed all the necessary components. Istio offers numerous advanced features for traffic management, security, and observability that will allow you to build a robust and scalable application.
