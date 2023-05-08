package validation

var (
	// PerCallLimit specify the actual cost limit per CEL validation call
	// current PerCallLimit gives roughly 0.1 second for each expression validation call
	// PerCallLimit enables cost tracking and sets configures program evaluation to exit early with a
	// "runtime cost limit exceeded" error if the runtime cost exceeds the costLimit.
	// The PerCallLimit is a metric that corresponds to the number and estimated expense of operations
	// performed while evaluating an expression. It is indicative of CPU usage, not memory usage.
	PerCallLimit uint64 = 1000000
)
