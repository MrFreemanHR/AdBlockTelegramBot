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

func (s *Subcmd) GetRequiredArg(position byte) *RequiredArg {
	if val, ok := s.RequiredArgs[position]; ok {
		return &val
	}
	return nil
}

func (s *Subcmd) GetOptionalArg(position byte) *OptionalArg {
	if val, ok := s.OptionalArgs[position]; ok {
		return &val
	}
	return nil
}
