package api

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"
	apitool "github.com/wayne011872/api-toolkit"
	"github.com/wayne011872/api-toolkit/auth"
	apiErr "github.com/wayne011872/api-toolkit/errors"
	"github.com/wayne011872/microservice/cfg"
	"github.com/wayne011872/myTodoBackend/dao"
	"github.com/wayne011872/myTodoBackend/input"
	"github.com/wayne011872/myTodoBackend/model"
	"github.com/wayne011872/myTodoBackend/todoItem"
	"github.com/wayne011872/wayneLib/auth/perm"
)

func NewTodoItemAPI() apitool.GinAPI {
	return &TodoItemAPI{}
}

type TodoItemAPI struct {
	apiErr.CommonApiErrorHandler
}

func (a *TodoItemAPI) GetAPIs() []*apitool.GinApiHandler {
	return [] *apitool.GinApiHandler{
		{Path: "/v1/todoItem",Handler: a.getEndpoint,Method: "GET", Auth: true, Group: []auth.ApiPerm{perm.Admin, perm.Editor}},
		{Path: "/v1/todoItem",Handler: a.postEndpoint,Method: "POST", Auth: true, Group: []auth.ApiPerm{perm.Admin, perm.Editor}},
	}
}

func(a *TodoItemAPI) getEndpoint(s *gin.Context) {
	cfg,ok := cfg.GetFromGinCtx[*model.Config](s)
	if !ok {
		a.GinApiErrorHandler(s, apiErr.New(http.StatusBadRequest, "get config failed"))
	}
	conn,err := cfg.NewPgxConn(s)
	if err != nil{
		a.GinApiErrorHandler(s, apiErr.New(http.StatusBadRequest, "get config failed"))
	}

	defer conn.Close()
	queries := todoItem.New(conn.GetPgxConn())
	items, err := queries.GetAllTodoItem(s)
	if err != nil {
		a.GinApiErrorHandler(s, apiErr.New(http.StatusBadRequest, err.Error()))
		return
	}
	itemsDao := []dao.TodoItem{}
	for _, item := range items {
		itemsDao = append(itemsDao, ConvertSQLCTodoItem(item))
	}
	fmt.Println("listOK")
	s.JSON(http.StatusOK, map[string]any{
		"result": itemsDao,
	})
}

func(a *TodoItemAPI) postEndpoint(s *gin.Context) {
	in := &input.TodoItemInput{}
	err := s.BindJSON(in)
	if err != nil {
		error := apiErr.New(http.StatusBadRequest, err.Error())
		a.GinApiErrorHandler(s, error)
		return
	}
	cfg,ok := cfg.GetFromGinCtx[*model.Config](s)
	if !ok {
		a.GinApiErrorHandler(s, apiErr.New(http.StatusBadRequest, "get config failed"))
	}
	conn,err := cfg.NewPgxConn(s)
	if err != nil{
		a.GinApiErrorHandler(s, apiErr.New(http.StatusBadRequest, "get config failed"))
	}
	defer conn.Close()
	queries := todoItem.New(conn.GetPgxConn())
	err = queries.InsertTodoItem(s,todoItem.InsertTodoItemParams{
		Title: in.Title,
		Detail: pgtype.Text{
			String: in.Detail,
			Valid:  in.Detail != "",
		},
		Completed: in.Completed,
		Starttime: pgtype.Timestamp{
			Time: time.Now(),
			Valid: true,
		},
		Endtime: pgtype.Timestamp{
			Time: time.Now().Add(time.Hour * 24),
			Valid: true,
		},
	})
	if err != nil {
		a.GinApiErrorHandler(s, apiErr.New(http.StatusBadRequest, err.Error()))
		return
	}
	fmt.Println("save ok")
	s.JSON(http.StatusOK, map[string]any{
		"result": "ok",
	})
}

func ConvertSQLCTodoItem(sqlcItem todoItem.Todoitem) dao.TodoItem {
	detail:=""
	if sqlcItem.Detail.Valid {
		detail = sqlcItem.Detail.String
	}
	var createTime time.Time
	if sqlcItem.Createdtime.Valid {
		createTime = sqlcItem.Createdtime.Time
	}else{
		createTime = time.Time{}
	}

	var updateTime time.Time
	if sqlcItem.Updatedtime.Valid {
		updateTime = sqlcItem.Updatedtime.Time
	}else{
		updateTime = time.Time{}
	}

	return dao.TodoItem{
		ID:          sqlcItem.ID,
		Title:       sqlcItem.Title,
		Detail:      detail,   // Use String field of sql.NullString
		Completed:   sqlcItem.Completed,
		StartTime:   sqlcItem.Starttime.Time,
		EndTime:     sqlcItem.Endtime.Time,
		CreateTime:  createTime, // Use Time field of sql.NullTime
		UpdateTime:  updateTime, // Use Time field of sql.NullTime
	}
}