package services

import (
	"context"
	"errors"
	"fmt"

	"github.com/Romasmi/social-network/internal/models"
	"github.com/Romasmi/social-network/internal/repository"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	cityRepo           *repository.CityRepository
	profileRepo        *repository.ProfileRepository
	profileFriendsRepo *repository.ProfileFriendsRepository
	uow                repository.UnitOfWork
}

func CreateUserService(
	cityRepo *repository.CityRepository,
	profileRepo *repository.ProfileRepository,
	profileFriendRepo *repository.ProfileFriendsRepository,
	uow repository.UnitOfWork) *UserService {
	return &UserService{cityRepo: cityRepo, profileRepo: profileRepo, profileFriendsRepo: profileFriendRepo, uow: uow}
}

func (s *UserService) RegisterUser(ctx context.Context, payload *models.CreateProfileModel) (*models.Profile, error) {
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
		if !errors.Is(err, repository.ErrNotFound) {
			return nil, err
		}
		cityId, err := uuid.NewV7()
		if err != nil {
			return nil, err
		}

		city, err = s.cityRepo.CreateCity(ctx, &models.City{
			ID:   cityId,
			Name: payload.City,
		})
		if err != nil {
			return nil, err
		}

	}

	var newUser models.User
	newUser.ID = userId
	newUser.Email = gofakeit.Email() // Generate email just as example
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
	newProfile.City = city.Name

	err = s.uow.WithTransaction(ctx, func(ctx context.Context, txUoW repository.UnitOfWork) error {
		userRepo := txUoW.User()
		profileRepo := txUoW.Profile()

		_, err = userRepo.CreateUser(ctx, &newUser)
		if err != nil {
			return fmt.Errorf("can't create user: %w", err)
		}

		profile, err := profileRepo.CreateProfile(ctx, &newProfile)
		if err != nil {
			return fmt.Errorf("can't create profile: %w", err)
		}
		profile.City = city.Name
		return nil
	})
	return &newProfile, err
}

func (s *UserService) GetUserByProfileId(ctx context.Context, profileId uuid.UUID) (*models.Profile, error) {
	return s.profileRepo.GetProfileByProfileId(ctx, profileId)
}

func (s *UserService) SearchUsers(ctx context.Context, queryParams *models.UserSearchParams) ([]*models.Profile, error) {
	return s.profileRepo.SearchProfile(ctx, queryParams)
}

func (s *UserService) SetFriend(ctx context.Context, profileId1, profileId2 uuid.UUID) error {
	return s.profileFriendsRepo.SetFriend(ctx, profileId1, profileId2)
}
