# State Machine Library

The State Machine Library provides a structured way to manage state transitions using a chain of responsibility pattern. It ensures that state transitions act like transactions, meaning that if any handler in the chain fails, the system will attempt to roll back to the original state.

Certainly! The State Machine Library is designed with specific objectives in mind and addresses several challenges commonly encountered in system design and event-driven architectures. Here's a summary:

## Goals:

1. **Structured State Management**: Provide a structured approach to manage state transitions in systems, ensuring that transitions occur in a controlled and predictable manner.

2. **Idempotency**: Ensure that state transitions are idempotent, meaning they can be retried without side effects, thus enhancing system reliability.

3. **Flexibility**: Offer a framework that is flexible enough to handle various types of state transitions, accommodating both default and custom transition logic.

4. **Transaction-like Behavior**: Make sure that state transitions are atomic, i.e., either they complete fully or get rolled back, ensuring data consistency.

5. **Extensibility**: Allow users to extend the default behavior with custom handlers and logic, catering to specific application needs.

6. **Logging & Monitoring**: Provide built-in logging and monitoring capabilities to track state transitions, making it easier to debug and monitor system behavior.

## Problems It Solves:

1. **Complex State Management**: Managing state can become complex, especially in large systems with numerous possible states. The library offers a structured way to manage these transitions, ensuring that only valid transitions occur.

2. **Handling Failures**: In distributed systems, failures are a norm rather than an exception. The library ensures that if a transition fails, the system can either retry or roll back to a consistent state.

3. **Ensuring Idempotency**: In event-driven architectures, events might get delivered more than once. By using unique event IDs and checking for processed events, the library ensures that state transitions are idempotent.

4. **Custom Logic Integration**: Every application has unique requirements. The library allows users to integrate custom logic seamlessly, ensuring that the state machine can cater to diverse needs.

5. **Manual Interventions**: Some failures require manual intervention. The library includes a `ManualReview` state to flag such scenarios, allowing operations teams to step in.

6. **Consistent Logging**: Keeping track of state transitions can be challenging. The structured logging ensures that transitions, failures, and successes are consistently logged, aiding in debugging and monitoring.

## Key Concepts

- **StateObject**: A generic object that holds data and its current state. It is designed to be flexible so users can attach any data they need for state transitions.
- **Handlers**: Handlers are functions that perform specific tasks in the state transition process. Each handler must also have a rollback mechanism defined in case of failure.
- **Chain of Responsibility Pattern**: Handlers are designed in a chain where each handler passes the request to the next handler in the chain. If any handler fails, the rollback for that handler is executed, ensuring a consistent state.
- **StateTransition**: Represents a transition from one state to another.
- **EventID**: A unique identifier for a state transition to ensure idempotency.
- **Logging**: Structured logging to stdout, capturing the start, any failures, and the conclusion of transitions.

### Hashing the Event ID
To generate a unique Event ID, the library hashes the event content combined with the current date. This ensures that the same event content on different days will have distinct Event IDs. The SHA-256 cryptographic hash function is used, providing a good balance of speed and security.

### Rollbacks

Each handler in the chain has a rollback function defined. If a handler fails, the system will attempt to roll back to the original state using the rollback functions of the handlers executed before the failure.

In certain scenarios, if the rollback function of a handler determines that an action can't be reverted, the state is moved to `ManualReview` for manual intervention.

## ManualReview State

The `ManualReview` state is a special state that indicates a need for manual intervention. This state is entered when there's a failure in the handler chain that can't be automatically reverted. It allows operations teams to manually investigate and correct any issues.

## Using the State Machine

### Initialization

```go
stateMachine := statemachine.NewStateMachine()
```

### Configuration

```go
stateMachine.LogTransitions = true
stateMachine.DebugLogging = true
```

### Registering Transitions

```go
stateMachine.RegisterTransition("FromState", "ToState")
```

#### Creating a Chain of Handlers

```go
handlers := statemachine.CreateHandlerChain(handler1, handler2, ...)
```

### Handling State Transitions

To handle a state transition:

```go
stateObject := statemachine.NewStateObject(map[string]interface{}{
    "PhoneNumber": "1234567890",
    "Carrier": "TelecomProvider",
})
err := stateObject.TransitionTo("NewState", "UniqueEventID", stateMachine)
```

### Using the Helper Function to Create Handler Chain

The `CreateHandlerChain` function helps in creating a chain of handlers:

```go
handlerChainStart := statemachine.CreateHandlerChain(handler1, handler2, handler3)
```

Absolutely, here's a dedicated section for configuring the `StateMachine` object:

---

### Configuring the StateMachine

The `StateMachine` object is central to managing state transitions. It provides various configurations to tailor its behavior according to the application's needs. Here's how you can customize and configure your `StateMachine`:

1. **Initialization**:

   First, initialize a new `StateMachine` instance:

   ```go
   stateMachine := statemachine.NewStateMachine()
   ```

2. **Logging Configuration**:

   The library offers structured logging to stdout. You can control the logging of transitions and decide whether to include debug information:

   ```go
   // Enable logging of transitions
   stateMachine.LogTransitions = true
   
   // Enable debug logging (includes caller info)
   stateMachine.DebugLogging = true
   ```

