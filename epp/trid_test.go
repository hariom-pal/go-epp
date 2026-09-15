package epp

import "testing"

// TestClientTRIDIsUniqueAcrossClients covers the NIXI OT&E finding that two
// sessions opened in the same second emitted identical client transaction
// identifiers, because the identifier was built only from a command prefix, a
// one-second timestamp and a per-client counter that always starts at zero.
// RFC 5730 section 2.5 requires uniqueness: the identifier is how a registrar
// correlates a response and reconciles a transform whose reply was lost, so a
// collision between pooled connections makes an ambiguous transform
// unreconcilable.
func TestClientTRIDIsUniqueAcrossClients(t *testing.T) {
	const clients = 8
	const perClient = 25

	seen := map[string]bool{}
	for i := 0; i < clients; i++ {
		client := &Client{}
		for j := 0; j < perClient; j++ {
			trid := client.nextTRID("TRANSFER")
			if len(trid) > 64 {
				t.Fatalf("clTRID %q exceeds the RFC 5730 clTRIDType limit of 64 characters", trid)
			}
			if len(trid) < 3 {
				t.Fatalf("clTRID %q is below the RFC 5730 clTRIDType minimum of 3 characters", trid)
			}
			if seen[trid] {
				t.Fatalf("clTRID %q was emitted by two different clients", trid)
			}
			seen[trid] = true
		}
	}
	if len(seen) != clients*perClient {
		t.Fatalf("expected %d unique identifiers, got %d", clients*perClient, len(seen))
	}
}

// TestClientTRIDIsStablePerClient confirms one client keeps a single nonce, so
// its identifiers stay correlatable within a session.
func TestClientTRIDIsStablePerClient(t *testing.T) {
	client := &Client{}

	first := client.nextTRID("CHECK")
	second := client.nextTRID("CHECK")

	if first == second {
		t.Fatal("a client must not repeat a transaction identifier")
	}
	if client.tridNonce() == "" {
		t.Fatal("client nonce was not generated")
	}
	if nonce := client.tridNonce(); client.tridNonce() != nonce {
		t.Fatal("client nonce must be stable for the life of the client")
	}
}
