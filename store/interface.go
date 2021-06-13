package store

type Store interface {
	Set(assets []Asset) error
	Get() ([]Asset, error)
}
