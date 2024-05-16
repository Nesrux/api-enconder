package repositories

import "github.com/Nesrux/api-enconder/domain"

type VideoRepository interface {
	Insert(video *domain.Video) (*domain.Video, error)
	Find(id string) (*domain.Video, error)
}

type VideoRepositoryDb struct {
}

func (vr VideoRepositoryDb) Insert(video *domain.Video) (*domain.Video, error) {
	return nil, nil
}

func (vr VideoRepositoryDb) Find(id string) (*domain.Video, error) {
	return nil, nil
}
