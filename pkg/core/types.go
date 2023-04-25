package core

// Command represents a command or a sub-command of the `grofer` CLI.
type Command int

const (
	// RootCommand is the root command of grofer, i.e.
	// `grofer`.
	RootCommand Command = iota
	// ProcCommand is `grofer proc` and its variants.
	ProcCommand
	// ContainerCommand is `grofer container` and its
	// variants.
	ContainerCommand
	// ExportCommand is `grofer export` and its variants.
	ExportCommand
)