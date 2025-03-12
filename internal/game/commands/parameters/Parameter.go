package parameters

type Parameter interface {
	Name() string
	IsRequired() bool
	IsMatch(text string) bool
	Consume(text string) (string, string)
}
