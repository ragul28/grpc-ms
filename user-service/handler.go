package main

import (
	"context"
	"errors"
	"fmt"
	"log"

	pb "github.com/grpc-ms/user-service/proto/user"
	"golang.org/x/crypto/bcrypt"
)

type handler struct {
	repo         repository
	tokenService Authable
}

func (s *handler) Get(ctx context.Context, req *pb.User) (*pb.Response, error) {
	user, err := s.repo.Get(req.Id)
	if err != nil {
		return nil, err
	}

	return &pb.Response{User: user}, nil
}

func (s *handler) GetAll(ctx context.Context, req *pb.Request) (*pb.Response, error) {
	users, err := s.repo.GetAll()
	if err != nil {
		return nil, err
	}

	return &pb.Response{Users: users}, nil
}

func (s *handler) Auth(ctx context.Context, req *pb.User) (*pb.Token, error) {
	log.Println("Logging in with:", req.Email, req.Password)
	user, err := s.repo.GetByEmail(req.Email)
	log.Println(user)
	if err != nil {
		return nil, err
	}

	// comapares input pwd with db hashed pwd
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, err
	}

	token, err := s.tokenService.Encode(user)
	if err != nil {
		return nil, err
	}

	return &pb.Token{Token: token}, nil
}

func (s *handler) Create(ctx context.Context, req *pb.User) (*pb.Response, error) {

	log.Println("Creating user:", req)

	hashedPass, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("error hashing password: %v", err))
	}

	req.Password = string(hashedPass)
	if err := s.repo.Create(req); err != nil {
		return nil, errors.New(fmt.Sprintf("error creating user: %v", err))
	}

	token, err := s.tokenService.Encode(req)
	if err != nil {
		return nil, err
	}

	return &pb.Response{User: req, Token: &pb.Token{Token: token}}, nil
}

func (s *handler) ValidateToken(ctx context.Context, req *pb.Token) (*pb.Token, error) {

	claims, err := s.tokenService.Decode(req.Token)
	if err != nil {
		return nil, err
	}

	if claims.User.Id == "" {
		return nil, errors.New("Invalid user")
	}

	return &pb.Token{Valid: true}, nil
}
