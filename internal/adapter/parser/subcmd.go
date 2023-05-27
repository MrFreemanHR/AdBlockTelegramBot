package parser

type Subcmd struct {
	Name         string
	Child        *Subcmd
	RequiredArgs map[byte]RequiredArg
	OptionalArgs map[byte]OptionalArg
}

func NewSubCmd(name string) Subcmd {
	c := Subcmd{
		Name: name,
	}
	c.RequiredArgs = make(map[byte]RequiredArg)
	c.OptionalArgs = make(map[byte]OptionalArg)
	return c
}
