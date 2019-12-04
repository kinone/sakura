package console

type CommandInterface interface {
	Configure()
	Execute() error
	Name() string
	AddArgument(arg *Argument)
	Args() []*Argument
	SetApp(a *Application)
	App() *Application
}

type CommandTraits struct {
	app  *Application
	args []*Argument
}

func (c *CommandTraits) AddArgument(arg *Argument) {
	c.args = append(c.args, arg)
}

func (c *CommandTraits) Args() []*Argument {
	return c.args
}

func (c *CommandTraits) SetApp(a *Application) {
	c.app = a
}

func (c *CommandTraits) App() *Application {
	return c.app
}
