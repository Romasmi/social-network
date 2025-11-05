package services

import (
	"context"

	"github.com/Romasmi/social-network/internal/models"
	"github.com/Romasmi/social-network/internal/repository"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	cityRepo *repository.CityRepository
	uow      repository.UnitOfWork
}

func CreateUserService(cityRepo *repository.CityRepository, uow repository.UnitOfWork) *UserService {
	return &UserService{cityRepo: cityRepo, uow: uow}
}

func (s *UserService) RegisterUser(ctx context.Context, payload *models.CreateUserModel) (*models.Profile, error) {
	userId, err := uuid.NewV7()
	if err != nil {
		return nil, err
	}

	profileId, err := uuid.NewV7()
	if err != nil {
		return nil, err
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(payload.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	city, err := s.cityRepo.GetCityByName(ctx, payload.City)
	if err != nil {
		return nil, err
	}
	var newUser models.User
	newUser.ID = userId
	newUser.Email = ""
	newUser.IsActive = true
	newUser.PasswordHash = string(passwordHash)

	var newProfile models.Profile
	newProfile.ID = profileId
	newProfile.UserId = userId
	newProfile.FirstName = payload.FirstName
	newProfile.SecondName = payload.SecondName
	newProfile.Birthdate = payload.Birthdate
	newProfile.Gender = payload.Gender
	newProfile.Biography = payload.Biography
	newProfile.CityId = city.ID

	userRepo := s.uow.User()
	profileRepo := s.uow.Profile()

	err = s.uow.WithTransaction(ctx, func(ctx context.Context) error {
		_, err = userRepo.CreateUser(ctx, &newUser)
		if err != nil {
			return err
		}

		profile, err := profileRepo.CreateProfile(ctx, &newProfile)
		if err != nil {
			return err
		}
		profile.City = city.Name
		return nil
	})
	return &newProfile, err
}
