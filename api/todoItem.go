package api

import (
	"fmt"
	"net/http"
	"time"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/gin-gonic/gin"
	"github.com/wayne011872/goSterna/api"
	apiErr "github.com/wayne011872/goSterna/api/err"
	"github.com/wayne011872/goSterna/db"
	"github.com/wayne011872/goSterna/log"
	"github.com/wayne011872/myTodoBackend/input"
	"github.com/wayne011872/myTodoBackend/todoItem"
	"github.com/wayne011872/myTodoBackend/dao"
)

func NewTodoItemAPI(service string) api.GinAPI {
	return &TodoItemAPI{
		ErrorOutputAPI: api.NewErrorOutputAPI(service),
	}
}

type TodoItemAPI struct {
	api.ErrorOutputAPI
}

func (a *TodoItemAPI) GetName() string{
	return "todoItem"
}

func (a *TodoItemAPI) GetAPIs() []*api.GinApiHandler {
	return [] *api.GinApiHandler{
		{Path: "/v1/todoItem",Handler: a.getEndpoint,Method: "GET"},
		{Path: "/v1/todoItem",Handler: a.postEndpoint,Method: "POST"},
	}
}

func(a *TodoItemAPI) getEndpoint(c *gin.Context) {
	logger := log.GetLogByGin(c)
	pgxClient := db.GetPgxClientByGin(c)
	conn := pgxClient.AcquireConnection(c)
	defer conn.Release()
	queries := todoItem.New(conn)
	saveLog := fmt.Sprintf("[%s] Get System Info to DB\n",time.Now().Format("2006-01-02 15:04:05"))
	logger.Info(saveLog)
	items, err := queries.GetAllTodoItem(c)
	if err != nil {
		a.GinOutputErr(c, err)
		return
	}
	itemsDao := []dao.TodoItem{}
	for _, item := range items {
		itemsDao = append(itemsDao, ConvertSQLCTodoItem(item))
	}
	fmt.Println("listOK")
	c.JSON(http.StatusOK, map[string]any{
		"result": itemsDao,
	})
}

func(a *TodoItemAPI) postEndpoint(c *gin.Context) {
	in := &input.TodoItemInput{}
	err := c.BindJSON(in)
	if err != nil {
		error := apiErr.NewApiError(http.StatusBadRequest, err.Error())
		a.GinOutputErr(c, error)
		return
	}
	logger := log.GetLogByGin(c)
	pgxClient := db.GetPgxClientByGin(c)
	queries := todoItem.New(pgxClient.AcquireConnection(c))
	saveLog := fmt.Sprintf("[%s] Save System Info to DB\n",time.Now().Format("2006-01-02 15:04:05"))
	logger.Info(saveLog)
	err = queries.InsertTodoItem(c,todoItem.InsertTodoItemParams{
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
		a.GinOutputErr(c, err)
		return
	}
	fmt.Println("save ok")
	c.JSON(http.StatusOK, map[string]any{
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