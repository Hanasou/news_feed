package core

import (
	"testing"

	"github.com/Hanasou/news_feed/go/todo/models"
	"github.com/stretchr/testify/require"
)

func TestInitializeService(t *testing.T) {
	tests := []struct {
		name     string
		dbType   string
		rootPath string
		wantErr  bool
	}{
		{
			name:     "Valid service initialization",
			dbType:   "mem",
			rootPath: "",
			wantErr:  false,
		},
		{
			name:     "Invalid db type",
			dbType:   "mysql",
			rootPath: "",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		newService, err := InitializeService(tt.dbType, tt.rootPath, false)
		if tt.wantErr {
			require.Error(t, err, "Expected error for db type: %s", tt.dbType)
		} else {
			require.NoError(t, err, "Did not expect error for db type: %s", tt.dbType)
			require.NotNil(t, newService, "Expected service to be initialized")
			require.NotNil(t, newService.todoTable, "Expected todo table to be initialized")
			require.NotNil(t, newService.userTable, "Expected user table to be initialized")
		}
	}
}

func TestTodoService_CreateTodo(t *testing.T) {
	// Setup test service
	tests := []struct {
		name    string
		todo    *models.Todo
		wantErr bool
	}{
		{
			name: "Valid todo creation",
			todo: &models.Todo{
				Id:     "todo1",
				Text:   "Test todo item",
				Done:   false,
				UserId: "user1",
			},
			wantErr: false,
		},
		{
			name: "Another valid todo",
			todo: &models.Todo{
				Id:     "todo2",
				Text:   "Another test todo",
				Done:   true,
				UserId: "user2",
			},
			wantErr: false,
		},
		{
			name: "Todo with empty text",
			todo: &models.Todo{
				Id:     "todo3",
				Text:   "",
				Done:   false,
				UserId: "user1",
			},
			wantErr: false, // Assuming empty text is allowed
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service, err := InitializeService("mem", "", false)
			require.NoError(t, err)

			err = service.CreateTodo(tt.todo)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateTodo() error = %v, wantErr %v", err, tt.wantErr)
			} else if err == nil {
				t.Logf("Created todo: %v", tt.todo)
			}
		})
	}
}

func TestTodoService_GetTodos(t *testing.T) {
	service, err := InitializeService("mem", "", false)
	require.NoError(t, err)

	// Create some todos
	todos := []*models.Todo{
		{Id: "todo1", Text: "First todo", Done: false, UserId: "user1"},
		{Id: "todo2", Text: "Second todo", Done: true, UserId: "user2"},
	}

	// Mock user data
	users := []*models.User{
		{Id: "user1", Name: "User One"},
		{Id: "user2", Name: "User Two"},
	}
	for _, user := range users {
		err = service.userTable.Upsert(user)
		require.NoError(t, err)
	}

	for _, todo := range todos {
		err = service.CreateTodo(todo)
		require.NoError(t, err)
	}

	retrievedTodos, err := service.GetTodos("user1")
	for _, todo := range retrievedTodos {
		t.Logf("Retrieved todo: %v", todo)
	}
	require.NoError(t, err)
	require.Len(t, retrievedTodos, 1)

	for _, todo := range retrievedTodos {
		require.Equal(t, "todo1", todo.Id)
		require.Equal(t, "First todo", todo.Text)
		require.Equal(t, false, todo.Done)
	}
}