3. **Handler Configuration**:

   The `StateMachine` allows you to set configurations for the default handlers:

   ```go
   config := statemachine.HandlerConfig{
       CheckEventID:   true,
       CheckProcessed: true,
       Telemetry:      true,
       Alerting:       true,
       MarkProcessed:  true,
   }
   stateMachine.SetHandlerConfig(config)
   ```

   You can turn off any default handler by setting its value to `false` in the `HandlerConfig`.

4. **Database and Event Emitters**:

   If your application integrates with databases or event emitters, you can configure the `StateMachine` to use these connections:

   ```go
   // Set a database connection
   stateMachine.SetDB(myDBConnection)
   
   // Set event emitters (Kafka, NATS, RabbitMQ, etc.)
   stateMachine.SetKafkaConn(myKafkaConnection)
   stateMachine.SetNatsConn(myNatsConnection)
   stateMachine.SetRabbitMQConn(myRabbitMQConnection)
   ```

5. **Registering Transitions**:

   After configuring, you can register state transitions and associate them with handler chains:

   ```go
   stateMachine.RegisterTransition("FromState", "ToState", customHandler1, customHandler2, ...)
   ```


### Default Handlers

The State Machine Library provides a series of default handlers designed to manage common aspects of state transitions. These handlers are executed in a specific sequence to ensure structured transitions.

**Order of Execution:**
1. **CheckEventIDHandler**: Checks and creates the eventID if missing.
2. **CheckProcessedHandler**: Verifies if the event was already processed.
3. **TelemetryHandler**: Handles telemetry logic.
4. **AlertingHandler**: Manages alerting logic.
5. **Custom Handlers (if any, but not overriding default ones)**
6. **MarkProcessedHandler**: Marks the event as processed.

#### Overriding Default Handlers

You can override any of the default handlers by providing a custom handler when registering a transition. When overridden, the custom handler takes the place of the default handler in the chain.

**Type Definitions for Overriding:**
- For `CheckEventIDHandler`:

```go
type CustomCheckEventIDHandler struct {
    next statemachine.Handler
}
```

- For `CheckProcessedHandler`:

```go
type CustomCheckProcessedHandler struct {
    next statemachine.Handler
}
```

- And similarly for other handlers.

#### How to Override

1. Create a custom handler following the respective type definition:

```go
type CustomCheckEventIDHandler struct {
    next statemachine.Handler
}

func (h *CustomCheckEventIDHandler) Handle(so *statemachine.StateObject, state string) bool {
    // Your custom logic here
    return true // or false based on your logic
}

func (h *CustomCheckEventIDHandler) SetNext(handler statemachine.Handler) {
    h.next = handler
}

func (h *CustomCheckEventIDHandler) Rollback(so *statemachine.StateObject, state string) bool {
    // Your rollback logic here
    return true
}
```

2. Register the transition with your custom handler:

```go
customHandler := &CustomCheckEventIDHandler{}
stateMachine.RegisterTransition("FromState", "ToState", customHandler)
```

By overriding a default handler, you can customize its behavior while still maintaining the structure and order of the handler chain. Non-default custom handlers will be placed in the chain right before the `MarkProcessedHandler`, ensuring that your custom logic is executed before finalizing the state transition.


## Examples

### Basic State Transition

```go
stateMachine := statemachine.NewStateMachine()
stateMachine.RegisterTransition(statemachine.SIMNotActivated, statemachine.SIMActivated)
user := statemachine.NewUser("JohnDoe")
err := user.TransitionTo(statemachine.SIMActivated, "")
if err != nil {
    log.Println("Error:", err)
}
```

### Using Custom Handlers

Assuming you have a custom handler for telemetry:

```go
customTelemetryHandler := &CustomTelemetryHandler{}
stateMachine.RegisterTransition(statemachine.SIMNotActivated, statemachine.SIMActivated, customTelemetryHandler)
```

This will replace the default telemetry handler with your custom implementation for the specified transition.

## Telecom SIM Activation Example

Consider a telecom system with states for SIM cards. Use the state machine library to handle transitions:

1. **Initialization**:

   ```go
   stateMachine := statemachine.NewStateMachine()
   ```

2. **Handling a SIM Activation**:

   ```go
   simData := map[string]interface{}{
       "PhoneNumber": "1234567890",
       "Carrier": "TelecomProvider",
   }
   sim := statemachine.NewStateObject(simData)
   err := sim.TransitionTo(statemachine.SIMActivated, "UniqueEventID", stateMachine)
   ```

3. **Using simData in Custom Handlers**:

   ```go
   func (h *CustomCarrierCheckHandler) Handle(so *statemachine.StateObject, state string) bool {
       carrier, exists := so.Data["Carrier"]
       if !exists || carrier != "ValidCarrier" {
           return false
       }
       return true
   }
   ```

The State Machine Library ensures idempotent state transitions with built-in rollback mechanisms. This provides a safe and structured way to manage complex workflows.
