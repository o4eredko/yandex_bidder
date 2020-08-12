package usecase

import "gitlab.jooble.com/marketing_tech/yandex_bidder/domain"

type (
	GroupRepo interface {
		GetAll() ([]*domain.Group, error)
	}
)
