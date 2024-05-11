package constant

const APP_NAME = "headline"

type Source string

const (
	SourceThebell Source = "thebell"
)

var Sources = []struct {
	Value  Source
	TSName string
}{
	{SourceThebell, "thebell"},
}
