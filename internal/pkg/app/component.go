package app

import "fyne.io/fyne/v2"

type Component struct {
	Name           string
	Object         *fyne.Container
	Enabled        bool
	Props          map[string]interface{}
	ChildComponent *Component
}

func (c *Component) Enable() {
	c.Enabled = true
}

func (c *Component) Disable() {
	c.Enabled = false
}

func (c *Component) Child() *Component {
	return c.ChildComponent
}

func (c *Component) AddChild(child *Component) {
	c.Object.Add(child.Object)
	c.ChildComponent = child
}

func (c *Component) AddChildDeep(child *Component) {
	marker := c
	for marker.ChildComponent != nil {
		marker = marker.ChildComponent
	}
	marker.AddChild(child)
}

func (c *Component) RemoveChild(child *Component) {
	c.Object.Remove(child.Object)
	c.ChildComponent = nil
}

func NewComponent(obj *fyne.Container) *Component {
	return &Component{
		Object:  obj,
		Props:   make(map[string]interface{}),
		Enabled: true,
	}
}
