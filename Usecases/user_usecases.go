package usecases

import (
    "errors"
    "github.com/Hailemari/clean_architecture_task_manager/Domain"
)

type UserUseCase struct {
    repo domain.UserRepository
}

func NewUserUseCase(repo domain.UserRepository) domain.UserUseCaseInterface {
    return &UserUseCase{repo: repo}
}

func (uc *UserUseCase) CreateUser(user *domain.User) error {
    if err := user.ValidateUser(); err != nil {
        return err
    }
    if err := user.HashPassword(); err != nil {
        return err
    }
    return uc.repo.CreateUser(user)
}

func (uc *UserUseCase) GetUserByUsername(username string) (*domain.User, error) {
    return uc.repo.GetUserByUsername(username)
}

func (uc *UserUseCase) PromoteUser(username string) error {
    user, err := uc.repo.GetUserByUsername(username)
    if err != nil {
        return err
    }
    if user.Role == "admin" {
        return errors.New("user is already an admin")
    }
    return uc.repo.PromoteUser(username)
}
