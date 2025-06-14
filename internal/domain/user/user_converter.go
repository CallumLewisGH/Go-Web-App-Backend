package userModel

func (user *User) ToUserDTO() *UserDTO {
	return &UserDTO{
		ID:             user.ID,
		AuthId:         user.AuthId,
		CreatedAt:      user.CreatedAt,
		UpdatedAt:      user.UpdatedAt,
		Username:       user.Username,
		Email:          user.Email,
		EmailVerified:  user.EmailVerified,
		LastLogin:      user.LastLogin,
		ProfilePicture: user.ProfilePicture,
		Bio:            user.Bio,
		Timezone:       user.Timezone,
		IsActive:       user.IsActive,
		IsBanned:       user.IsBanned,
		DeactivatedAt:  user.DeactivatedAt,
	}
}

func ToUserDTOs(users []User) []UserDTO {
	userDTOs := []UserDTO{}
	if len(users) > 0 {
		println(len(users))
		for _, i := range users {
			userDTOs = append(userDTOs, *i.ToUserDTO())
		}
	}
	return userDTOs
}
