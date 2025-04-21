package repository

import (
	"github.com/nurbol/cinema/internal/models"
	"gorm.io/gorm"
)

type CinemaRepository struct {
	db *gorm.DB
}

func NewCinemaRepository(db *gorm.DB) *CinemaRepository {
	return &CinemaRepository{db: db}
}

func (r *CinemaRepository) GetAll() ([]models.Cinema, error) {
	var cinemas []models.Cinema
	err := r.db.Find(&cinemas).Error
	return cinemas, err
}

func (r *CinemaRepository) GetByID(id int) (*models.Cinema, error) {
	var cinema models.Cinema
	err := r.db.First(&cinema, id).Error
	return &cinema, err
}

func (r *CinemaRepository) Create(cinema *models.Cinema) error {
	return r.db.Create(cinema).Error
}

func (r *CinemaRepository) Update(id int, cinema *models.Cinema) error {
	return r.db.Model(&models.Cinema{}).Where("id = ?", id).Updates(cinema).Error
}

func (r *CinemaRepository) Delete(id int) error {
	return r.db.Delete(&models.Cinema{}, id).Error
}
