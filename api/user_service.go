package api

import "github.com/simple-jwt-auth/models"

type UserService struct {
	UserRepository models.UserRepository
}

func ProvideUserService(repo models.UserRepository) UserService {
	return UserService{UserRepository: repo}

}

func (s *UserService) FindAll() ([]models.User, error) {
	return s.UserRepository.FindAll()
}

func (s *UserService) FindByID(id int) (models.User, error) {
	return s.UserRepository.FindByID(id)
}

func (s *UserService) Create(user models.User) (models.User, error) {
	return s.UserRepository.Create(user)
}

func (s *UserService) Update(user models.User) (models.User, error) {
	return s.UserRepository.Update(user)
}

func (s *UserService) Delete(user models.User) (models.User, error) {
	return s.UserRepository.Delete(user)

}
