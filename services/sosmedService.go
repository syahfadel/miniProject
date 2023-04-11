package services

import (
	"errors"
	"fmt"
	"miniProject/entities"

	"gorm.io/gorm"
)

type SosmedService struct {
	DB *gorm.DB
}

func (ss *SosmedService) CreateUser(user entities.User) (entities.User, error) {
	if err := ss.DB.Debug().Create(&user).Error; err != nil {
		return entities.User{}, err
	}
	return user, nil
}

func (ss *SosmedService) Login(user entities.User) (entities.User, error) {
	if err := ss.DB.Debug().Where("email = ?", user.Email).Take(&user).Error; err != nil {
		return entities.User{}, err
	}
	return user, nil
}

func (ss *SosmedService) CreateFollow(follow entities.Follow) error {
	err := ss.DB.Debug().Where("followed_id = ? AND following_id = ?", follow.FollowedID, follow.FollowingID).First(&follow).Error
	if err == nil {
		return errors.New(fmt.Sprintf("user id %v already follow id %v", follow.FollowedID, follow.FollowingID))
	}

	if err = ss.DB.Debug().Create(&follow).Error; err != nil {
		return err
	}

	return nil
}

func (ss *SosmedService) CreatePost(post entities.Post) (entities.Post, error) {
	if err := ss.DB.Debug().Create(&post).Error; err != nil {
		return entities.Post{}, err
	}

	return post, nil
}

func (ss *SosmedService) GetAllUser(userID uint) ([]entities.User, error) {
	var users []entities.User

	if err := ss.DB.Debug().Select("id", "user_name").Find(&users).Error; err != nil {
		return []entities.User{}, err
	}
	return users, nil
}

func (ss *SosmedService) GetFollowing(userID uint) ([]entities.GetFollowingResponse, error) {

	var follow []entities.Follow
	var result []entities.GetFollowingResponse

	if err := ss.DB.Debug().Model(&follow).Select("follows.following_id, users.user_name").Joins("left join users on users.id = follows.following_id").Where(`followed_id = ?`, userID).Scan(&result).Error; err != nil {
		return []entities.GetFollowingResponse{}, err
	}
	return result, nil
}

func (ss *SosmedService) GetFollowed(userID uint) ([]entities.GetFollowedResponse, error) {

	var follow []entities.Follow
	var result []entities.GetFollowedResponse

	if err := ss.DB.Debug().Model(&follow).Select("follows.followed_id, users.user_name").Joins("left join users on users.id = follows.followed_id").Where(`following_id = ?`, userID).Scan(&result).Error; err != nil {
		return []entities.GetFollowedResponse{}, err
	}
	return result, nil
}

func (ss *SosmedService) GetPostByFollowing(userID uint) ([]entities.Post, error) {
	var result []entities.Post
	err := ss.DB.Debug().Model(&entities.Post{}).Select("posts.id, posts.body, posts.user_id, posts.created_at").Joins("left join follows on follows.following_id = posts.user_id").Where("follows.followed_id = ?", userID).Scan(&result).Error

	if err != nil {
		return []entities.Post{}, err
	}
	return result, nil
}

func (ss *SosmedService) GetPostByUserId(userID uint) ([]entities.Post, error) {
	var result []entities.Post
	if err := ss.DB.Debug().Where("user_id = ?", userID).Find(&result).Error; err != nil {
		return []entities.Post{}, err
	}
	return result, nil
}

func (ss *SosmedService) DeleteFollow(userId, unfollowId uint) error {
	result := ss.DB.Debug().Where("following_id = ? AND followed_id = ?", unfollowId, userId).Delete(entities.Follow{})
	if result.RowsAffected == 0 {
		return errors.New(fmt.Sprintf("User with id %v not follow id %v", userId, unfollowId))
	}
	return nil
}

func (ss *SosmedService) DeletePost(userId, postId uint) error {
	result := ss.DB.Debug().Where("user_id = ? AND id = ?", userId, postId).Delete(&entities.Post{})
	if result.RowsAffected == 0 {
		return errors.New(fmt.Sprintf("Post with userId %v and postId %v not available", userId, postId))
	}
	return nil
}
