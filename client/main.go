package main

import (
	"flag"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"net/http"
	"strconv"
	pb "testHamkor/protos"
)

var (
	addr = flag.String("addr", "localhost:8080", "the address to connect to")
)

type User struct {
	ID    int64  `json:"id"`
	Name  string `json:"name"`
	Age   int64  `json:"age"`
	Phone string `json:"phone"`
}

func main() {
	flag.Parse()
	conn, err := grpc.Dial(*addr, grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}

	defer conn.Close()
	client := pb.NewUserServiceClient(conn)

	r := gin.Default()
	r.GET("/users", func(ctx *gin.Context) {
		res, err := client.GetUsers(ctx, &pb.ReadUserRequest{})
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": err,
			})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{
			"user": res.User,
		})
	})
	r.GET("/user/:id", func(ctx *gin.Context) {
		id := ctx.Param("id")
		i, _ := strconv.ParseInt(id, 10, 64)
		res, err := client.GetUser(ctx, &pb.ReadUserRequest{Id: i})
		if err != nil {
			ctx.JSON(http.StatusNotFound, gin.H{
				"message": err.Error(),
			})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{
			"user": res.User,
		})
	})
	r.POST("/user", func(ctx *gin.Context) {
		var user User

		err := ctx.ShouldBind(&user)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": err,
			})
			return
		}
		data := &pb.User{}
		res, err := client.CreateUser(ctx, &pb.CreateUserRequest{
			User: data,
		})
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": err,
			})
			return
		}
		ctx.JSON(http.StatusCreated, gin.H{
			"user": res.User,
		})
	})

	r.Run(":8081")

}
