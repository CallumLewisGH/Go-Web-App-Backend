package query

import (
	"context"
	"fmt"

	cqrs "github.com/CallumLewisGH/Generic-Service-Base/internal/domain"
	repos "github.com/CallumLewisGH/Generic-Service-Base/internal/domain/repositories"
	userModel "github.com/CallumLewisGH/Generic-Service-Base/internal/domain/user"
	"gorm.io/gorm"
)

// This serves as an example for how to do concurrent queries in this architecture and is not tested
// This should not be used as is
// By all means use it as a template just don't make me cry by using it directly
// I don't think this will ever get used as it sacrifices too much type safety but it is here just in case
func GetTwoThingsQuery() (*userModel.UserDTO, *userModel.UserDTO, error) {
	queryFuncUserById := func(db *gorm.DB, ctx context.Context) (any, error) {
		user := userModel.User{}
		if err := repos.NewUserRepo(db).WithUsername("callum").First(&user); err != nil {
			return nil, err
		}
		return user.ToUserDTO(), nil
	}

	queryFuncUserByAuthId := func(db *gorm.DB, ctx context.Context) (any, error) {
		user := userModel.User{}
		user.AuthId = "example-auth-id"
		if err := repos.NewUserRepo(db).First(&user); err != nil {
			return nil, err
		}
		return user, nil
	}

	results := cqrs.ConcurrentQueries(
		queryFuncUserById,
		queryFuncUserByAuthId,
	)

	// Check errors first
	if results[0].Err != nil || results[1].Err != nil {
		return nil, nil, fmt.Errorf("queries failed: %v, %v", results[0].Err, results[1].Err)
	}

	// Type assert both results
	userById, ok1 := results[0].Data.(*userModel.UserDTO)
	userByAuthId, ok2 := results[1].Data.(*userModel.UserDTO)

	if !ok1 || !ok2 {
		return nil, nil, fmt.Errorf("type assertion failed")
	}

	return userById, userByAuthId, nil
}
