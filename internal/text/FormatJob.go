package text

type FormatJob interface {
	Run() (string, error)
}
