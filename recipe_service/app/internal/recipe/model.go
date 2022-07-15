package recipe

import (
	"github.com/WTC-SYSTEM/wtc_system/libs/utils"
	"time"
)

type Recipe struct {
	ID          string        `json:"id,omitempty"`
	Title       string        `json:"title"`
	Description string        `json:"description"`
	Calories    int           `json:"calories"`
	Steps       []Step        `json:"steps"`
	Photos      []string      `json:"photos"` // S3 Bucket URLs of photos
	Tags        []string      `json:"tags"`
	TakesTime   time.Duration `json:"takes_time"`
	Hidden      bool          `json:"hidden"`
	CreatedAt   time.Time     `json:"created_at"`
	UpdatedAt   time.Time     `json:"updated_at"`
	DeletedAt   time.Time     `json:"deleted_at"`
}

type Step struct {
	ID          string        `json:"id,omitempty"`
	RecipeID    string        `json:"recipe_id,omitempty"`
	Title       string        `json:"title"`
	Description string        `json:"description"`
	Photos      []string      `json:"photos"` // S3 Bucket URL of a photo
	TakesTime   time.Duration `json:"takes_time"`
	Required    bool          `json:"required"`
}

type CreateRecipeDTO struct {
	Title       string          `json:"title"`
	Description string          `json:"description"`
	Calories    int             `json:"calories"`
	Steps       []CreateStepDTO `json:"steps"`
	Photos      []string        `json:"photos"` // urls
	Tags        []string        `json:"tags"`
	TakesTime   time.Duration   `json:"takes_time"`
	Hidden      bool            `json:"hidden"`
}

type EditRecipeDTO struct {
	ID          string        `json:"id,omitempty"`
	Title       string        `json:"title"`
	Description string        `json:"description"`
	Calories    int           `json:"calories"`
	Steps       []EditStepDTO `json:"steps"`
	Photos      []string      `json:"photos"` // urls
	Tags        []string      `json:"tags"`
	TakesTime   time.Duration `json:"takes_time"`
	Hidden      bool          `json:"hidden"`
}

type EditStepDTO struct {
	Title       string        `json:"title"`
	Description string        `json:"description"`
	Photos      []string      `json:"photos"` // urls
	TakesTime   time.Duration `json:"takes_time"`
	Required    bool          `json:"required"`
}

type CreateStepDTO struct {
	Title       string        `json:"title"`
	Description string        `json:"description"`
	Photos      []string      `json:"photos"` // urls
	TakesTime   time.Duration `json:"takes_time"`
	Required    bool          `json:"required"`
}

type Photo struct {
	MimeType string `json:"mime_type"`
	Filename string `json:"title"`
	Size     int64  `json:"size"`
	Data     []byte `json:"data"`
}

func (s CreateStepDTO) ToStep() Step {
	return Step{
		Title:       s.Title,
		Description: s.Description,
		Photos:      s.Photos,
		TakesTime:   s.TakesTime,
		Required:    s.Required,
	}
}

func (s EditStepDTO) ToStep() Step {
	return Step{
		Title:       s.Title,
		Description: s.Description,
		Photos:      s.Photos,
		TakesTime:   s.TakesTime,
		Required:    s.Required,
	}
}

func (r CreateRecipeDTO) ToRecipe() Recipe {
	return Recipe{
		Title:       r.Title,
		Description: r.Description,
		Calories:    r.Calories,
		Steps: utils.Map(r.Steps, func(s CreateStepDTO) Step {
			return s.ToStep()
		}),
		Photos:    r.Photos,
		Tags:      r.Tags,
		TakesTime: r.TakesTime,
		Hidden:    r.Hidden,
	}
}

func (r EditRecipeDTO) ToRecipe() Recipe {
	return Recipe{
		ID:          r.ID,
		Title:       r.Title,
		Description: r.Description,
		Calories:    r.Calories,
		Steps: utils.Map(r.Steps, func(s EditStepDTO) Step {
			return s.ToStep()
		}),
		Photos:    r.Photos,
		Tags:      r.Tags,
		TakesTime: r.TakesTime,
		Hidden:    r.Hidden,
	}
}
