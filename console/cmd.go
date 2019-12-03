package console

type CommandInterface interface {
	Configure()
	Execute() error
	Name() string
	AddArgument(arg *Argument)
	Args() []*Argument
}

type Command struct {
	args []*Argument
}

func (c *Command) AddArgument(arg *Argument) {
	c.args = append(c.args, arg)
}

func (c *Command) Args() []*Argument {
	return c.args
}

func (c *Command) Configure()           {}
func (c *Command) Execute() (err error) { return }
func (c *Command) Name() (n string)     { return }
