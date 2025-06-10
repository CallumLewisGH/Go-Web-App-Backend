package userModel

type UpdateUserRequest struct {
	Username       *string
	Email          *string
	Timezone       *string
	ProfilePicture *string
	Bio            *string
}
