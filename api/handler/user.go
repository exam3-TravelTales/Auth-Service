package handler

import (
	"auth/api/auth"
	pb "auth/genproto/users"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// Register godoc
// @Summary Register user
// @Description create new users
// @Tags auth
// @Param info body users.RegisterRequest true "User info"
// @Success 200 {object} users.RegisterResponse
// @Failure 400 {object} string "Invalid data"
// @Failure 500 {object} string "Server error"
// @Router /api/v1/auth/register [post]
func (h Handler) Register(c *gin.Context) {
	h.Log.Info("Register is starting")
	req := pb.RegisterRequest{}
	if err := c.BindJSON(&req); err != nil {
		h.Log.Error(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
	res, err := h.User.Register(c, &req)
	if err != nil {
		h.Log.Error(err.Error())
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	h.Log.Info("Register ended")
	c.JSON(http.StatusOK, res)
}

// Login godoc
// @Summary login user
// @Description it generates new access and refresh tokens
// @Tags auth
// @Param userinfo body users.LoginRequest true "username and password"
// @Success 200 {object} users.Tokens
// @Failure 400 {object} string "Invalid date"
// @Failure 500 {object} string "error while reading from server"
// @Router /api/v1/auth/login [post]
func (h Handler) Login(c *gin.Context) {
	h.Log.Info("Login is working")
	req := pb.LoginRequest{}

	if err := c.BindJSON(&req); err != nil {
		h.Log.Error(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error1": err.Error()})
	}

	res, err := h.User.Login(c, &req)
	if err != nil {
		h.Log.Error(err.Error())
		c.JSON(500, gin.H{"error2": err.Error()})
		return
	}
	var token pb.Tokens
	err = auth.GeneratedAccessJWTToken(res, &token)

	if err != nil {
		h.Log.Error(err.Error())
		c.JSON(500, gin.H{"error3": err.Error()})
	}
	err = auth.GeneratedRefreshJWTToken(res, &token)
	if err != nil {
		h.Log.Error(err.Error())
		c.JSON(500, gin.H{"error4": err.Error()})
	}

	c.JSON(http.StatusOK, &token)
	h.Log.Info("login is succesfully ended")

}

// ResetPassword godoc
// @Security ApiKeyAuth
// @Summary ResetPass user
// @Description it changes your password to new one
// @Tags userAuth
// @Param userinfo body users.EmailRecoveryRequest true "passwords"
// @Success 200 {object} string
// @Failure 400 {object} string "Invalid date"
// @Failure 401 {object} string "Invalid token"
// @Failure 500 {object} string "error while reading from server"
// @Router /api/v1/auth/reset-password [post]
func (h Handler) ResetPassword(c *gin.Context) {
	h.Log.Info("ResetPassword is working")

	accessToken := c.GetHeader("Authorization")
	id, err := auth.GetUserIdFromAccessToken(accessToken)
	req := pb.EmailRecoveryRequest{UserId: id}
	if err != nil {
		h.Log.Error(err.Error())
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	if err := c.BindJSON(&req); err != nil {
		h.Log.Error(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error1": err.Error()})
		return
	}
	_, err = h.User.EmailRecovery(c, &req)
	if err != nil {
		h.Log.Error(err.Error())
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "password successfully reset"})

	h.Log.Info("ResetPassword ended")
}

// Refresh godoc
// @Summary Refresh token
// @Description it changes your access token
// @Tags auth
// @Param userinfo body users.CheckRefreshTokenRequest true "token"
// @Success 200 {object} users.Tokens
// @Failure 400 {object} string "Invalid date"
// @Failure 401 {object} string "Invalid token"
// @Failure 500 {object} string "error while reading from server"
// @Router /api/v1/auth/refresh [post]
func (h Handler) Refresh(c *gin.Context) {
	h.Log.Info("Refresh is working")
	req := pb.CheckRefreshTokenRequest{}
	if err := c.BindJSON(&req); err != nil {
		h.Log.Error(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
	_, err := auth.ValidateRefreshToken(req.RefreshToken)
	if err != nil {
		h.Log.Error(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
	id, err := auth.GetUserIdFromRefreshToken(req.RefreshToken)
	if err != nil {
		h.Log.Error(err.Error())
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
	}
	res := pb.Tokens{Refreshtoken: req.RefreshToken}

	err = auth.GeneratedAccessJWTToken(&pb.UserInfo{Id: id}, &res)
	if err != nil {
		h.Log.Error(err.Error())
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
	}
	c.JSON(http.StatusOK, &res)
}

// Logout godoc
// @Summary Logout user
// @Description you log out
// @Tags auth
// @Success 200 {object} string
// @Router /api/v1/auth/logout [post]
func (h Handler) Logout(c *gin.Context) {
	h.Log.Info("Logout is working")
	c.JSON(http.StatusOK, gin.H{"message": "logged out"})
	h.Log.Info("Logout ended")
}

// Profile godoc
// @Security ApiKeyAuth
// @Summary get user
// @Description you can see your profile
// @Tags users
// @Success 200 {object} users.GetProfileResponse
// @Failure 401 {object} string "Invalid token"
// @Failure 500 {object} string "error while reading from server"
// @Router /api/v1/users/profile [get]
func (h Handler) Profile(c *gin.Context) {
	h.Log.Info("Profile is working")
	accessToken := c.GetHeader("Authorization")
	id, err := auth.GetUserIdFromAccessToken(accessToken)
	if err != nil {
		h.Log.Error(err.Error())
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
	}
	res, err := h.User.GetProfile(c, &pb.UserId{Id: id})
	if err != nil {
		h.Log.Error(err.Error())
		c.JSON(500, gin.H{"error1": err.Error()})
	}
	c.JSON(http.StatusOK, res)
	h.Log.Info("Profile ended")
}

// UserProfileUpdate godoc
// @Security ApiKeyAuth
// @Summary ResetPass user
// @Description you can update your profile
// @Tags users
// @Param userinfo body users.UpdateProfileRequest true "info"
// @Success 200 {object} users.UpdateProfileResponse
// @Failure 400 {object} string "Invalid data"
// @Failure 500 {object} string "error while reading from server"
// @Router /api/v1/users/profile [put]
func (h Handler) UserProfileUpdate(c *gin.Context) {
	h.Log.Info("UserProfileUpdate is working")
	accessToken := c.GetHeader("Authorization")
	id, err := auth.GetUserIdFromAccessToken(accessToken)
	if err != nil {
		h.Log.Error(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
	req := pb.UpdateProfileRequest{Id: id}
	if err := c.BindJSON(&req); err != nil {
		h.Log.Error(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
	res, err := h.User.UpdateProfile(c, &req)
	if err != nil {
		h.Log.Error(err.Error())
		c.JSON(500, gin.H{"error": err.Error()})
	}
	c.JSON(http.StatusOK, res)
	h.Log.Info("UserProfileUpdate ended")
}

// GetAllUsers godoc
// @Security ApiKeyAuth
// @Summary all users
// @Description you can see all users
// @Tags users
// @Param limit query string false "Number of users to fetch"
// @Param offset query string false "Number of users to omit"
// @Success 200 {object} users.GetUsersResponse
// @Failure 400 {object} string "Invalid data"
// @Failure 500 {object} string "error while reading from server"
// @Router /api/v1/users [get]
func (h Handler) GetAllUsers(c *gin.Context) {
	h.Log.Info("GetAllUsers is working")
	req := pb.GetUsersRequest{}
	limitStr := c.Query("limit")
	offsetStr := c.Query("offset")

	if limitStr != "" {
		limit, err := strconv.Atoi(limitStr)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest,
				gin.H{"error": err.Error()})
			h.Log.Error(err.Error())
			return
		}
		req.Limit = int64(limit)
	} else {
		req.Limit = 0
	}

	if offsetStr != "" {
		offset, err := strconv.Atoi(offsetStr)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest,
				gin.H{"error": err.Error()})
			h.Log.Error(err.Error())
			return
		}
		req.Offset = int64(offset)
	} else {
		req.Offset = 0
	}

	res, err := h.User.GetUsers(c, &req)
	if err != nil {
		h.Log.Error(err.Error())
		c.JSON(500, gin.H{"error": err.Error()})
	}
	c.JSON(http.StatusOK, res)
	h.Log.Info("GetAllUsers ended")
}

// Delete godoc
// @Security ApiKeyAuth
// @Summary delete user
// @Description you can delete your profile
// @Tags users
// @Param user_id path string true "user_id"
// @Success 200 {object} string
// @Failure 400 {object} string "Invalid data"
// @Failure 500 {object} string "error while reading from server"
// @Router /api/v1/users/{user_id} [delete]
func (h Handler) Delete(c *gin.Context) {
	h.Log.Info("Delete is working")
	id := c.Param("user_id")
	_, err := uuid.Parse(id)
	if err != nil {
		h.Log.Error(err.Error())
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user id is incorrect"})
		return
	}

	_, err = h.User.DeleteUser(c, &pb.UserId{Id: id})
	if err != nil {
		h.Log.Error(err.Error())
		c.JSON(500, gin.H{"error": err.Error()})
	}
	c.JSON(http.StatusOK, gin.H{"message": "user deleted"})
	h.Log.Info("Delete ended")
}

// ActivityOfUser godoc
// @Security ApiKeyAuth
// @Summary Activities user
// @Description you can see your profile activity
// @Tags users
// @Param user_id path string true "user_id"
// @Success 200 {object} users.ActivityResponse
// @Failure 400 {object} string "Invalid data"
// @Failure 500 {object} string "error while reading from server"
// @Router /api/v1/users/{user_id}/activity [get]
func (h Handler) ActivityOfUser(c *gin.Context) {
	h.Log.Info("Activity is working")
	id := c.Param("user_id")
	_, err := uuid.Parse(id)
	if err != nil {
		h.Log.Error(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": "user id is incorrect"})
		return
	}
	req := pb.UserId{
		Id: id,
	}

	res, err := h.User.Activity(c, &req)

	if err != nil {
		h.Log.Error(err.Error())
		c.JSON(500, gin.H{"error": err.Error()})
	}

	c.JSON(http.StatusOK, &res)
	h.Log.Info("Activity ended")
}

// Follow godoc
// @Security ApiKeyAuth
// @Summary follow user
// @Description you can follow another user
// @Tags users
// @Param user_id path string true "user_id"
// @Success 200 {object} users.FollowResponse
// @Failure 400 {object} string "Invalid data"
// @Failure 500 {object} string "error while reading from server"
// @Router /api/v1/users/{user_id}/follow [post]
func (h Handler) Follow(c *gin.Context) {
	h.Log.Info("Follow is working")
	id := c.Param("user_id")
	_, err := uuid.Parse(id)
	if err != nil {
		h.Log.Error(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": "user id is incorrect"})
	}

	accessToken := c.GetHeader("Authorization")
	idFollower, err := auth.GetUserIdFromAccessToken(accessToken)
	if err != nil {
		h.Log.Error(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": "unauthorized"})
	}

	res, err := h.User.Follow(c, &pb.FollowRequest{FollowerId: idFollower, FollowingId: id})
	if err != nil {
		h.Log.Error(err.Error())
		c.JSON(500, gin.H{"error": err.Error()})
	}
	c.JSON(http.StatusOK, res)
	h.Log.Info("Follow ended")
}

// GetFollowers godoc
// @Security ApiKeyAuth
// @Summary get followers
// @Description you can see your followers
// @Tags users
// @Param user_id path string true "user_id"
// @Param limit query string false "Number of users to fetch"
// @Param offset query string false "Number of users to omit"
// @Success 200 {object} users.FollowersResponse
// @Failure 400 {object} string "Invalid data"
// @Failure 500 {object} string "error while reading from server"
// @Router /api/v1/users/{user_id}/followers [get]
func (h Handler) GetFollowers(c *gin.Context) {
	h.Log.Info("Followers is working")
	id := c.Param("user_id")
	_, err := uuid.Parse(id)
	if err != nil {
		h.Log.Error(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": "user id is incorrect"})
	}
	req := pb.FollowersRequest{UserId: id}

	limitStr := c.Query("limit")
	offsetStr := c.Query("offset")

	if limitStr != "" {
		limit, err := strconv.Atoi(limitStr)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest,
				gin.H{"error": err.Error()})
			h.Log.Error(err.Error())
			return
		}
		req.Limit = int64(limit)
	} else {
		req.Limit = 0
	}

	if offsetStr != "" {
		offset, err := strconv.Atoi(offsetStr)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest,
				gin.H{"error": err.Error()})
			h.Log.Error(err.Error())
			return
		}
		req.Offset = int64(offset)
	} else {
		req.Offset = 0
	}

	res, err := h.User.Followers(c, &req)
	if err != nil {
		h.Log.Error(err.Error())
	}
	c.JSON(http.StatusOK, res)
	h.Log.Info("Followers ended")
}
