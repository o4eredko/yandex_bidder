package repository

import (
	dbx "github.com/go-ozzo/ozzo-dbx"

	"gitlab.jooble.com/marketing_tech/yandex_bidder/domain"
	"gitlab.jooble.com/marketing_tech/yandex_bidder/usecase"
)

type strategyRepo struct {
	db dbx.Builder
}

func NewStrategyRepo(db dbx.Builder) usecase.StrategyRepo {
	return &strategyRepo{
		db: db,
	}
}

func (r *strategyRepo) GetAll() ([]*domain.Strategy, error) {
	rows, err := r.db.
		Select("id", "name").
		From("strategies").
		Rows()
	if err != nil {
		return nil, err
	}

	strategies := make([]*domain.Strategy, 0)
	for rows.Next() {
		strategy := new(domain.Strategy)
		if err := rows.ScanStruct(strategy); err != nil {
			return nil, err
		}
		strategies = append(strategies, strategy)
	}

	return strategies, nil
}
