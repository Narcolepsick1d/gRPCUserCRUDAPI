package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	_ "github.com/lib/pq"
	"google.golang.org/grpc"
	"log"
	"net"
	pb "testHamkor/protos"
)

func init() {
	DatabaseConnection()
}

var DB *sql.DB
var err error

type User struct {
	ID    int64 `gorm:"primarykey"`
	name  string
	age   int64
	phone string
}

func DatabaseConnection() {
	host := "localhost"
	port := "5432"
	dbName := "testP"
	dbUser := "postgres"
	password := "12345"
	dsn := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=disable",
		host,
		port,
		dbUser,
		dbName,
		password,
	)
	DB, err = sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal("Error connecting to the database...", err)
	}
	fmt.Println("Database connection successful...")
}

var (
	port = flag.Int("port", 8080, "gRPC server port")
)

type server struct {
	pb.UnimplementedUserServiceServer
}

func (*server) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	fmt.Println("Create User")
	user := req.GetUser()

	data := User{

		name:  user.GetName(),
		age:   user.GetAge(),
		phone: user.GetPhone(),
	}

	_, err := DB.Exec("insert into users (name, age, phone) values ($1,$2,$3)",
		&data.name, &data.age, &data.phone)
	if err != nil {
		return nil, err
	}
	return &pb.CreateUserResponse{
		User: &pb.User{
			Id:    user.GetId(),
			Name:  user.GetName(),
			Age:   user.GetAge(),
			Phone: user.GetPhone(),
		},
	}, nil
}

func (*server) GetUser(ctx context.Context, req *pb.ReadUserRequest) (*pb.ReadUserResponse, error) {
	fmt.Println("Read user", req.GetId())
	var user User
	err := DB.QueryRow("select id,name,age,phone from users WHERE id=$1", req.GetId()).Scan(
		&user.ID, &user.name, &user.age, &user.phone)
	if err == sql.ErrNoRows {
		return nil, err
	}
	return &pb.ReadUserResponse{
		User: &pb.User{
			Id:    user.ID,
			Name:  user.name,
			Age:   user.age,
			Phone: user.phone,
		},
	}, nil
}

func (*server) GetUsers(ctx context.Context, req *pb.ReadUserRequest) (*pb.ReadUsersResponse, error) {
	fmt.Println("Read Users")
	users := []*pb.User{}
	rows, err := DB.Query("select id,name,age,phone from users")
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var user *pb.User = new(pb.User)
		err := rows.Scan(&user.Id, &user.Name, &user.Age, &user.Phone)
		log.Println(err)

		users = append(users, user)
	}
	return &pb.ReadUsersResponse{
		User: users,
	}, nil
}

func main() {
	fmt.Println("gRPC server running ...")

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))

	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := grpc.NewServer()

	pb.RegisterUserServiceServer(s, &server{})

	log.Printf("Server listening at %v", lis.Addr())

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve : %v", err)
	}
}
