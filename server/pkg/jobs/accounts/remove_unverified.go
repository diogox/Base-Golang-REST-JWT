package accounts

import (
	"context"
	"github.com/diogox/REST-JWT/server"
	"github.com/labstack/echo"
)

func NewRemoveIfUnverifiedAccountJob(logger echo.Logger, db server.SqlDB, userId string) RemoveIfUnverifiedAccountJob {
	return RemoveIfUnverifiedAccountJob{
		logger: logger,
		db:     db,
		userId: userId,
	}
}

type RemoveIfUnverifiedAccountJob struct {
	logger echo.Logger
	db     server.SqlDB
	userId string
}

func (j RemoveIfUnverifiedAccountJob) Run() {
	ctx := context.Background()

	user, err :=j.db.GetUserByID(ctx, j.userId)
	if err != nil {
		j.logger.Error(err.Error())
		return
	}

	if !user.IsEmailVerified {
		_, err = j.db.DeleteUserByID(ctx, j.userId)
		if err != nil {
			j.logger.Error(err.Error())
			return
		}
	}

	j.logger.Info("Successfully deleted unverified user!")
}
