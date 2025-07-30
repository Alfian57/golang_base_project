package factory

import (
	"github.com/Alfian57/belajar-golang/internal/constants"
	"github.com/Alfian57/belajar-golang/internal/model"
	"github.com/bluele/factory-go/factory"
	"github.com/brianvoe/gofakeit/v6"
	"github.com/google/uuid"
)

func NewMemberFactory() *factory.Factory {
	return factory.NewFactory(
		&model.User{},
	).SeqInt("ID", func(n int) (any, error) {
		return uuid.New(), nil
	}).Attr("Email", func(args factory.Args) (any, error) {
		return gofakeit.Email(), nil
	}).Attr("Username", func(args factory.Args) (any, error) {
		return gofakeit.Username(), nil
	}).Attr("Role", func(args factory.Args) (any, error) {
		roles := []string{model.UserRoleAdmin, model.UserRoleMember}
		return gofakeit.RandomString(roles), nil
	}).OnCreate(func(args factory.Args) error {
		user := args.Instance().(*model.User)
		// Set a default password for the user
		return user.SetHashedPassword(constants.DefaultPassword)
	})
}

func NewAdminFactory() *factory.Factory {
	return factory.NewFactory(
		&model.User{},
	).SeqInt("ID", func(n int) (any, error) {
		return uuid.New(), nil
	}).Attr("Email", func(args factory.Args) (any, error) {
		return gofakeit.Email(), nil
	}).Attr("Username", func(args factory.Args) (any, error) {
		return "admin_" + gofakeit.Username(), nil
	}).Attr("Role", func(args factory.Args) (any, error) {
		return model.UserRoleAdmin, nil
	}).OnCreate(func(args factory.Args) error {
		user := args.Instance().(*model.User)
		// Set a default password for the admin user
		return user.SetHashedPassword(constants.DefaultPassword)
	})
}
