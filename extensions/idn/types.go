// Package idn implements the idn-1.0 EPP extension, which carries the
// language table and Unicode form of an internationalised domain name.
//
// Registries that accept IDNs generally require the table tag on a create
// command so they can apply the right script and variant rules; the A-label
// in <domain:name> alone does not say which table the name was validated
// against.
package idn

import "github.com/hariom-pal/go-epp/constants"

// Namespace is the idn-1.0 XML namespace.
const Namespace = constants.IDNNamespace

// Data contains the idn-1.0 extension fields.
type Data struct {
	// Table is the registry's language/script table tag, such as "hi".
	Table string

	// UName is the Unicode (U-label) form of the domain name. It is
	// optional; registries derive it from the A-label when it is omitted.
	UName string
}

// CreateRequest contains optional IDN data for a domain create command.
type CreateRequest struct {
	Data
}
