package usecases

import "github.com/Hailemari/clean_architecture_task_manager/Domain"

type UserUseCase struct {
    repo domain.UserRepository
}

func NewUserUseCase(repo domain.UserRepository) *UserUseCase {
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
    return uc.repo.PromoteUser(username)
}