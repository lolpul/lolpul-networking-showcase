# Discussing the examples

Use the files and tests to reason about the decisions, rather than memorize a script.

## Writer session

Problem: a writer may be unsafe for concurrent access, and closing it while a write is active violates ownership. Constraint: `io.WriteCloser` has no cancellation contract. Decision: serialize writes and close under one mutex, transfer ownership at construction, and retain the first close result. Alternative: track active operations separately, allowing simultaneous writes only when the backend supports them. Trade-off: the small contract is easy to review, but a blocked write also blocks shutdown and state inspection.

Questions: Where does ownership transfer? What happens on a short write? Why does close not retry? Could calling back into the session deadlock? How would cancellable I/O change this design?

## Cancellation task

Problem: cancelling a context is a request, not evidence that a goroutine and its resources are gone. Constraint: work must cooperate and clean up before returning. Decision: derive the context, own one goroutine, and join using a closed channel that also publishes the result. Alternative: a synchronous operation when background execution adds no value; an error group when coordinating multiple workers. Trade-off: no forced termination or panic recovery is provided, and successful work remains successful even if cancellation arrives after its result.

Questions: Why is result access race-free? What if cancellation arrives before the worker starts? What if work ignores context? Who must call Wait or Stop? How would you handle a worker panic or multiple tasks?

## Scoped resource cleanup

Problem: the second acquisition can fail after the first succeeds, and an acquisition can return a resource with an error. Constraint: each non-nil return transfers ownership; use only borrows it. Decision: register each cleanup immediately and unwind in reverse order, preserving errors with `errors.Join`. Alternative: a reusable cleanup stack for a genuinely variable resource count. Trade-off: two explicit resources keep the example readable; this is not a general transaction manager, and Close has no deadline.

Questions: Which failures still require cleanup? Why does a typed nil break the contract? Can cleanup replace the operation error? Would a context-aware cleanup inherit a cancelled request or receive a separate budget? When would a cleanup stack justify its added abstraction?
