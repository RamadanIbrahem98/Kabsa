package kabsa

import "github.com/RamadanIbrahem98/kabsa/db"

type Kabsa struct {
	Presses int64
	StartAt int64
	EndAt   int64
	DB      *db.DB
}

func New() (*Kabsa, error) {
	db, err := db.New()

	if err != nil {
		return nil, err
	}

	return &Kabsa{
		Presses: 0,
		StartAt: 0,
		EndAt:   0,
		DB:      db,
	}, nil
}
