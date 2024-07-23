package service

import (
	"context"
	"fmt"
	"grpc-redis-postgres/internal/db"
	"grpc-redis-postgres/internal/redis"
	"strconv"

	"grpc-redis-postgres/proto"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// userService implements the proto.UserServiceServer interface
type userService struct {
	proto.UnimplementedUserServiceServer
	redisClient redis.RedisClient
	dbClient    db.DatabaseClient
}

func NewUserService(dbClient db.DatabaseClient, redisClient redis.RedisClient) *userService {
	return &userService{
		redisClient: redisClient,
		dbClient:    dbClient,
	}
}

// GetUser retrieves a user from Redis or PostgreSQL.
func (s *userService) GetUser(ctx context.Context, req *proto.GetUserRequest) (*proto.User, error) {
	// return nil, status.Errorf(codes.NotFound, "user with ID %s not found", req.Id)

	// 1. Try to retrieve user from Redis
	key := "user_" + strconv.FormatInt(req.Id, 10)
	user, err := s.redisClient.Get(ctx, key)
	if err == nil {
		id, _ := strconv.ParseInt(user["id"], 10, 64)
		var user = &proto.User{
			Id:    id,
			Name:  user["name"],
			Email: user["email"],
		}
		return user, nil
	}else{
		return nil, status.Errorf(codes.Internal, "error caching user in Redis: %v", err)
	}

	// 2. If not found in Redis, retrieve from PostgreSQL
	// user, err = s.dbClient.GetUser(ctx, req.Id)
	// if err != nil {
	// 	return nil, status.Errorf(codes.NotFound, "user with ID %s not found", req.Id)
	// }

	// // 3. Cache user in Redis
	// if err := s.redisClient.Set(ctx, "user", user, time.Minute); err != nil {
	// 	// Error caching, but still return user
	// 	log.Printf("Error caching user in Redis: %v", err)
	// }

	// return user, nil
}

// CreateUser creates a new user and stores it in Redis and PostgreSQL.
func (s *userService) CreateUser(ctx context.Context, req *proto.CreateUserRequest) (*proto.User, error) {
	// 1. Create user in PostgreSQL
	var newUser = &proto.User{
		Name:  req.Name,
		Email: req.Email,
	}
	user, err := s.dbClient.CreateUser(newUser)
	if err != nil {
		fmt.Println("err:",err)
		return nil, status.Errorf(codes.Internal, "error creating user: %v", err)
	}

	// 2. Cache user in Redis
	if err := s.redisClient.Set(ctx, "user_"+strconv.FormatInt(user.Id, 10), user); err != nil {
		return nil, status.Errorf(codes.Internal, "error caching user in Redis: %v", err)
	}

	return user, nil
}
