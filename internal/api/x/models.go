package x

// UserInfo represents the user information from X (Twitter)
type UserInfo struct {
	ID              string `json:"id"`
	Name            string `json:"name"`
	ScreenName      string `json:"screen_name"`
	ProfileImageURL string `json:"profile_image_url"`
	Description     string `json:"description"`
	FollowersCount  int    `json:"followers_count"`
	FollowingCount  int    `json:"following_count"`
	TweetsCount     int    `json:"tweets_count"`
	IsVerified      bool   `json:"is_verified"`
	Location        string `json:"location"`
	URL             string `json:"url"`
	CreatedAt       string `json:"created_at"`
	Protected       bool   `json:"protected"`
}

// UserInfoResponse represents the response structure for user info endpoint
type UserInfoResponse struct {
	Success bool      `json:"success"`
	Data    *UserInfo `json:"data,omitempty"`
	Error   string    `json:"error,omitempty"`
}
