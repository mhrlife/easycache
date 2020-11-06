package easycache

import "fmt"

type ResourceDidntFindInCache struct {
	slug string
}

func (p ResourceDidntFindInCache) Error() string {
	return fmt.Sprintf("%s didn't found in cache", p.slug)
}

type ResourceLayerUndefined struct {
	slug  string
	layer int
}

func (p ResourceLayerUndefined) Error() string {
	return fmt.Sprintf("%s dont have %d layer", p.slug, p.layer)
}

type ResourceProvideError struct {
	slug string
}

func (r ResourceProvideError) Error() string {
	return fmt.Sprintf("%s coudnl't provide data", r.slug)
}

type CannotProvide struct {
}

func (r CannotProvide) Error() string {
	return fmt.Sprintf("cannot provide")
}

type ResourceNotFound struct {
	slug string
}

func (r ResourceNotFound) Error() string {
	return fmt.Sprintf("%s resource not found", r.slug)
}
