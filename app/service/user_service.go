package service

import (
	"errors"
	"github.com/Lim018/arek-muhammadiyah-be/app/model"
	"github.com/Lim018/arek-muhammadiyah-be/app/repository"
	"github.com/Lim018/arek-muhammadiyah-be/helper"

	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	userRepo *repository.UserRepository
}

func NewUserService() *UserService {
	return &UserService{
		userRepo: repository.NewUserRepository(),
	}
}

func (s *UserService) GetAllUsers(page, limit int) ([]model.User, model.Pagination, error) {
	offset := (page - 1) * limit
	users, total, err := s.userRepo.GetAll(limit, offset)
	if err != nil {
		return nil, model.Pagination{}, err
	}

	pagination := helper.CreatePagination(int64(page), int64(limit), total)
	return users, pagination, nil
}

func (s *UserService) GetUserByID(id string) (*model.User, error) {
	return s.userRepo.GetByID(id)
}

func (s *UserService) CreateUser(req *model.CreateUserRequest) (*model.User, error) {
	// Check if user already exists
	existing, _ := s.userRepo.GetByID(req.ID)
	if existing != nil {
		return nil, errors.New("user already exists")
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &model.User{
		ID:        req.ID,
		Name:      req.Name,
		Password:  string(hashedPassword),
		RoleID:    req.RoleID,
		VillageID: req.VillageID,
		NIK:       req.NIK,
		Address:   req.Address,
	}

	err = s.userRepo.Create(user)
	if err != nil {
		return nil, err
	}

	return s.userRepo.GetByID(user.ID)
}

func (s *UserService) UpdateUser(id string, req *model.UpdateUserRequest) (*model.User, error) {
	// Check if user exists
	existing, err := s.userRepo.GetByID(id)
	if err != nil {
		return nil, errors.New("user not found")
	}

	updateData := &model.User{
		Name:      helper.GetStringValue(req.Name, existing.Name),
		RoleID:    helper.GetUintPointer(req.RoleID, existing.RoleID),
		VillageID: helper.GetUintPointer(req.VillageID, existing.VillageID),
		NIK:       helper.GetStringPointer(req.NIK, existing.NIK),
		Address:   helper.GetStringPointer(req.Address, existing.Address),
		Photo:     helper.GetStringPointer(req.Photo, existing.Photo),
	}

	err = s.userRepo.Update(id, updateData)
	if err != nil {
		return nil, err
	}

	return s.userRepo.GetByID(id)
}

func (s *UserService) DeleteUser(id string) error {
	// Check if user exists
	_, err := s.userRepo.GetByID(id)
	if err != nil {
		return errors.New("user not found")
	}

	return s.userRepo.Delete(id)
}

func (s *UserService) Login(req *model.LoginRequest) (*model.User, error) {
	user, err := s.userRepo.GetByID(req.ID)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	return user, nil
}

func (s *UserService) BulkCreateUsers(users []model.CreateUserRequest) ([]model.User, error) {
	var createdUsers []model.User
	var failedUsers []string

	for _, userReq := range users {
		user, err := s.CreateUser(&userReq)
		if err != nil {
			failedUsers = append(failedUsers, userReq.ID)
			continue
		}
		createdUsers = append(createdUsers, *user)
	}

	if len(failedUsers) > 0 {
		return createdUsers, errors.New("some users failed to create")
	}

	return createdUsers, nil
}