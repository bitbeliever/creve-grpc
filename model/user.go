package model

import (
	"errors"
	"github.com/xblymmx/creve-grpc/pkg/id"
	userservicepb "github.com/xblymmx/creve-grpc/proto/user/pb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
	"time"
)

type User struct {
	ID        int64 `gorm:"primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`

	UserID   int64
	Username string `gorm:"type:varchar(30);unique_index"`
	Password string
	Avatar   string
	Email    string
	Salt     string
	Phone    string
	Status   int
}

func (u User) TableName() string {
	return "t_users"
}

func init() {
	if err := db.AutoMigrate(&User{}); err != nil {
		panic(err)
	}
}

var (
	ErrUsernameExists = errors.New("error username exists")
)

func newDefaultUser() User {
	return User{
		Avatar: "https://gw.alipayobjects.com/zos/antfincdn/XAosXuNZyF/BiazfanxmamNRoxxVxka.png",
	}
}

func NewUser(in *userservicepb.NewUserRequest) (User, error) {
	if len(in.Password) < 4 {
		return User{}, status.Error(codes.InvalidArgument, "password length must be greater than 4")
	}
	if len(in.Username) < 4 {
		return User{}, status.Error(codes.InvalidArgument, "username length must be greater than 4")
	}

	var user User
	d := db.Where("username = ?", in.Username).Find(&user)
	if d.Error != nil {
		return user, d.Error
	}
	if d.RowsAffected > 0 {
		return user, status.New(codes.AlreadyExists, ErrUsernameExists.Error()).Err()
	}
	if d.RowsAffected == 0 {
		user = newDefaultUser()
		user.Username = in.Username
		user.Password = in.Password
		user.UserID = id.GenerateInt64()
		if err := db.Create(&user).Error; err != nil {
			return user, status.New(codes.Internal, err.Error()).Err()
		}
	} else {
		panic("unreachable")
	}

	return user, nil
}

func GetUserByName(username string) (User, error) {
	user := User{}
	err := db.Where("username = ?", username).First(&user).Error
	return user, err
}

func GetUserByID(id int64) (User, error) {
	user := User{}
	err := db.Where("id = ?", id).First(&user).Error
	return user, err
}

func UpdateUserInfo(in *userservicepb.UpdateUserInfoRequest) (User, error) {
	user := User{}
	err := db.Model(&user).Where("id = ?", in.Id).Updates(map[string]interface{}{
		"username": in.Username,
		"avatar":   in.Avatar,
		"email":    in.Email,
	}).Error
	return user, err
}

//func UserLogin(in *userservicepb.LoginRequest) (User, error) {
//	user := User{}
//	err := db.Where("username = ? and password = ?", in.Username, in.Password).First(&user).Error
//	return user, err
//}
