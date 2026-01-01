package services

import (
	"context"
	"errors"
	"fmt"

	"github.com/Romasmi/social-network/internal/domain/city"
	"github.com/Romasmi/social-network/internal/domain/profile"
	"github.com/Romasmi/social-network/internal/domain/user"
	"github.com/Romasmi/social-network/internal/events/publisher"
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
	publisher          publisher.Publisher
}

func CreateUserService(
	cityRepo *repository.CityRepository,
	profileRepo *repository.ProfileRepository,
	profileFriendRepo *repository.ProfileFriendsRepository,
	uow repository.UnitOfWork,
	publisher publisher.Publisher,
) *UserService {
	return &UserService{cityRepo: cityRepo, profileRepo: profileRepo, profileFriendsRepo: profileFriendRepo, uow: uow, publisher: publisher}
}

func (s *UserService) RegisterUser(ctx context.Context, payload *profile.CreateProfileModel) (*profile.Profile, error) {
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

	cityModel, err := s.cityRepo.GetCityByName(ctx, payload.City)
	if err != nil {
		if !errors.Is(err, repository.ErrNotFound) {
			return nil, err
		}
		cityId, err := uuid.NewV7()
		if err != nil {
			return nil, err
		}

		cityModel, err = s.cityRepo.CreateCity(ctx, &city.City{
			ID:   cityId,
			Name: payload.City,
		})
		if err != nil {
			return nil, err
		}

	}

	var newUser user.User
	newUser.ID = userId
	newUser.Email = gofakeit.Email() // Generate email just as example
	newUser.IsActive = true
	newUser.PasswordHash = string(passwordHash)

	var newProfile profile.Profile
	newProfile.ID = profileId
	newProfile.UserId = userId
	newProfile.FirstName = payload.FirstName
	newProfile.SecondName = payload.SecondName
	newProfile.Birthdate = payload.Birthdate
	newProfile.Gender = payload.Gender
	newProfile.Biography = payload.Biography
	newProfile.CityId = cityModel.ID
	newProfile.City = cityModel.Name

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
		profile.City = cityModel.Name
		return nil
	})
	return &newProfile, err
}

func (s *UserService) GetUserByProfileId(ctx context.Context, profileId uuid.UUID) (*profile.Profile, error) {
	return s.profileRepo.GetProfileByProfileId(ctx, profileId)
}

func (s *UserService) SearchUsers(ctx context.Context, queryParams *profile.UserSearchParams) ([]*profile.Profile, error) {
	return s.profileRepo.SearchProfile(ctx, queryParams)
}

func (s *UserService) SetFriend(ctx context.Context, profileId1, profileId2 uuid.UUID) error {
	if err := s.profileFriendsRepo.SetFriend(ctx, profileId1, profileId2); err != nil {
		return err
	}
	return s.publisher.Publish(ctx, profile.NewFriendAddedEvent(&profile.FriendAddedEventData{
		ProfileID: profileId1,
		FriendID:  profileId2,
	}))
}

func (s *UserService) DeleteFriend(ctx context.Context, profileId1, profileId2 uuid.UUID) error {
	if s.profileFriendsRepo.DeleteFriend(ctx, profileId1, profileId2) != nil {
		return s.profileFriendsRepo.DeleteFriend(ctx, profileId1, profileId2)
	}
	return s.publisher.Publish(ctx, profile.NewFriendDeletedEvent(&profile.FriendDeletedEventData{
		ProfileID: profileId1,
		FriendID:  profileId2,
	}))
}
