package user_handler

import "github.com/Romasmi/social-network/internal/services"

type UserHandler struct {
	userService *services.UserService
}

func CreateUserHandler(userService *services.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}
