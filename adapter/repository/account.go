package repository

import (
	dbx "github.com/go-ozzo/ozzo-dbx"

	"gitlab.jooble.com/marketing_tech/yandex_bidder/domain"
	"gitlab.jooble.com/marketing_tech/yandex_bidder/usecase"
)

type (
	accountRepo struct {
		db dbx.Builder
	}
)

func NewAccountRepo(db dbx.Builder) usecase.AccountRepo {
	return &accountRepo{
		db: db,
	}
}

func (r *accountRepo) Campaigns(account *domain.Account) ([]*domain.Campaign, error) {
	campaigns := make([]*domain.Campaign, 0)
	rows, err := r.db.
		Select("c.id", "c.name").
		From("campaigns AS c").
		InnerJoin("accounts AS a", dbx.NewExp("a.id = c.account_id")).
		Where(dbx.HashExp{"a.id": account.ID}).
		Rows()
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		campaign := new(domain.Campaign)
		if err := rows.ScanStruct(campaign); err != nil {
			return nil, err
		}
		campaigns = append(campaigns, campaign)
	}

	return campaigns, nil
}
