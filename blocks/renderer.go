package blocks

type Renderer interface {
	Activate(string)
	Disable(string)
}
type Notifier interface {
	Push(string, string, string, string) error
}
