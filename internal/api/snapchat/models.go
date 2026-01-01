package snapchat



// import "encoding/json"

// read input json from next_data_json['props']['pageProps']['userProfile']


// UserInfoResponse represents the public profile information of a Snapchat user
type UserInfoResponse struct {
	Case              string             `json:"$case"`
	PublicProfileInfo *UserInfo          `json:"publicProfileInfo"`
}



type UserInfo struct {
	Username                 string                 `json:"username"`
	Title                    string                 `json:"title"`
	SnapcodeImageURL         string                 `json:"snapcodeImageUrl"`
	Badge                    int                    `json:"badge"`
	CategoryStringID         string                 `json:"categoryStringId"`
	SubcategoryStringID      string                 `json:"subcategoryStringId"`
	SubscriberCount          string                 `json:"subscriberCount"`
	Bio                      string                 `json:"bio"`
	WebsiteURL               string                 `json:"websiteUrl"`
	ProfilePictureURL        string                 `json:"profilePictureUrl"`
	Address                  string                 `json:"address"`
	Bitmoji3D                interface{}            `json:"bitmoji3d"` // always null in sample
	HasCuratedHighlights     bool                   `json:"hasCuratedHighlights"`
	HasSpotlightHighlights   bool                   `json:"hasSpotlightHighlights"`
	MutableName              string                 `json:"mutableName"`
	PublisherType            string                 `json:"publisherType"`
	SquareHeroImageURL       string                 `json:"squareHeroImageUrl"`
	PrimaryColor             string                 `json:"primaryColor"`
	HasStory                 bool                   `json:"hasStory"`
	RelatedAccountsInfo      []RelatedAccountInfo   `json:"relatedAccountsInfo"`
	CreationTimestampMs      *TimestampWrapper      `json:"creationTimestampMs"`
	LastUpdateTimestampMs    *TimestampWrapper      `json:"lastUpdateTimestampMs"`
	BusinessProfileID        string                 `json:"businessProfileId"`
	SameAsLinks              []string               `json:"sameAsLinks"`
	ShouldHideUsername       bool                   `json:"shouldHideUsername"`
}


type RelatedAccountInfo struct {
	PublicProfileInfo *interface{} `json:"publicProfileInfo"`
	SubscribeLink     *SubscribeLink `json:"subscribeLink"`
}


type SubscribeLink struct {
	OneLinkBaseURL         string   `json:"oneLinkBaseUrl"`
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
