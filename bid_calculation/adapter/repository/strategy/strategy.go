package strategy

import (
	"gitlab.jooble.com/marketing_tech/yandex_bidder/domain/entities"
	sqlStore "gitlab.jooble.com/marketing_tech/yandex_bidder/infrastructure/store/sql"
	"gitlab.jooble.com/marketing_tech/yandex_bidder/usecase"
)

type (
	repo struct {
		store *sqlStore.Store
	}
)

func New(store *sqlStore.Store) usecase.StrategyRepo {
	return &repo{
		store: store,
	}
}

func (r *repo) GetAll() ([]entities.Strategy, error) {
	rows, err := r.store.DB.
		Select("name").
		From("strategies").
		Rows()
	if err != nil {
		return nil, err
	}

	strategies := make([]entities.Strategy, 0)
	for rows.Next() {
		var strategy entities.Strategy
		if err := rows.Scan(&strategy); err != nil {
			return nil, err
		}
		strategies = append(strategies, strategy)
	}

	return strategies, nil
}
