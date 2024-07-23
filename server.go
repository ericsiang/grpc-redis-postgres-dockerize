package main

import (
	"context"
	"grpc-redis-postgres/internal/db"
	"grpc-redis-postgres/internal/redis"
	"grpc-redis-postgres/internal/service"
	"grpc-redis-postgres/proto"
	"log"
	"net"
	"os"
	"strconv"
	"time"

	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/recovery"
	// grpc_validator "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/validator"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func InterceptorLogger(l *zap.Logger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		start := time.Now()
		resp, err := handler(ctx, req)
		duration := time.Since(start)

		if err != nil {
			l.Error("RPC failed",
				zap.String("method", info.FullMethod),
				zap.Duration("duration", duration),
				zap.Error(err),
			)
		} else {
			l.Info("RPC succeeded",
				zap.String("method", info.FullMethod),
				zap.Duration("duration", duration),
			)
		}

		return resp, err
	}
}

func main() {
	// Load environment variables from .env file
	// local 測試，使用 docker compose 要隱藏此段
	// if err := godotenv.Load(".env"); err != nil {
	// 	log.Fatal("Error loading .env file: ", err)
	// }

	// Initialize logger
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	// Configure Redis
	redisAddr := os.Getenv("REDIS_ADDR")
	redisPassword := os.Getenv("REDIS_PASSWORD")
	redisDB, err := strconv.Atoi(os.Getenv("REDIS_DB"))
	if err != nil {
		log.Fatalf("Error parsing REDIS_DB: %v", err)
	}
	redisClient, err := redis.NewRedisClient(redisAddr, redisPassword, redisDB)
	if err != nil {
		log.Fatalf("Error creating Redis client: %v", err)
	}

	// // Configure PostgreSQL
	dbDSN := os.Getenv("DB_DSN")
	dbClient, err := db.NewDatabaseClient(dbDSN)
	if err != nil {
		log.Fatalf("Error creating PostgreSQL client: %v", err)
	}

	userService := service.NewUserService(dbClient, redisClient)

	// Create gRPC server with interceptors
	server := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			// grpc_validator.UnaryServerInterceptor(),
			grpc_recovery.UnaryServerInterceptor(),
			InterceptorLogger(logger),
			// tags.UnaryServerInterceptor(),
		),
	)
	reflection.Register(server)
	// Register User service
	proto.RegisterUserServiceServer(server, userService)

	port := os.Getenv("APP_PORT")
	// Listen on port
	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	log.Printf("Server listening on port: %s", ":"+port)

	// Serve gRPC requests
	if err := server.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
