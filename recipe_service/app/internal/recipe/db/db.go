package db

import (
	"bytes"
	"context"
	"fmt"
	sq "github.com/Masterminds/squirrel"
	"github.com/WTC-SYSTEM/wtc_system/recipe_service/internal/recipe"
	"github.com/WTC-SYSTEM/wtc_system/recipe_service/pkg/client/aws"
	"github.com/WTC-SYSTEM/wtc_system/recipe_service/pkg/client/postgresql"
	"github.com/WTC-SYSTEM/wtc_system/recipe_service/pkg/logging"
	"github.com/aws/aws-sdk-go/service/s3"
	"strconv"
	"strings"
	"time"
)

type db struct {
	client postgresql.Client
	logger logging.Logger
	awsCfg aws.Aws
}

func NewStorage(c postgresql.Client, l logging.Logger, awsCfg aws.Aws) recipe.Storage {
	return &db{
		client: c,
		logger: l,
		awsCfg: awsCfg,
	}
}

func (d db) UploadFile(ctx context.Context, photo recipe.Photo, folder string) (string, error) {
	fName := folder + "/" + strings.ToLower(strings.ReplaceAll(photo.Filename, " ", "")) + "-" + strconv.FormatInt(time.Now().Unix(), 10) + ".jpg"
	params := &s3.PutObjectInput{
		Bucket:      &d.awsCfg.Config.Bucket,
		Key:         &fName,
		Body:        bytes.NewReader(photo.Data),
		ContentType: &photo.MimeType,
	}
	_, err := d.awsCfg.S3.PutObject(params)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s", d.awsCfg.Config.Bucket, d.awsCfg.Config.Region, fName), nil
}

func (d db) Create(ctx context.Context, r recipe.Recipe) error {
	// save recipe to db
	sql, args, err := sq.Insert("recipes").
		Columns("title", "description", "calories", "takes_time").
		Values(r.Title, r.Description, r.Calories, r.TakesTime).
		PlaceholderFormat(sq.Dollar).
		Suffix("RETURNING id").
		ToSql()

	if err != nil {
		return err
	}
	var id string

	if err := d.client.QueryRow(ctx, sql, args...).Scan(&id); err != nil {
		return err
	}

	// save recipe images to db
	for _, photo := range r.Photos {
		sql, args, err := sq.Insert("recipe_photos").
			Columns("recipe_id", "url").
			Values(id, photo).
			PlaceholderFormat(sq.Dollar).
			ToSql()

		if err != nil {
			return err
		}
		if _, err := d.client.Exec(ctx, sql, args...); err != nil {
			return err
		}
	}

	// save steps to db with recipe id
	for _, step := range r.Steps {
		sql, args, err := sq.Insert("steps").
			Columns("recipe_id", "title", "description", "takes_time", "required").
			Values(id, step.Title, step.Description, step.TakesTime, step.Required).
			PlaceholderFormat(sq.Dollar).
			Suffix("RETURNING id").
			ToSql()

		if err != nil {
			return err
		}
		var id string

		if err := d.client.QueryRow(ctx, sql, args...).Scan(&id); err != nil {
			return err
		}

		// save step photos to db
		for _, photo := range step.Photos {
			sql, args, err := sq.Insert("step_photos").
				Columns("step_id", "url").
				Values(id, photo).
				PlaceholderFormat(sq.Dollar).
				ToSql()

			if err != nil {
				return err
			}
			if _, err := d.client.Exec(ctx, sql, args...); err != nil {
				return err
			}
		}

	}

	return nil
}

func (d db) FindOne(ctx context.Context, id string) error {
	//TODO implement me
	panic("implement me")
}

func (d db) Update(ctx context.Context) error {
	//TODO implement me
	panic("implement me")
}

func (d db) Delete(ctx context.Context, id string) error {
	//TODO implement me
	panic("implement me")
}
