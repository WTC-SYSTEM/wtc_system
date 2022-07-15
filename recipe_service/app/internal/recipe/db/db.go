package db

import (
	"context"
	sq "github.com/Masterminds/squirrel"
	"github.com/WTC-SYSTEM/wtc_system/libs/logging"
	"github.com/WTC-SYSTEM/wtc_system/recipe_service/internal/recipe"
	"github.com/WTC-SYSTEM/wtc_system/recipe_service/pkg/client/postgresql"
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

func (d db) FindOne(ctx context.Context, id string) (recipe.Recipe, error) {

	var r recipe.Recipe
	// fill recipe
	sql, args, err := sq.
		Select("title", "description", "calories", "takes_time", "hidden").
		From("recipes").
		Where(sq.Eq{"id": id}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return recipe.Recipe{}, err
	}

	if err := d.client.QueryRow(ctx, sql, args...).
		Scan(&r.Title, &r.Description, &r.Calories, &r.TakesTime, &r.Hidden); err != nil {
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

func (d db) Update(ctx context.Context, r recipe.Recipe) error {
	// save recipe to db
	sql, args, err := sq.Update("recipes").
		Set("title", r.Title).
		Set("description", r.Description).
		Set("calories", r.Calories).
		Set("takes_time", r.TakesTime).
		Set("hidden", r.Hidden).
		Where(sq.Eq{"id": r.ID}).
		PlaceholderFormat(sq.Dollar).
		ToSql()

	if err != nil {
		return err
	}

	if _, err := d.client.Exec(ctx, sql, args...); err != nil {
		return err
	}

	// delete resources
	sql, args, err = sq.Delete("recipe_photos").
		Where(sq.Eq{"recipe_id": r.ID}).
		PlaceholderFormat(sq.Dollar).
		ToSql()

	if err != nil {
		return err
	}

	if _, err := d.client.Exec(ctx, sql, args...); err != nil {
		return err
	}

	// delete resources
	sql, args, err = sq.Delete("recipe_tags").
		Where(sq.Eq{"recipe_id": r.ID}).
		PlaceholderFormat(sq.Dollar).
		ToSql()

	if err != nil {
		return err
	}

	if _, err := d.client.Exec(ctx, sql, args...); err != nil {
		return err
	}

	// delete steps associated with recipe
	sql, args, err = sq.Delete("steps").
		Where(sq.Eq{"recipe_id": r.ID}).
		PlaceholderFormat(sq.Dollar).
		ToSql()

	if err != nil {
		return err
	}

	if _, err := d.client.Exec(ctx, sql, args...); err != nil {
		return err
	}

	// save recipe images to db
	for _, photo := range r.Photos {
		sql, args, err := sq.Insert("recipe_photos").
			Columns("recipe_id", "url").
			Values(r.ID, photo).
			PlaceholderFormat(sq.Dollar).
			ToSql()

		if err != nil {
			return err
		}
		if _, err := d.client.Exec(ctx, sql, args...); err != nil {
			return err
		}
	}

	for _, t := range r.Tags {
		sql, args, err := sq.Insert("recipe_tags").
			Columns("recipe_id", "tag").
			Values(r.ID, t).
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
			Values(r.ID, step.Title, step.Description, step.TakesTime, step.Required).
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

func (d db) Delete(ctx context.Context, id string) error {
	//TODO implement me
	panic("implement me")
}
