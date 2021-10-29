package flattener

import "strings"

const serperator = "."

// Generic way of identifying a location in a provided document and matching locations
type Path struct {
	stack []string
}

func NewPath(path string) *Path {
	return &Path{
		stack: strings.Split(path, serperator),
	}
}

func (p *Path) Enter(namespace string) {
	p.stack = append(p.stack, namespace)
}

func (p *Path) Exit() {
	p.stack = p.stack[:len(p.stack)-1]
}

func (p *Path) Includes(target *Path) bool {
	for offset := 0; offset < len(target.stack); offset++ {
		if target.stack[len(target.stack)-offset-1] != p.stack[len(p.stack)-offset-1] {
			return false
		}
	}
	return true
}

func (p *Path) toString() string {
	return strings.Join(p.stack, serperator)
}
