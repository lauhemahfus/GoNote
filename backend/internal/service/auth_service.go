package service

import (
    "errors"
    "gonote/internal/domain"
    "gonote/internal/utils"
    "golang.org/x/crypto/bcrypt"
)

type AuthService interface {
    SignUp(email, password, name string) (*domain.User, error)
    Login(email, password string) (string, error)
}

type authService struct {
    userRepo domain.UserRepository
    jwtUtil  *utils.JWTUtil
}

func NewAuthService(userRepo domain.UserRepository, jwtUtil *utils.JWTUtil) AuthService {
    return &authService{
        userRepo: userRepo,
        jwtUtil:  jwtUtil,
    }
}

func (s *authService) SignUp(email, password, name string) (*domain.User, error) {
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
    if err != nil {
        return nil, err
    }
    
    user := &domain.User{
        Email:    email,
        Password: string(hashedPassword),
        Name:     name,
    }
    
    if err := s.userRepo.Create(user); err != nil {
        return nil, err
    }
    
    return user, nil
}

func (s *authService) Login(email, password string) (string, error) {
    user, err := s.userRepo.GetByEmail(email)
    if err != nil {
        return "", errors.New("invalid credentials")
    }
    
    if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
        return "", errors.New("invalid credentials")
    }
    
    token, err := s.jwtUtil.GenerateToken(user.ID)
    if err != nil {
        return "", err
    }
    
    return token, nil
}