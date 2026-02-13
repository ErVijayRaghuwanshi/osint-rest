package snapchat

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

// ========================================================
//
//	User Info Models
//
// ========================================================
// UserInfoResponse represents the public profile information of a Snapchat user
type UserInfoResponse struct {
	Case              string    `json:"$case"`
	PublicProfileInfo *UserInfo `json:"publicProfileInfo"`
}

type UserInfo struct {
	Username               string               `json:"username"`
	Title                  string               `json:"title"`
	SnapcodeImageURL       string               `json:"snapcodeImageUrl"`
	Badge                  int                  `json:"badge"`
	CategoryStringID       string               `json:"categoryStringId"`
	SubcategoryStringID    string               `json:"subcategoryStringId"`
	SubscriberCount        string               `json:"subscriberCount"`
	Bio                    string               `json:"bio"`
	WebsiteURL             string               `json:"websiteUrl"`
	ProfilePictureURL      string               `json:"profilePictureUrl"`
	Address                string               `json:"address"`
	Bitmoji3D              interface{}          `json:"bitmoji3d"` // always null in sample
	HasCuratedHighlights   bool                 `json:"hasCuratedHighlights"`
	HasSpotlightHighlights bool                 `json:"hasSpotlightHighlights"`
	MutableName            string               `json:"mutableName"`
	PublisherType          string               `json:"publisherType"`
	SquareHeroImageURL     string               `json:"squareHeroImageUrl"`
	PrimaryColor           string               `json:"primaryColor"`
	HasStory               bool                 `json:"hasStory"`
	RelatedAccountsInfo    []RelatedAccountInfo `json:"relatedAccountsInfo"`
	CreationTimestampMs    *TimestampWrapper    `json:"creationTimestampMs"`
	LastUpdateTimestampMs  *TimestampWrapper    `json:"lastUpdateTimestampMs"`
	BusinessProfileID      string               `json:"businessProfileId"`
	SameAsLinks            []string             `json:"sameAsLinks"`
	ShouldHideUsername     bool                 `json:"shouldHideUsername"`
}

type RelatedAccountInfo struct {
	PublicProfileInfo *interface{}   `json:"publicProfileInfo"`
	SubscribeLink     *SubscribeLink `json:"subscribeLink"`
}

type SubscribeLink struct {
	OneLinkBaseURL        string   `json:"oneLinkBaseUrl"`
	PIDKeys               []string `json:"pidKeys"`
	PIDFallbackValue      string   `json:"pidFallbackValue"`
	CampaignKeys          []string `json:"campaignKeys"`
	CampaignFallbackValue string   `json:"campaignFallbackValue"`
	GoogleClickIDParam    string   `json:"googleClickIdParam"`
	DeepLinkURL           string   `json:"deepLinkUrl"`
	IOSAppStoreURL        string   `json:"iosAppStoreUrl"`
	DesktopPageURL        string   `json:"desktopPageUrl"`
}

type TimestampWrapper struct {
	Value string `json:"value"`
}

// ========================================================
//						Highlights Models
// ========================================================

// SpotlightHighlight represents a Snapchat spotlight highlight object
type SpotlightHighlight struct {
	CanonicalURLSuffix *string         `json:"canonicalUrlSuffix"`
	Emoji              *string         `json:"emoji"`
	HighlightID        *string         `json:"highlightId"`
	IsAttributed       *BoolValue      `json:"isAttributed"`
	SnapList           []SpotlightSnap `json:"snapList"`
	StoryID            *StringValue    `json:"storyId"`
	StoryShareID       *string         `json:"storyShareId"`
	StorySubtitle      *string         `json:"storySubtitle"`
	StoryTapID         string          `json:"storyTapId"`
	StoryTitle         *string         `json:"storyTitle"`
	StoryType          int             `json:"storyType"`
	ThumbnailURL       *StringValue    `json:"thumbnailUrl"`
	VideoTrackURL      *string         `json:"videoTrackUrl"`
}

type SpotlightSnap struct {
	AudioTranscriptionObjectURL *string      `json:"audioTranscriptionObjectUrl"`
	HasAttachment               *bool        `json:"hasAttachment"`
	IntervalStartTimeMs         *int64       `json:"intervalStartTimeMs"`
	IsSponsored                 *bool        `json:"isSponsored"`
	Lat                         *float64     `json:"lat"`
	Lng                         *float64     `json:"lng"`
	SnapID                      *StringValue `json:"snapId"`
	SnapIndex                   int          `json:"snapIndex"`
	SnapMediaType               int          `json:"snapMediaType"`
	SnapSubtitles               *string      `json:"snapSubtitles"`
	SnapTitle                   *string      `json:"snapTitle"`
	SnapURLs                    SnapURLs     `json:"snapUrls"`
	TimestampInSec              *StringValue `json:"timestampInSec"`
}

// CuratedHighlightsResponse represents the response for curated highlights
type CuratedHighlightsResponse struct {
	Highlights []CuratedHighlight `json:"highlights"`
}

// CuratedHighlight represents a Snapchat curated highlight
type CuratedHighlight struct {
	CanonicalURLSuffix *string       `json:"canonicalUrlSuffix"`
	Emoji              *string       `json:"emoji"`
	HighlightID        *StringValue  `json:"highlightId"`
	IsAttributed       *BoolValue    `json:"isAttributed"`
	SnapList           []CuratedSnap `json:"snapList"`
	StoryTapID         *string       `json:"storyTapId,omitempty"`
	StoryTitle         *string       `json:"storyTitle,omitempty"`
	StorySubtitle      *string       `json:"storySubtitle,omitempty"`
	StoryType          *int          `json:"storyType,omitempty"`
	ThumbnailURL       *StringValue  `json:"thumbnailUrl,omitempty"`
	VideoTrackURL      *string       `json:"videoTrackUrl,omitempty"`
}
type CuratedSnap struct {
	AudioTranscriptionObjectURL *string      `json:"audioTranscriptionObjectUrl"`
	HasAttachment               *bool        `json:"hasAttachment"`
	IntervalStartTimeMs         *int64       `json:"intervalStartTimeMs"`
	IsSponsored                 *bool        `json:"isSponsored"`
	Lat                         *float64     `json:"lat"`
	Lng                         *float64     `json:"lng"`
	SnapID                      *StringValue `json:"snapId"`
	SnapIndex                   int          `json:"snapIndex"`
	SnapMediaType               int          `json:"snapMediaType"`
	SnapSubtitles               *string      `json:"snapSubtitles"`
	SnapTitle                   *string      `json:"snapTitle"`
	SnapURLs                    SnapURLs     `json:"snapUrls"`
	TimestampInSec              *StringValue `json:"timestampInSec"`
}

// Commnon types used in both Curated and Spotlight highlights
type SnapURLs struct {
	AttachmentURL   *string      `json:"attachmentUrl"`
	MediaPreviewURL *StringValue `json:"mediaPreviewUrl"`
	MediaURL        *string      `json:"mediaUrl"`
	OverlayURL      *string      `json:"overlayUrl"`
}
type StringValue struct {
	Value string `json:"value"`
}

type BoolValue struct {
	Value bool `json:"value"`
}
