package services

import (
	"encoding/json"
	"myapp/db"
	"myapp/models"
	"myapp/repositories"
	"myapp/utils/logic"
	"net/http"
)

/*
 Todoリストを取得しレスポンス用に変換
*/
func GetAllTodos(w http.ResponseWriter, userId int) ([]models.BaseTodoResponse, error) {
	var todos []models.Todo
	// todoリストデータ取得
	if err := repositories.GetAllTodos(&todos, userId); err != nil {
		logic.SendResponse(w, logic.CreateErrorStringResponse("データ取得に失敗"), http.StatusInternalServerError)
		return nil, err
	}

	// レスポンス用の構造体に変換
	var responseTodos []models.BaseTodoResponse
	for _, todo := range todos {
		var newTodo models.BaseTodoResponse
		newTodo.BaseModel.ID = todo.BaseModel.ID
		newTodo.BaseModel.CreatedAt = todo.BaseModel.CreatedAt
		newTodo.BaseModel.UpdatedAt = todo.BaseModel.UpdatedAt
		newTodo.BaseModel.DeletedAt = todo.BaseModel.DeletedAt
		newTodo.Title = todo.Title
		newTodo.Comment = todo.Comment
		responseTodos = append(responseTodos, newTodo)
	}

	return responseTodos, nil
}

func GetTodoById(todo *models.Todo, id string, userId int) {
	db := db.GetDB()
	db.Joins("User").Where("user_id=?", userId).First(&todo, id)
}

func InsertTodo(todo *models.Todo) {
	db := db.GetDB()
	db.Create(&todo)
}

func DeleteTodo(id string, userId int) {
	db := db.GetDB()
	db.Where("id=? AND user_id=?", id, userId).Delete(&models.Todo{})
}

func UpdateTodo(todo *models.Todo, id string) {
	db := db.GetDB()
	db.Model(&todo).Where("id=? AND user_id=?", id, todo.UserId).Updates(
        map[string]interface{}{
            "title":     todo.Title,
            "comment":    todo.Comment,
			"user_id": todo.UserId,
        })
}


/*
 Todoリスト取得APIのレスポンス送信処理
*/
func SendAllTodoResponse(w http.ResponseWriter, todos *[]models.BaseTodoResponse) {
	var response models.AllTodoResponse
	response.Todos = *todos
	// レスポンスデータ作成
	responseBody, _ := json.Marshal(response)

	// レスポンス送信
	logic.SendResponse(w, responseBody, http.StatusCreated)
}