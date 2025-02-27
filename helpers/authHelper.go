package helper

import (
	"errors"

	"github.com/gin-gonic/gin"
)

func CheckUserType(c *gin.Context, role string) (err error) {
	userType := c.GetString("user_type")
	err = nil

	if userType != role {
		err = errors.New("Unauthorized access to this resource.")
		return err
	}

	return err
}

func MatchUserTypeToUid(c *gin.Context, userId string) (err error) {
	userType := c.GetString("user_type")
	uid := c.GetString("uid")
	err = nil

	// Do not allow a user to look at other users' data
	if userType == "USER" && uid != userId {
		err = errors.New("Unauthorized access to this resource.")
		return err
	}

	err = CheckUserType(c, userType)

	return err
}
