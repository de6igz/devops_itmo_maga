package main

import (
	"reflect"

	"game-catalog-backend/internal/game"

	"github.com/golang/mock/gomock"
)

type MockGameRepository struct {
	ctrl     *gomock.Controller
	recorder *MockGameRepositoryMockRecorder
}

type MockGameRepositoryMockRecorder struct {
	mock *MockGameRepository
}

func NewMockGameRepository(ctrl *gomock.Controller) *MockGameRepository {
	mock := &MockGameRepository{ctrl: ctrl}
	mock.recorder = &MockGameRepositoryMockRecorder{mock}
	return mock
}

func (m *MockGameRepository) EXPECT() *MockGameRepositoryMockRecorder {
	return m.recorder
}

func (m *MockGameRepository) Create(entity game.Game) (game.Game, error) {
	m.ctrl.T.Helper()
	results := m.ctrl.Call(m, "Create", entity)
	created, _ := results[0].(game.Game)
	err, _ := results[1].(error)
	return created, err
}

func (mr *MockGameRepositoryMockRecorder) Create(entity any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockGameRepository)(nil).Create), entity)
}

func (m *MockGameRepository) List(filters game.Filters) ([]game.Game, error) {
	m.ctrl.T.Helper()
	results := m.ctrl.Call(m, "List", filters)
	games, _ := results[0].([]game.Game)
	err, _ := results[1].(error)
	return games, err
}

func (mr *MockGameRepositoryMockRecorder) List(filters any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "List", reflect.TypeOf((*MockGameRepository)(nil).List), filters)
}

func (m *MockGameRepository) GetByID(id int64) (game.Game, error) {
	m.ctrl.T.Helper()
	results := m.ctrl.Call(m, "GetByID", id)
	entity, _ := results[0].(game.Game)
	err, _ := results[1].(error)
	return entity, err
}

func (mr *MockGameRepositoryMockRecorder) GetByID(id any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByID", reflect.TypeOf((*MockGameRepository)(nil).GetByID), id)
}

func (m *MockGameRepository) Update(id int64, entity game.Game) (game.Game, error) {
	m.ctrl.T.Helper()
	results := m.ctrl.Call(m, "Update", id, entity)
	updated, _ := results[0].(game.Game)
	err, _ := results[1].(error)
	return updated, err
}

func (mr *MockGameRepositoryMockRecorder) Update(id, entity any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockGameRepository)(nil).Update), id, entity)
}

func (m *MockGameRepository) Delete(id int64) error {
	m.ctrl.T.Helper()
	results := m.ctrl.Call(m, "Delete", id)
	err, _ := results[0].(error)
	return err
}

func (mr *MockGameRepositoryMockRecorder) Delete(id any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockGameRepository)(nil).Delete), id)
}
