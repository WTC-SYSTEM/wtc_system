package recipe

import "time"

type Recipe struct {
	ID          string        `json:"id"`
	Title       string        `json:"title"`
	Description string        `json:"description"`
	Calories    int           `json:"calories"`
	Steps       []Step        `json:"steps"`
	Photos      []string      `json:"photos"` // S3 Bucket URLs of photos
	Tags        []string      `json:"tags"`
	TakesTime   time.Duration `json:"takes_time"` // how much time is needed to complete this recipe in minutes
	CreatedAt   int           `json:"created_at"`
	UpdatedAt   int           `json:"updated_at"`
	DeletedAt   int           `json:"deleted_at"`
}

type Step struct {
	ID          string        `json:"id"`
	RecipeID    string        `json:"recipe_id"`
	Title       string        `json:"title"`
	Description string        `json:"description"`
	Photos      []string      `json:"photos"`     // S3 Bucket URL of a photo
	TakesTime   time.Duration `json:"takes_time"` // how much time is needed to complete this step in minutes
	Required    bool          `json:"required"`
}

type CreateRecipeDTO struct {
	Title       string          `json:"title"`
	Description string          `json:"description"`
	Calories    int             `json:"calories"`
	Steps       []CreateStepDTO `json:"steps"`
	Photos      []Photo         `json:"photos"` // raw photos
	Tags        []string        `json:"tags"`
	TakesTime   time.Duration   `json:"takes_time"` // how much time is needed to complete this recipe in minutes
}

type CreateStepDTO struct {
	Title       string        `json:"title"`
	Description string        `json:"description"`
	Photos      []Photo       `json:"photos"`     // raw photo
	TakesTime   time.Duration `json:"takes_time"` // how much time is needed to complete this step in minutes
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
		Photos:      []string{},
		TakesTime:   s.TakesTime,
		Required:    s.Required,
	}
}

func (r CreateRecipeDTO) ToRecipe() Recipe {
	return Recipe{
		Title:       r.Title,
		Description: r.Description,
		Calories:    r.Calories,
		Steps:       []Step{},
		Photos:      []string{},
		Tags:        r.Tags,
		TakesTime:   r.TakesTime,
	}
}
