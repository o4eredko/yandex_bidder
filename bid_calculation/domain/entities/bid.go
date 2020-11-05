package entities

type (
	Bid struct {
		CampaignID int `json:"CampaignId" db:"id"`
		Bid        *int
		OldBid     *int `json:"-"`
	}
)
