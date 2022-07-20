package db

import (
	"context"
	sq "github.com/Masterminds/squirrel"
	"github.com/WTC-SYSTEM/logging"
	"github.com/WTC-SYSTEM/wtc_system/recipe_service/internal/recipe"
	"github.com/WTC-SYSTEM/wtc_system/recipe_service/pkg/client/postgresql"
	"time"
)

type db struct {
	client postgresql.Client
	logger logging.Logger
}

func NewStorage(c postgresql.Client, l logging.Logger) recipe.Storage {
	return &db{
		client: c,
		logger: l,
	}
}

/*
FindAll
*/

func (d db) FindAll(ctx context.Context, filter recipe.Filter) ([]recipe.Recipe, error) {
	// get recipes with filter
	sql, args, _ := sq.Select(
		"title",
		"description",
		"calories",
		"takes_time",
		"hidden",
		"created_at",
		"updated_at",
	).
		From("recipes").
		Where(sq.Eq{"hidden": false, "deleted_at": ""}).
		Where(sq.Expr("id IN (SELECT recipe_id FROM recipe_tags WHERE tag IN ($1)) OR $1 IS NULL", filter.Tags)).
		Offset(uint64(filter.Limit * filter.Page)).
		PlaceholderFormat(sq.Dollar).
		ToSql()

	d.logger.Info(sql, args)
	return []recipe.Recipe{}, nil
}

/*
Create
*/

func (d db) Create(ctx context.Context, r recipe.Recipe) (string, error) {
	// save recipe to db
	sql, args, err := sq.Insert("recipes").
		Columns("title", "description", "calories", "takes_time", "hidden").
		Values(r.Title, r.Description, r.Calories, r.TakesTime, r.Hidden).
		PlaceholderFormat(sq.Dollar).
		Suffix("RETURNING id").
		ToSql()

	if err != nil {
		return "", err
	}
	var id string

	if err := d.client.QueryRow(ctx, sql, args...).Scan(&id); err != nil {
		return "", err
	}

	// save recipe images to db
	for _, photo := range r.Photos {
		sql, args, err := sq.Insert("recipe_photos").
			Columns("recipe_id", "url").
			Values(id, photo).
			PlaceholderFormat(sq.Dollar).
			ToSql()

		if err != nil {
			return "", err
		}
		if _, err := d.client.Exec(ctx, sql, args...); err != nil {
			return "", err
		}
	}

	for _, t := range r.Tags {
		sql, args, err := sq.Insert("recipe_tags").
			Columns("recipe_id", "tag").
			Values(id, t).
			PlaceholderFormat(sq.Dollar).
			ToSql()

		if err != nil {
			return "", err
		}
		if _, err := d.client.Exec(ctx, sql, args...); err != nil {
			return "", err
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
			return "", err
		}
		var id string

		if err := d.client.QueryRow(ctx, sql, args...).Scan(&id); err != nil {
			return "", err
		}

		// save step photos to db
		for _, photo := range step.Photos {
			sql, args, err := sq.Insert("step_photos").
				Columns("step_id", "url").
				Values(id, photo).
				PlaceholderFormat(sq.Dollar).
				ToSql()

			if err != nil {
				return "", err
			}
			if _, err := d.client.Exec(ctx, sql, args...); err != nil {
				return "", err
			}
		}

	}

	return id, nil
}

/*
FindOne
*/
func (d db) FindOne(ctx context.Context, id string) (recipe.Recipe, error) {

	var r recipe.Recipe
	// fill recipe
	sql, args, err := sq.
		Select("title", "description", "calories", "takes_time", "hidden", "created_at", "updated_at").
		From("recipes").
		Where(sq.Eq{"id": id}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return recipe.Recipe{}, err
	}

	if err := d.client.QueryRow(ctx, sql, args...).
		Scan(&r.Title, &r.Description, &r.Calories, &r.TakesTime, &r.Hidden, &r.CreatedAt, &r.UpdatedAt); err != nil {
		return recipe.Recipe{}, err
	}

	// fill steps
	sql, args, err = sq.
		Select("s.id", "s.title", "s.description", "s.takes_time", "s.required").
		From("recipes").
		LeftJoin("steps s ON recipes.id=s.recipe_id").
		Where(sq.Eq{"recipes.id": id}).
		PlaceholderFormat(sq.Dollar).
		ToSql()

	if err != nil {
		return recipe.Recipe{}, err
	}

	rows, err := d.client.Query(ctx, sql, args...)
	if err != nil {
		return recipe.Recipe{}, err
	}

	for rows.Next() {
		var step recipe.Step
		if err := rows.Scan(&step.ID, &step.Title, &step.Description, &step.TakesTime, &step.Required); err != nil {
			return recipe.Recipe{}, err
		}

		sql, args, err := sq.
			Select("url").
			From("step_photos").
			Where(sq.Eq{"step_id": step.ID}).
			PlaceholderFormat(sq.Dollar).
			ToSql()

		if err != nil {
			return recipe.Recipe{}, err
		}

		rows, err := d.client.Query(ctx, sql, args...)
		if err != nil {
			return recipe.Recipe{}, err
		}
		for rows.Next() {
			var url string
			if err := rows.Scan(&url); err != nil {
				return recipe.Recipe{}, err
			}
			step.Photos = append(step.Photos, url)
		}
		r.Steps = append(r.Steps, step)
	}

	// fill recipe tags
	sql, args, err = sq.
		Select("tag").
		From("recipe_tags").
		Where(sq.Eq{"recipe_id": id}).
		PlaceholderFormat(sq.Dollar).
		ToSql()

	if err != nil {
		return recipe.Recipe{}, err
	}

	rows, err = d.client.Query(ctx, sql, args...)
	if err != nil {
		return recipe.Recipe{}, err
	}

	for rows.Next() {
		var tag string
		if err := rows.Scan(&tag); err != nil {
			return recipe.Recipe{}, err
		}
		r.Tags = append(r.Tags, tag)
	}

	// fill recipe photos
	sql, args, err = sq.
		Select("url").
		From("recipe_photos").
		Where(sq.Eq{"recipe_id": id}).
		PlaceholderFormat(sq.Dollar).
		ToSql()

	if err != nil {
		return recipe.Recipe{}, err
	}

	rows, err = d.client.Query(ctx, sql, args...)
	if err != nil {
		return recipe.Recipe{}, err
	}

	for rows.Next() {
		var url string
		if err := rows.Scan(&url); err != nil {
			return recipe.Recipe{}, err
		}
		r.Photos = append(r.Photos, url)
	}
	return r, nil
}

