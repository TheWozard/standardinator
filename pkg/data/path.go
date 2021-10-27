package data

import (
	"fmt"
	"strings"
)

func CreateNewRootDataPath() *DataPath {
	return &DataPath{
		path: []interface{}{"$"},
	}
}

type DataPath struct {
	path []interface{}
}

func (p *DataPath) EnterNamespace(namespace string) {
	p.path = append(p.path, namespace)
}

func (p *DataPath) EnterIndex(index int) {
	p.path = append(p.path, index)
}

func (p *DataPath) ExitNamespace() {
	if p.HeadIsNamespace() {
		p.path = p.path[:len(p.path)-1]
	}
}

func (p *DataPath) ExitIndex() {
	if p.HeadIsIndex() {
		p.path = p.path[:len(p.path)-1]
	}
}

func (p *DataPath) IncrementIndex() {
	index, ok := p.path[len(p.path)-1].(int)
	if !ok {
		return
	}
	p.path[len(p.path)-1] = index + 1
}

func (p *DataPath) HeadIsIndex() bool {
	_, ok := p.path[len(p.path)-1].(int)
	return ok
}

func (p *DataPath) HeadIsNamespace() bool {
	_, ok := p.path[len(p.path)-1].(string)
	return ok
}

func (p *DataPath) ToString() string {
	builder := strings.Builder{}
	for i, elem := range p.path {
		switch cast := elem.(type) {
		case string:
			if i != 0 {
				builder.WriteRune('.')
			}
			builder.WriteString(cast)
		case int:
			builder.WriteString(fmt.Sprintf("[%d]", cast))
		}
	}
	return builder.String()
}
