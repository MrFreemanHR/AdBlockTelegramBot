package parser

type OptionalArg struct {
	value string
}

func (a *OptionalArg) SetValue(value string) {
	a.value = value
}

func (a *OptionalArg) GetValue() string {
	return a.value
}
