package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"
	"math/rand"

	"github.com/EnglederLucas/nvs-dood/graph/generated"
	"github.com/EnglederLucas/nvs-dood/graph/models"
)

func (r *mutationResolver) AddStudent(ctx context.Context, input models.NewStudent) (*models.Student, error) {
	student := &models.Student{
		ID:    fmt.Sprint("T%d", rand.Int()),
		Name:  input.Name,
		Role:  input.Role,
		Class: input.Class,
	}

	r.students = append(r.students, student)
	return student, nil
}

func (r *mutationResolver) AddShift(ctx context.Context, newShift models.InputShift) (*models.Student, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) Students(ctx context.Context) ([]*models.Student, error) {
	return r.students, nil
}

func (r *queryResolver) AllWithRole(ctx context.Context) ([]*models.Student, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) AllInClass(ctx context.Context) ([]*models.Student, error) {
	panic(fmt.Errorf("not implemented"))
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
