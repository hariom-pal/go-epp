//go:build live

package live

import (
	"context"
	"errors"
	"testing"

	"github.com/hariom-pal/go-epp/epp"
	"github.com/hariom-pal/go-epp/types"
)

// TestCancelDoesNotAffectNextCommand covers a cross-command hazard in the
// per-command cancellation watcher. The watcher goroutine selects on both
// ctx.Done() and a completion channel; when a caller cancels with the
// idiomatic `defer cancel()`, both become ready at once and Go picks a case at
// random. If it picks ctx.Done() after the command already returned, it sets
// the shared connection's read deadline into the past, aborting whichever
// command runs next on that session.
func TestCancelDoesNotAffectNextCommand(t *testing.T) {
	client := session(t)

	for i := 0; i < 60; i++ {
		func() {
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			if _, err := client.DomainCheckContext(ctx, types.DomainCheckRequest{
				Domains: []string{domainName("cr")},
			}); err != nil {
				t.Fatalf("iteration %d: command failed: %v", i, err)
			}
		}()

		// The very next command must not inherit the previous context's
		// cancellation through the shared read deadline.
		if _, err := client.DomainCheck(types.DomainCheckRequest{
			Domains: []string{domainName("cr")},
		}); err != nil {
			var sdkErr *epp.SDKError
			if errors.As(err, &sdkErr) && sdkErr.Kind == epp.ErrorKindTimeout {
				t.Fatalf("iteration %d: a previous command's cancellation aborted this one: %v", i, err)
			}
			t.Fatalf("iteration %d: follow-up command failed: %v", i, err)
		}
	}
}
