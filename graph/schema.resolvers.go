package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"
	"github.com/EnglederLucas/nvs-dood/auth"
	generated1 "github.com/EnglederLucas/nvs-dood/graph/generated"
	"github.com/EnglederLucas/nvs-dood/graph/models"
	"github.com/EnglederLucas/nvs-dood/repository"
	"time"
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

func (r *mutationResolver) AddMeAsStudent(ctx context.Context, input models.InputStudent) (*models.Student, error) {
	user := auth.ForContext(ctx)
	if user == nil || !user.Admin {
		return nil, fmt.Errorf("Access denied")
	}

	student := &models.Student{
		Class: input.Class,
		Role:  input.Role,
		ID:    user.ID,
		Name:  *user.Name,
	}

	err := r.DB.Create(student).Error
	if err != nil {
		fmt.Printf("Could not create student from user %s", user.ID)
		return nil, fmt.Errorf("Could not create user")
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

func (r *mutationResolver) GroupEntersRoom(ctx context.Context, enterRoom models.EnterRoomInput) (*models.RoomStay, error) {
	var _, err = repository.GetStudentByID(r.DB, enterRoom.StudentID)

	if err != nil {
		return nil, err
	}

	roomStay := &models.RoomStay{
		Room:      enterRoom.Room,
		StudentID: enterRoom.StudentID,
		GroupSize: enterRoom.GroupSize,
		Start:     &enterRoom.Start,
	}

	//When passed with pointer reference, the id will be added to the roomStay struct
	err = r.DB.Create(&roomStay).Error

	if err != nil {
		return nil, err
	}

	return roomStay, nil
}

func (r *mutationResolver) GroupLeavesRoom(ctx context.Context, leaveRoom models.LeaveRoomInput) (*models.RoomStay, error) {
	var room models.RoomStay

	err := r.DB.Where(&models.RoomStay{StayID: leaveRoom.RoomStayID}).Find(&room).Error
	if err != nil {
		return nil, err
	}

	if room.Start.After(leaveRoom.End) {
		return nil, fmt.Errorf("Endtime has to be after Starttime. This API is not for time travellers")
	}

	if room.End != nil {
		return nil, fmt.Errorf("Endtime was already set. You can't leave a room twice in a row")
	}

	room.End = &leaveRoom.End
	err = r.DB.Save(&room).Error
	if err != nil {
		return nil, err
	}

	return &room, nil
}

func (r *mutationResolver) AddRoomStay(ctx context.Context, input models.RoomStayInput) (*models.RoomStay, error) {
	roomStay := &models.RoomStay{
		Room:      input.Room,
		StudentID: input.StudentID,
		GroupSize: input.GroupSize,
		Start:     &input.Start,
		End:       &input.End,
	}

	if roomStay.Start.After(*roomStay.End) {
		return nil, fmt.Errorf("Endtime has to be after Starttime. This API is not for time travellers")
	}

	//When passed with pointer reference, the id will be added to the roomStay struct
	err := r.DB.Create(&roomStay).Error

	if err != nil {
		return nil, err
	}

	return roomStay, nil
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

func (r *queryResolver) AllVisitsInRoom(ctx context.Context, room string) ([]*models.RoomStay, error) {
	if user := auth.ForContext(ctx); user == nil || !user.Admin {
		return nil, fmt.Errorf("Access denied")
	}

	var stays []*models.RoomStay

	err := r.DB.Where("room LIKE ?", room).Find(&stays).Error
	if err != nil {
		return nil, err
	}
	return stays, nil
}

func (r *queryResolver) AllRoomActivities(ctx context.Context) ([]*models.RoomStay, error) {
	if user := auth.ForContext(ctx); user == nil || !user.Admin {
		return nil, fmt.Errorf("Access denied")
	}

	var activities []*models.RoomStay

	err := r.DB.Find(&activities).Error
	if err != nil {
		return nil, err
	}
	return activities, nil
}

func (r *queryResolver) CurrentNumberOfPeople(ctx context.Context, room string) (int, error) {
	var people = 0

	err := r.DB.Model(&models.RoomStay{}).Where("room LIKE ? AND start < ?  AND end is NULL", room, time.Now()).Count(&people).Error
	if err != nil {
		return -1, err
	}
	return people, nil
}

// Mutation returns generated1.MutationResolver implementation.
func (r *Resolver) Mutation() generated1.MutationResolver { return &mutationResolver{r} }

// Query returns generated1.QueryResolver implementation.
func (r *Resolver) Query() generated1.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
