package parser

type Cmd struct {
	Name         string
	ChildSubcmd  *Subcmd
	RequiredArgs map[byte]RequiredArg
	OptionalArgs map[byte]OptionalArg
}

func NewCmd(name string) Cmd {
	c := Cmd{
		Name: name,
	}
	c.RequiredArgs = make(map[byte]RequiredArg)
	c.OptionalArgs = make(map[byte]OptionalArg)
	return c
}

func GetCommandNames() []string {
	return []string{
		"locales",
	}
}
