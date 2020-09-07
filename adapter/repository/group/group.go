package group

import (
	"database/sql"

	dbx "github.com/go-ozzo/ozzo-dbx"

	"gitlab.jooble.com/marketing_tech/yandex_bidder/domain"
	"gitlab.jooble.com/marketing_tech/yandex_bidder/domain/entities"
	sqlStore "gitlab.jooble.com/marketing_tech/yandex_bidder/infrastructure/store/sql"
	"gitlab.jooble.com/marketing_tech/yandex_bidder/usecase"
)

type (
	repo struct {
		store *sqlStore.Store
	}
)

func New(store *sqlStore.Store) usecase.GroupRepo {
	return &repo{
		store: store,
	}
}

func (r *repo) Accounts(group *entities.Group) ([]*entities.Account, error) {
	accounts := make([]*entities.Account, 0)
	rows, err := r.store.DB.
		Select("id", "name").
		From("accounts AS a").
		Where(dbx.HashExp{"group_id": group.ID}).
		Rows()
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		account := new(entities.Account)
		if err := rows.ScanStruct(account); err != nil {
			return nil, err
		}
		accounts = append(accounts, account)
	}

	return accounts, nil
}

func (r *repo) GetAll() ([]*entities.Group, error) {
	groups := make([]*entities.Group, 0)
	rows, err := r.store.DB.
		Select("g.id", "g.name", "schedule_start", "schedule_interval", "s.name as strategy").
		From("groups AS g").
		LeftJoin("strategies AS s", dbx.NewExp("g.strategy_id = s.id")).
		Rows()
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var group *entities.Group
		if err := rows.ScanStruct(&group); err != nil {
			return nil, err
		}
		groups = append(groups, group)
	}

	return groups, nil
}

func (r *repo) GetByID(id int) (*entities.Group, error) {
	group := new(entities.Group)
	query := r.store.DB.
		Select("g.id", "g.name", "schedule_start", "schedule_interval", "s.name AS strategy").
		From("groups AS g").
		LeftJoin("strategies AS s", dbx.NewExp("g.strategy_id = s.id")).
		Where(dbx.HashExp{"g.id": id})

	if err := query.One(group); err != nil {
		if err == sql.ErrNoRows {
			err = domain.ErrGroupNotFound
		}
		return nil, err
	}

	return group, nil
}

func (r *repo) Update(group *entities.Group) error {
	fieldsToUpdate := dbx.Params{
		"schedule_start":    group.Start,
		"schedule_interval": group.Interval,
		"strategy":          group.Strategy,
	}
	whereClause := dbx.HashExp{"id": group.ID}

	query := r.store.DB.Update("groups", fieldsToUpdate, whereClause)
	if _, err := query.Execute(); err != nil {
		return err
	}

	return nil
}
