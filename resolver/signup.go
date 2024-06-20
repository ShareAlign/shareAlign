package resolver

import (
	"context"

	"github.com/graph-gophers/graphql-go"
)

type signupArgs struct {
	Username string
	Email    string
	Password string
}

func (r *QueryResolver) Signup(ctx context.Context, args signupArgs) (graphql.ID, error) {
	return graphql.ID(0), nil
}
