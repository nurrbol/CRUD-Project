package services

import (
	"github.com/nurbol/cinema/internal/models"
	"github.com/nurbol/cinema/internal/repository"
)

type CinemaService struct {
	repo *repository.CinemaRepository
}

func NewCinemaService(repo *repository.CinemaRepository) *CinemaService {
	return &CinemaService{repo: repo}
}

func (s *CinemaService) GetAllCinemas() ([]models.Cinema, error) {
	return s.repo.GetAll()
}

func (s *CinemaService) GetCinemaByID(id int) (*models.Cinema, error) {
	return s.repo.GetByID(id)
}

func (s *CinemaService) CreateCinema(cinema *models.Cinema) error {
	return s.repo.Create(cinema)
}

func (s *CinemaService) UpdateCinema(id int, cinema *models.Cinema) error {
	return s.repo.Update(id, cinema)
}

func (s *CinemaService) DeleteCinema(id int) error {
	return s.repo.Delete(id)
}
