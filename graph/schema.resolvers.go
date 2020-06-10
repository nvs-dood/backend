package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"github.com/EnglederLucas/nvs-dood/auth"
	"github.com/EnglederLucas/nvs-dood/graph/generated"
	"github.com/EnglederLucas/nvs-dood/graph/models"
	"github.com/EnglederLucas/nvs-dood/repository"
)

func (r *mutationResolver) AddStudent(ctx context.Context, input models.NewStudent) (*models.Student, error) {
	student := &models.Student{
		Name:  input.Name,
		Role:  input.Role,
		Class: input.Class,
	}

	err := r.DB.Create(&student).Error
	if err != nil {
		return nil, err
	}
	return student, nil
}

func (r *mutationResolver) AddShift(ctx context.Context, studentID string, newShift models.InputShift) (*models.Student, error) {
	var shifts []*models.InputShift

	return r.AddShifts(ctx, studentID, append(shifts, &newShift))
}

func (r *mutationResolver) AddShifts(ctx context.Context, studentID string, newShifts []*models.InputShift) (*models.Student, error) {
	var student, err = repository.GetStudentByID(r.DB, studentID)

	if err != nil {
		return nil, err
	}

	for _, element := range newShifts {
		shift := &models.Shift{
			Start: element.Start,
			End:   element.End,
		}

		student.Shifts = append(student.Shifts, shift)
	}

	err = r.DB.Save(student).Error

	//err = r.DB.Create(&shift).Error
	if err != nil {
		return nil, err
	}

	return student, nil
}

func (r *mutationResolver) UpdateShifts(ctx context.Context, studentID string, newShifts []*models.InputShift) (*models.Student, error) {
	var student, err = repository.GetStudentByID(r.DB, studentID)

	if err != nil {
		return nil, err
	}

	err = r.DB.Where("student_id = ?", studentID).Delete(models.Shift{}).Error

	if err != nil {
		return nil, err
	}

	for _, element := range newShifts {
		curShift := &models.Shift{
			StudentID: studentID,
			Start:     element.Start,
			End:       element.End,
		}

		err = r.DB.Create(curShift).Error
		if err != nil {
			return nil, err
		}
	}

	return student, nil
}

func (r *queryResolver) Students(ctx context.Context) ([]*models.Student, error) {
	if user := auth.ForContext(ctx); user == nil || !user.Admin {
		return nil, fmt.Errorf("Access denied")
	}

	var students []*models.Student

	err := r.DB.Find(&students).Error
	if err != nil {
		return nil, err
	}
	return students, nil
}

func (r *queryResolver) AllWithRole(ctx context.Context, role models.Role) ([]*models.Student, error) {
	var students []*models.Student

	err := r.DB.Where("role = ?", role.String()).Find(&students).Error
	if err != nil {
		return nil, err
	}
	return students, nil
}

func (r *queryResolver) AllInClass(ctx context.Context, class string) ([]*models.Student, error) {
	var students []*models.Student

	err := r.DB.Where("class LIKE ?", class).Find(&students).Error
	if err != nil {
		return nil, err
	}
	return students, nil
}

func (r *queryResolver) GetStudentByID(ctx context.Context, studentID string) (*models.Student, error) {
	return repository.GetStudentByID(r.DB, studentID)
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
