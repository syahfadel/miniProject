package controller

import (
	"fmt"
	"miniProject/entities"
	"miniProject/helpers"
	"miniProject/services"
	"net/http"
	"strconv"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type SosmedController struct {
	DB            *gorm.DB
	SosmedService *services.SosmedService
}

type UserDto struct {
	UserName string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type PostDto struct {
	Body string `json:"body"`
}

func (sc *SosmedController) Register(ctx *gin.Context) {
	var userDto UserDto

	if err := ctx.ShouldBindJSON(&userDto); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	user := entities.User{
		UserName: userDto.UserName,
		Email:    userDto.Email,
		Password: userDto.Password,
	}

	result, err := sc.SosmedService.CreateUser(user)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, result)
}

func (sc *SosmedController) Login(ctx *gin.Context) {
	var userDto UserDto

	if err := ctx.ShouldBindJSON(&userDto); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	user := entities.User{
		UserName: userDto.UserName,
		Email:    userDto.Email,
	}

	result, err := sc.SosmedService.Login(user)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error":   "Unauthorized",
			"message": "invalid email/password",
		})
		return
	}

	comparePass := helpers.ComparePass([]byte(result.Password), []byte(userDto.Password))

	if !comparePass {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error":   "Unauthorized",
			"message": "invalid email/password",
		})
		return
	}

	token := helpers.GenerateToken(result.ID, result.Email, result.UserName)

	ctx.JSON(http.StatusOK, gin.H{
		"token": token,
	})
}

func (sc *SosmedController) Follow(ctx *gin.Context) {
	userData := ctx.MustGet("userData").(jwt.MapClaims)
	followingID, err := strconv.Atoi(ctx.Param("id"))

	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error_status":  "Wrong Parameter",
			"error_message": fmt.Sprintf("%v not an integer", followingID),
		})
		return
	}

	followedID := uint(userData["id"].(float64))

	if uint(followingID) == followedID {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error_status":  "Wrong Parameter",
			"error_message": "FollowingID and FollowedID are same",
		})
		return
	}

	follow := entities.Follow{
		FollowingID: uint(followingID),
		FollowedID:  followedID,
	}

	err = sc.SosmedService.CreateFollow(follow)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error_status":  "Wrong Parameter",
			"error_message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"Status": "Success",
	})
}

func (sc *SosmedController) CreatePost(ctx *gin.Context) {
	var postDto PostDto
	userData := ctx.MustGet("userData").(jwt.MapClaims)
	userID := uint(userData["id"].(float64))

	if err := ctx.ShouldBindJSON(&postDto); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	post := entities.Post{
		Body:   postDto.Body,
		UserID: userID,
	}

	result, err := sc.SosmedService.CreatePost(post)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, result)
}

func (sc *SosmedController) GetAllUser(ctx *gin.Context) {
	userData := ctx.MustGet("userData").(jwt.MapClaims)
	userID := uint(userData["id"].(float64))

	var users []entities.User
	users, err := sc.SosmedService.GetAllUser(userID)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, users)
}

func (sc *SosmedController) GetFollowing(ctx *gin.Context) {

	userData := ctx.MustGet("userData").(jwt.MapClaims)
	userID := uint(userData["id"].(float64))

	result := []entities.GetFollowingResponse{}
	result, err := sc.SosmedService.GetFollowing(userID)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"following": len(result),
		"data":      result,
	})
}

func (sc *SosmedController) GetFollowed(ctx *gin.Context) {

	userData := ctx.MustGet("userData").(jwt.MapClaims)
	userID := uint(userData["id"].(float64))

	result := []entities.GetFollowedResponse{}
	result, err := sc.SosmedService.GetFollowed(userID)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"followed": len(result),
		"data":     result,
	})
}

func (sc *SosmedController) GetPostByFollowing(ctx *gin.Context) {
	userData := ctx.MustGet("userData").(jwt.MapClaims)
	userID := uint(userData["id"].(float64))

	result := []entities.Post{}
	result, err := sc.SosmedService.GetPostByFollowing(userID)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
	}

	ctx.JSON(http.StatusOK, result)
}

func (sc *SosmedController) GetPostByUserId(ctx *gin.Context) {
	userId, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error_status":  "Wrong Parameter",
			"error_message": fmt.Sprintf("%v not an integer", userId),
		})
		return
	}

	posts := []entities.Post{}
	posts, err = sc.SosmedService.GetPostByUserId(uint(userId))
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
	}

	ctx.JSON(http.StatusOK, posts)
}

func (sc *SosmedController) Unfollow(ctx *gin.Context) {
	userData := ctx.MustGet("userData").(jwt.MapClaims)
	userId := uint(userData["id"].(float64))
	unfollowId, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error_status":  "Wrong Parameter",
			"error_message": fmt.Sprintf("%v not an integer", unfollowId),
		})
		return
	}

	if uint(unfollowId) == userId {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error_status":  "Wrong Parameter",
			"error_message": "FollowingID and FollowedID are same",
		})
		return
	}

	err = sc.SosmedService.DeleteFollow(userId, uint(unfollowId))
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error_status":  "Wrong Parameter",
			"error_message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"Status": "Success",
	})
}

func (sc *SosmedController) DeletePost(ctx *gin.Context) {
	userData := ctx.MustGet("userData").(jwt.MapClaims)
	userId := uint(userData["id"].(float64))
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error_status":  "Wrong Parameter",
			"error_message": fmt.Sprintf("%v not an integer", id),
		})
		return
	}
	postId := uint(id)

	err = sc.SosmedService.DeletePost(userId, postId)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error_status":  "Wrong Parameter",
			"error_message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"Status": "Success",
	})
}
