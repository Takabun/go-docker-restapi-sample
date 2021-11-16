package controllers

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"myapp/models"
	"myapp/services"
	"myapp/utils/logic"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type DeleteTodoResponse struct {
    Id string `json:"id"`
}

/*
 Todoリスト取得
*/
func fetchAllTodos(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-type", "application/json")
	// トークンからuserIdを取得
	userId, err := logic.GetUserIdFromContext(r)
	if err != nil {
		// レスポンスデータ作成
		response := map[string]interface{}{
			"err": "認証エラー",
		}
		responseBody, err := json.Marshal(response)
		if err != nil {
			log.Fatal(err)
		}
		w.WriteHeader(http.StatusUnauthorized)
		w.Write(responseBody)
	}

	var todos []models.Todo
    services.GetAllTodos(&todos, userId)

	// レスポンスデータ作成
	response := map[string]interface{}{
		"todos": todos,
	}
    responseBody, err := json.Marshal(response)
    if err != nil {
        log.Fatal(err)
    }
	w.WriteHeader(http.StatusOK) // ステータスコード
    w.Write(responseBody)
}


/*
 idに紐づくTodoを取得
*/
func fetchTodoById(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    // トークンからuserIdを取得
	userId, err := logic.GetUserIdFromContext(r)
	if err != nil {
		// レスポンスデータ作成
		response := map[string]interface{}{
			"err": "認証エラー",
		}
		responseBody, err := json.Marshal(response)
		if err != nil {
			log.Fatal(err)
		}
		w.WriteHeader(http.StatusUnauthorized)
		w.Write(responseBody)
	}

    vars := mux.Vars(r)
    id := vars["id"] 

    var todo models.Todo
    services.GetTodoById(&todo, id, userId)

    if todo.ID == 0 {
        // レスポンスデータ作成
		response := map[string]interface{}{
			"err": "データがありません。",
		}
		responseBody, _ := json.Marshal(response)
        w.WriteHeader(http.StatusBadRequest)
		w.Write(responseBody)
        return
    }

    responseBody, err := json.Marshal(todo)
    if err != nil {
        log.Fatal(err)
    }

    w.WriteHeader(http.StatusOK) // ステータスコード
    w.Write(responseBody)
}

/*
 Todo新規登録
*/
func createTodo(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-type", "application/json")
    // トークンからuserIdを取得
	userId, err := logic.GetUserIdFromContext(r)
	if err != nil {
		// レスポンスデータ作成
		response := map[string]interface{}{
			"err": "認証エラー",
		}
		responseBody, err := json.Marshal(response)
		if err != nil {
			log.Fatal(err)
		}
		w.WriteHeader(http.StatusUnauthorized)
		w.Write(responseBody)
	}
    // ioutil: ioに特化したパッケージ
    reqBody,_ := ioutil.ReadAll(r.Body)
    var todo models.Todo
    // json.Unmarshal()
    // 第１引数で与えたjsonデータを、第二引数に指定した値にマッピングする
    // 返り値はerrorで、エラーが発生しない場合はnilになる
    if err := json.Unmarshal(reqBody, &todo); err != nil {
        log.Fatal(err)
    }

    todo.UserId = userId

    services.InsertTodo(&todo)

    responseBody, err := json.Marshal(todo)
    if err != nil {
        log.Fatal(err)
    }

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)
    w.Write(responseBody)
}

/*
 削除処理
*/
func deleteTodo(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
     // トークンからuserIdを取得
	userId, err := logic.GetUserIdFromContext(r)
	if err != nil {
		// レスポンスデータ作成
		response := map[string]interface{}{
			"err": "認証エラー",
		}
		responseBody, err := json.Marshal(response)
		if err != nil {
			log.Fatal(err)
		}
		w.WriteHeader(http.StatusUnauthorized)
		w.Write(responseBody)
	}
    vars := mux.Vars(r)
    id := vars["id"]

    var todo models.Todo

    // 削除データの確認
    services.GetTodoById(&todo, id, userId)
    if todo.ID == 0 {
        // レスポンスデータ作成
		response := map[string]interface{}{
			"err": "データがありません。",
		}
		responseBody, _ := json.Marshal(response)
        w.WriteHeader(http.StatusBadRequest)
		w.Write(responseBody)
        return
    }

    services.DeleteTodo(id, userId)
    // responseBody, err := json.Marshal(DeleteResponse{Id: id})
    // responseBody, err := json.Marshal(todo)
    // if err != nil {
    //     log.Fatal(err)
    // }

    w.WriteHeader(http.StatusNoContent)
    // w.Write(responseBody)
}

/*
 Todo更新処理
*/
func updateTodo(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-type", "application/json")
    // トークンからuserIdを取得
	userId, err := logic.GetUserIdFromContext(r)
	if err != nil {
		// レスポンスデータ作成
		response := map[string]interface{}{
			"err": "認証エラー",
		}
		responseBody, err := json.Marshal(response)
		if err != nil {
			log.Fatal(err)
		}
		w.WriteHeader(http.StatusUnauthorized)
		w.Write(responseBody)
	}

    vars := mux.Vars(r)
    id := vars["id"]

    reqBody, _ := ioutil.ReadAll(r.Body)
    var todo models.Todo
    if err := json.Unmarshal(reqBody, &todo); err != nil {
        log.Fatal(err)
    }
    var updateTodo models.Todo
    if err := json.Unmarshal(reqBody, &updateTodo); err != nil {
        log.Fatal(err)
    }

    updateTodo.UserId = userId

    // 更新データの確認
    services.GetTodoById(&todo, id, userId)
    if todo.ID == 0 {
        // レスポンスデータ作成
		response := map[string]interface{}{
			"err": "データがありません。",
		}
		responseBody, _ := json.Marshal(response)
        w.WriteHeader(http.StatusBadRequest)
		w.Write(responseBody)
        return
    }

    // データ更新
    services.UpdateTodo(&updateTodo, id)
    convertUnitId, _ := strconv.ParseUint(id, 10, 64)
    updateTodo.BaseModel.ID = uint(convertUnitId)

    responseBody, err := json.Marshal(updateTodo)
    if err != nil {
        log.Fatal(err)
    }

    w.WriteHeader(http.StatusOK) // ステータスコード
    w.Write(responseBody)
}


func SetTodoRouting(router *mux.Router) {
	router.Handle("/todo", logic.JwtMiddleware.Handler(http.HandlerFunc(fetchAllTodos))).Methods("GET")
    router.Handle("/todo/{id}", logic.JwtMiddleware.Handler(http.HandlerFunc(fetchTodoById))).Methods("GET")

    router.Handle("/todo", logic.JwtMiddleware.Handler(http.HandlerFunc(createTodo))).Methods("POST")
    router.Handle("/todo/{id}", logic.JwtMiddleware.Handler(http.HandlerFunc(deleteTodo))).Methods("DELETE")
    router.Handle("/todo/{id}", logic.JwtMiddleware.Handler(http.HandlerFunc(updateTodo))).Methods("PUT")
}