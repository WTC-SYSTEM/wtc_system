package recipe

type Recipe struct {
	ID          string   `json:"id"`
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Calories    int      `json:"calories"`
	Steps       []Step   `json:"steps"`
	Photos      []string `json:"photos"` // S3 Bucket URLs of photos
	Tags        []string `json:"tags"`
	TakesTime   int      `json:"takes_time"` // how much time is needed to complete this recipe in minutes
	CreatedAt   int      `json:"created_at"`
	UpdatedAt   int      `json:"updated_at"`
	DeletedAt   int      `json:"deleted_at"`
}

type Step struct {
	ID          string `json:"id"`
	RecipeID    string `json:"recipe_id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Photo       string `json:"photo"`      // S3 Bucket URL of a photo
	TakesTime   int    `json:"takes_time"` // how much time is needed to complete this step in minutes
	Required    bool   `json:"required"`
}
