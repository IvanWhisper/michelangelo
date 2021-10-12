package dependencies

type IComponent interface {
	GetName() string
	GetInstance() (interface{}, error)
	GetLifeTimeScope() string
}

type Component struct {
	Name          string
	LifeTimeScope string
	NewFunc       func() (interface{}, error)
}

func (c *Component) GetName() string {
	return c.Name
}
func (c *Component) GetInstance() (interface{}, error) {
	return c.NewFunc()
}
func (c *Component) GetLifeTimeScope() string {
	return c.LifeTimeScope
}
