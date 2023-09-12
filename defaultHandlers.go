package statemachine

type CheckEventIDHandler struct {
	next Handler
}

func (h *CheckEventIDHandler) Handle(u *StateObject, state string) bool {
	if u.EventID == "" {
		u.EventID = generateSignature(state)
	}
	return h.next.Handle(u, state)
}

func (h *CheckEventIDHandler) Rollback(u *StateObject, state string) bool {
	// Define rollback logic, return false if cannot rollback
	return true
}

func (h *CheckEventIDHandler) SetNext(handler Handler) {
	h.next = handler
}

func (h *CheckEventIDHandler) Next() Handler {
	return h.next
}

type CheckProcessedHandler struct {
	next         Handler
	stateMachine *StateMachine
}

func (h *CheckProcessedHandler) Handle(u *StateObject, state string) bool {
	if h.stateMachine.isProcessed(u.EventID) {
		return true
	}
	return h.next.Handle(u, state)
}

func (h *CheckProcessedHandler) Rollback(u *StateObject, state string) bool {
	// Define rollback logic, return false if cannot rollback
	return true
}

func (h *CheckProcessedHandler) SetNext(handler Handler) {
	h.next = handler
}

func (h *CheckProcessedHandler) Next() Handler {
	return h.next
}

type MarkProcessedHandler struct {
	stateMachine *StateMachine
}

func (h *MarkProcessedHandler) Handle(u *StateObject, state string) bool {
	h.stateMachine.markProcessed(u.EventID)
	return true
}

func (h *MarkProcessedHandler) Rollback(u *StateObject, state string) bool {
	// Define rollback logic, return false if cannot rollback
	return true
}

func (h *MarkProcessedHandler) SetNext(handler Handler) {}

func (h *MarkProcessedHandler) Next() Handler {
	return nil
}

type TelemetryHandler struct {
	next      Handler
	telemetry bool
}

func (h *TelemetryHandler) Handle(u *StateObject, state string) bool {
	if h.telemetry {
		// Fire telemetry logic here
	}
	if h.next != nil {
		return h.next.Handle(u, state)
	}
	return false
}

func (h *TelemetryHandler) Rollback(u *StateObject, state string) bool {
	// Define rollback logic, return false if cannot rollback
	return true
}

func (h *TelemetryHandler) SetNext(handler Handler) {
	h.next = handler
}

func (h *TelemetryHandler) Next() Handler {
	return h.next
}

type AlertingHandler struct {
	next     Handler
	alerting bool
}

func (h *AlertingHandler) Handle(u *StateObject, state string) bool {
	if h.alerting {
		// Fire alerting logic here
	}
	if h.next != nil {
		return h.next.Handle(u, state)
	}
	return false
}

func (h *AlertingHandler) Rollback(u *StateObject, state string) bool {
	// Define rollback logic, return false if cannot rollback
	return true
}

func (h *AlertingHandler) SetNext(handler Handler) {
	h.next = handler
}

func (h *AlertingHandler) Next() Handler {
	return h.next
}

// CustomCheckEventIDHandler
type CustomCheckEventIDHandler struct {
	next Handler
}

func (h *CustomCheckEventIDHandler) Handle(u *StateObject, state string) bool {
	// Your custom logic here
	return h.next.Handle(u, state)
}

func (h *CustomCheckEventIDHandler) Rollback(u *StateObject, state string) bool {
	// Your custom rollback logic
	return true
}

func (h *CustomCheckEventIDHandler) SetNext(handler Handler) {
	h.next = handler
}

func (h *CustomCheckEventIDHandler) Next() Handler {
	return h.next
}

// CustomCheckProcessedHandler
type CustomCheckProcessedHandler struct {
	next Handler
}

func (h *CustomCheckProcessedHandler) Handle(u *StateObject, state string) bool {
	// Your custom logic here
	return h.next.Handle(u, state)
}

func (h *CustomCheckProcessedHandler) Rollback(u *StateObject, state string) bool {
	// Your custom rollback logic
	return true
}

func (h *CustomCheckProcessedHandler) SetNext(handler Handler) {
	h.next = handler
}

func (h *CustomCheckProcessedHandler) Next() Handler {
	return h.next
}

// CustomTelemetryHandler
type CustomTelemetryHandler struct {
	next Handler
}

func (h *CustomTelemetryHandler) Handle(u *StateObject, state string) bool {
	// Your custom logic here
	return h.next.Handle(u, state)
}

func (h *CustomTelemetryHandler) Rollback(u *StateObject, state string) bool {
	// Your custom rollback logic
	return true
}

func (h *CustomTelemetryHandler) SetNext(handler Handler) {
	h.next = handler
}

func (h *CustomTelemetryHandler) Next() Handler {
	return h.next
}

// CustomAlertingHandler
type CustomAlertingHandler struct {
	next Handler
}

func (h *CustomAlertingHandler) Handle(u *StateObject, state string) bool {
	// Your custom logic here
	return h.next.Handle(u, state)
}

func (h *CustomAlertingHandler) Rollback(u *StateObject, state string) bool {
	// Your custom rollback logic
	return true
}

func (h *CustomAlertingHandler) SetNext(handler Handler) {
	h.next = handler
}

func (h *CustomAlertingHandler) Next() Handler {
	return h.next
}

// CustomMarkProcessedHandler
type CustomMarkProcessedHandler struct {
	next Handler
}

func (h *CustomMarkProcessedHandler) Handle(u *StateObject, state string) bool {
	// Your custom logic here
	return h.next.Handle(u, state)
}

func (h *CustomMarkProcessedHandler) Rollback(u *StateObject, state string) bool {
	// Your custom rollback logic
	return true
}

func (h *CustomMarkProcessedHandler) SetNext(handler Handler) {
	h.next = handler
}

func (h *CustomMarkProcessedHandler) Next() Handler {
	return h.next
}
