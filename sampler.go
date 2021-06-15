package thoth

type Sampler interface {
	Samples(name string) []Model
}
