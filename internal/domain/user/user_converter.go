package userModel

func (user *User) ToUserDTO() *UserDTO {
	return &UserDTO{
		ID:             user.ID,
		Username:       user.Username,
		Email:          user.Email,
		EmailVerified:  user.EmailVerified,
		LastLogin:      user.LastLogin,
		ProfilePicture: user.ProfilePicture,
		Bio:            user.Bio,
		Locale:         user.Locale,
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
		for i := range users {
			userDTOs = append(userDTOs, *users[i].ToUserDTO())
		}
	}
	return userDTOs
}
