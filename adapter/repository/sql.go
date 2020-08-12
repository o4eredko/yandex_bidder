package repository

type (
	SQLRepository struct {
	}
)

func NewSQLRepository() *SQLRepository {
	return &SQLRepository{}
}

func (r *SQLRepository) GetGroups() {
	//sql = "SELECT DISTINCT name FROM Yandex.dbo.accounts"
	return
}