/*
Update
*/
func (d db) Update(ctx context.Context, r recipe.Recipe) error {
	updBuilder := sq.Update("recipes")

	if r.Title != "" {
		updBuilder = updBuilder.Set("title", r.Title)
	}

	if r.Description != "" {
		updBuilder = updBuilder.Set("description", r.Description)
	}

	if r.Calories != 0 {
		updBuilder = updBuilder.Set("calories", r.Calories)
	}

	if r.TakesTime != 0 {
		updBuilder = updBuilder.Set("takes_time", r.TakesTime)
	}

	updBuilder = updBuilder.Set("hidden", r.Hidden)

	// Update fields in database for recipe
	sql, args, _ := updBuilder.
		Where(sq.Eq{"id": r.ID}).
		PlaceholderFormat(sq.Dollar).
		ToSql()

	if _, err := d.client.Exec(ctx, sql, args...); err != nil {
		return err
	}
	/*
		Recipe Photos
	*/
	if len(r.Photos) > 0 {
		// delete recipe photos
		sql, args, _ = sq.Delete("recipe_photos").
			Where(sq.Eq{"recipe_id": r.ID}).
			PlaceholderFormat(sq.Dollar).
			ToSql()

		if _, err := d.client.Exec(ctx, sql, args...); err != nil {
			return err
		}

		// save recipe photos to db
		for _, photo := range r.Photos {
			sql, args, _ := sq.Insert("recipe_photos").
				Columns("recipe_id", "url").
				Values(r.ID, photo).
				PlaceholderFormat(sq.Dollar).
				ToSql()

			if _, err := d.client.Exec(ctx, sql, args...); err != nil {
				return err
			}
		}
	}

	/*
		Recipe Tags
	*/
	if len(r.Tags) > 0 {
		// delete recipe tags
		sql, args, _ = sq.Delete("recipe_tags").
			Where(sq.Eq{"recipe_id": r.ID}).
			PlaceholderFormat(sq.Dollar).
			ToSql()

		if _, err := d.client.Exec(ctx, sql, args...); err != nil {
			return err
		}
		// create recipe tags
		for _, t := range r.Tags {
			sql, args, _ := sq.Insert("recipe_tags").
				Columns("recipe_id", "tag").
				Values(r.ID, t).
				PlaceholderFormat(sq.Dollar).
				ToSql()

			if _, err := d.client.Exec(ctx, sql, args...); err != nil {
				return err
			}
		}
	}
	/*
		Recipe Steps
	*/
	if len(r.Steps) > 0 {
		/*
			Delete all steps associated with recipe
			This will also delete all photos associated with step
		*/
		sql, args, _ = sq.Delete("steps").
			Where(sq.Eq{"recipe_id": r.ID}).
			PlaceholderFormat(sq.Dollar).
			ToSql()

		if _, err := d.client.Exec(ctx, sql, args...); err != nil {
			return err
		}

		/*
			Save steps and steps photos to db
		*/
		for _, step := range r.Steps {
			// save step to db
			sql, args, _ := sq.Insert("steps").
				Columns("recipe_id", "title", "description", "takes_time", "required").
				Values(r.ID, step.Title, step.Description, step.TakesTime, step.Required).
				PlaceholderFormat(sq.Dollar).
				Suffix("RETURNING id").
				ToSql()

			var id string

			if err := d.client.QueryRow(ctx, sql, args...).Scan(&id); err != nil {
				return err
			}

			// save step photos to db
			for _, photo := range step.Photos {
				sql, args, _ := sq.Insert("step_photos").
					Columns("step_id", "url").
					Values(id, photo).
					PlaceholderFormat(sq.Dollar).
					ToSql()

				if _, err := d.client.Exec(ctx, sql, args...); err != nil {
					return err
				}
			}

		}
	}

	return nil
}

/*
Delete
*/
func (d db) Delete(ctx context.Context, id string) error {
	/*
		Actually, we do not delete it.
		We just set the deleted_at to current timestamp.
	*/
	sql, args, _ := sq.Update("recipes").
		Set("deleted_at", time.Now()).
		Where(sq.Eq{"id": id}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	_, err := d.client.Exec(ctx, sql, args...)
	return err
}
