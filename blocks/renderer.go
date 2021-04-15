package blocks

type Renderer interface {
	Activate(string)
	Disable(string)
}
