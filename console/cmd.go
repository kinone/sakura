package console

type CommandInterface interface {
	Configure()
	Execute() error
	Name() string
	AddArgument(arg *Argument)
	Args() []*Argument
}

type CommandTraints struct {
	args []*Argument
}

func (c *CommandTraints) AddArgument(arg *Argument) {
	c.args = append(c.args, arg)
}

func (c *CommandTraints) Args() []*Argument {
	return c.args
}