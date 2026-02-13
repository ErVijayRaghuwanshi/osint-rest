package x

type PingResponse struct {
	Message string `json:"message" example:"Pong"`
}

type BadRequestError struct {
	Error string `json:"error" example:"username parameter is required"`
}

type NotFoundError struct {
	Error string `json:"error" example:"user not found"`
}

type InternalServerError struct {
	Error string `json:"error" example:"internal server error"`
}

type UserInfoResponse struct {
	Data struct {
		User struct {
			Result UserResult `json:"result"`
		} `json:"user"`
	} `json:"data"`
}

type UserResult struct {
	TypeName         string           `json:"__typename"`
	ID               string           `json:"id"`
	RestID           string           `json:"rest_id"`
	Core             UserCore         `json:"core"`
	Legacy           UserLegacy       `json:"legacy"`
	Avatar           UserAvatar       `json:"avatar"`
	Location         UserLocation     `json:"location"`
	IsBlueVerified   bool             `json:"is_blue_verified"`
	Verification     Verification     `json:"verification"`
	VerificationInfo VerificationInfo `json:"verification_info"`
	Privacy          Privacy          `json:"privacy"`
}

type UserCore struct {
	CreatedAt  string `json:"created_at"`
	Name       string `json:"name" example:"Mahesh Babu"`
	ScreenName string `json:"screen_name"`
}

type UserAvatar struct {
	ImageURL string `json:"image_url"`
}

type UserLocation struct {
	Location string `json:"location"`
}

type Privacy struct {
	Protected bool `json:"protected"`
}

type Verification struct {
	Verified bool `json:"verified"`
}

type VerificationInfo struct {
	IsIdentityVerified bool   `json:"is_identity_verified"`
	VerifiedSinceMsec  string `json:"verified_since_msec"`
}

type UserLegacy struct {
	Description         string `json:"description"`
	FollowersCount      int    `json:"followers_count"`
	FriendsCount        int    `json:"friends_count"`
	StatusesCount       int    `json:"statuses_count"`
	FavouritesCount     int    `json:"favourites_count"`
	MediaCount          int    `json:"media_count"`
	ProfileBannerURL    string `json:"profile_banner_url"`
	ProfileImageURL     string `json:"profile_image_url"`
	DefaultProfile      bool   `json:"default_profile"`
	DefaultProfileImage bool   `json:"default_profile_image"`
	PossiblySensitive   bool   `json:"possibly_sensitive"`
	Verified            bool   `json:"verified"`
	URL                 string `json:"url"`
}
