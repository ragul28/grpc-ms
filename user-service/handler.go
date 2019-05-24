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

func (s *handler) Get(ctx context.Context, req *pb.User, res *pb.Response) error {
	user, err := s.repo.Get(req.Id)
	if err != nil {
		return err
	}
	res.User = user
	return nil
}

func (s *handler) GetAll(ctx context.Context, req *pb.Request, res *pb.Response) error {
	users, err := s.repo.GetAll()
	if err != nil {
		return err
	}
	res.Users = users
	return nil
}

func (s *handler) Auth(ctx context.Context, req *pb.User, res *pb.Token) error {
	log.Println("Logging in with:", req.Email, req.Password)
	user, err := s.repo.GetByEmail(req.Email)
	log.Println(user)
	if err != nil {
		return err
	}

	// comapares input pwd with db hashed pwd
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return err
	}

	token, err := s.tokenService.Encode(user)
	if err != nil {
		return err
	}
	res.Token = token
	return nil
}

func (s *handler) Create(ctx context.Context, req *pb.User, res *pb.Response) error {

	log.Println("Creating user:", req)

	hashedPass, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return errors.New(fmt.Sprintf("error hashing password: %v", err))
	}

	req.Password = string(hashedPass)
	if err := s.repo.Create(req); err != nil {
		return errors.New(fmt.Sprintf("error creating user: %v", err))
	}

	token, err := s.tokenService.Encode(req)
	if err != nil {
		return err
	}

	res.User = req
	res.Token = &pb.Token{Token: token}
	return nil
}

func (s *handler) ValidateToken(ctx context.Context, req *pb.Token, res *pb.Token) error {

	claims, err := s.tokenService.Decode(req.Token)
	if err != nil {
		return err
	}

	if claims.User.Id == "" {
		return errors.New("Invalid user")
	}

	res.Valid = true
	return nil
}
