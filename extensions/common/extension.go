package common

// Namespace identifies an EPP extension XML namespace.
type Namespace string

// Builder is implemented by extension request builders that can report whether
// they have XML content to attach to an EPP command.
type Builder interface {
	Empty() bool
}
