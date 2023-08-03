package admin

import (
	"ca-reservaksin/app/middlewares"
	"ca-reservaksin/businesses"
	"ca-reservaksin/helpers/encrypt"
	"ca-reservaksin/helpers/nanoid"
	"strings"
)

type adminService struct {
	adminRepository Repository
	jwtAuth         *middlewares.ConfigJWT
}

func NewAdminService(adminRepo Repository, jwtAuth *middlewares.ConfigJWT) Service {
	return &adminService{
		adminRepository: adminRepo,
		jwtAuth:         jwtAuth,
	}
}

func (service *adminService) Register(dataAdmin *Domain) (Domain, error) {
	existedAdmin, err := service.adminRepository.GetByUsername(dataAdmin.Username)
	if err != nil {
		if !strings.Contains(err.Error(), "not found") {
			return Domain{}, businesses.ErrInternalServer
		}
	}

	if existedAdmin != (Domain{}) {
		return Domain{}, businesses.ErrDuplicateData
	}

	dataAdmin.Id, err = nanoid.GenerateNanoId()
	if err != nil {
		return Domain{}, businesses.ErrInternalServer
	}

	hashedPassword := encrypt.HashAndSalt([]byte(dataAdmin.Password))
	dataAdmin.Password = hashedPassword
	res, err := service.adminRepository.Register(dataAdmin)
	if err != nil {
		return Domain{}, businesses.ErrInternalServer
	}
	return res, nil
}

func (service *adminService) Login(username, password string) (string, error) {
	adminDomain, err := service.adminRepository.GetByUsername(username)
	if err != nil {
		if !strings.Contains(err.Error(), "not found") {
			return "", businesses.ErrInternalServer
		}
		return "", businesses.ErrUsernamePasswordNotFound
	}

	if !encrypt.ValidateHash(password, adminDomain.Password) {
		return "", businesses.ErrUsernamePasswordNotFound
	}

	token := service.jwtAuth.GenerateTokenAdmin(adminDomain.Id)
	return token, nil
}

func (service *adminService) GetByID(id string) (Domain, error) {
	adminDomain, err := service.adminRepository.GetByID(id)
	if err != nil {
		return Domain{}, businesses.ErrIDNotFound
	}

	return adminDomain, nil
}
