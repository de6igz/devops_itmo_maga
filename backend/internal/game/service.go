package game

type Repository interface {
	Create(Game) (Game, error)
	List(Filters) ([]Game, error)
	GetByID(int64) (Game, error)
	Update(int64, Game) (Game, error)
	Delete(int64) error
}

type Service struct {
	repository Repository
}

func NewService(repository Repository) *Service {
	return &Service{repository: repository}
}

func (s *Service) List(filters Filters) ([]Game, error) {
	return s.repository.List(filters)
}

func (s *Service) GetByID(id int64) (Game, error) {
	return s.repository.GetByID(id)
}

func (s *Service) Create(game Game) (Game, error) {
	Normalize(&game)
	if err := Validate(game); err != nil {
		return Game{}, err
	}

	return s.repository.Create(game)
}

func (s *Service) Update(id int64, game Game) (Game, error) {
	Normalize(&game)
	if err := Validate(game); err != nil {
		return Game{}, err
	}

	return s.repository.Update(id, game)
}

func (s *Service) Delete(id int64) error {
	return s.repository.Delete(id)
}
