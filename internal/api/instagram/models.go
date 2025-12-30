package instagram



import "encoding/json"

// UserInfoResponse represents the top-level response from Instagram API
type UserInfoResponse struct {
	Data struct {
		User UserInfo `json:"user"`
	} `json:"data"`
}

// UserInfo represents detailed information about an Instagram user
type UserInfo struct {
	AIAgentOwnerUsername       *string               `json:"ai_agent_owner_username"`
	AIAgentType                *string               `json:"ai_agent_type"`
	BioLinks                   []BioLink             `json:"bio_links"`
	Biography                  string                `json:"biography"`
	BiographyWithEntities      BiographyWithEntities `json:"biography_with_entities"`
	BlockedByViewer            bool                  `json:"blocked_by_viewer"`
	BusinessAddressJSON        *string               `json:"business_address_json"`
	BusinessCategoryName       *string               `json:"business_category_name"`
	BusinessContactMethod      string                `json:"business_contact_method"`
	BusinessEmail              *string               `json:"business_email"`
	BusinessPhoneNumber        *string               `json:"business_phone_number"`
	CategoryEnum               *string               `json:"category_enum"`
	CategoryName               *string               `json:"category_name"`
	CountryBlock               bool                  `json:"country_block"`
	// EdgeFelixVideoTimeline     EdgeConnection        `json:"edge_felix_video_timeline"`
	// EdgeFollow                 EdgeCount             `json:"edge_follow"`
	// EdgeFollowedBy             EdgeCount             `json:"edge_followed_by"`
	// EdgeMediaCollections       EdgeConnection        `json:"edge_media_collections"`
	// EdgeMutualFollowedBy       EdgeMutualFollowedBy  `json:"edge_mutual_followed_by"`
	// EdgeOwnerToTimelineMedia   EdgeConnection        `json:"edge_owner_to_timeline_media"`
	// EdgeRelatedProfiles        EdgeRelatedProfiles   `json:"edge_related_profiles"`
	// EdgeSavedMedia             EdgeConnection        `json:"edge_saved_media"`
	FBPageCallsCount           *int                  `json:"fb_page_calls_count"`
	FBPageEmail                *string               `json:"fb_page_email"`
	FollowedByViewer           bool                  `json:"followed_by_viewer"`
	FollowsViewer              bool                  `json:"follows_viewer"`
	FullName                   string                `json:"full_name"`
	GroupMetadata              *string               `json:"group_metadata"`
	HasAdjustedAds             *bool                 `json:"has_adjusted_ads"`
	HasArEffects               bool                  `json:"has_ar_effects"`
	HasBlockedViewer           bool                  `json:"has_blocked_viewer"`
	HasChannel                 bool                  `json:"has_channel"`
	HasGuides                  bool                  `json:"has_guides"`
	HasRequestedViewer         bool                  `json:"has_requested_viewer"`
	ID                         string                `json:"id"`
	InstafilmsCount            *int                  `json:"instafilms_count"`
	IsPrivate                  bool                  `json:"is_private"`
	IsProfilePrivate           bool                  `json:"is_profile_private"`
	IsVerified                 bool                  `json:"is_verified"`
	IsVerifiedByMV2            bool                  `json:"is_verified_by_mv2"`
	IgCustomerSalesChannel     bool                  `json:"ig_custom_customer_sales_channel"`
	InstagramLocation          *InstagramLocation    `json:"instagram_location"`
	JoinedRecently             bool                  `json:"joined_recently"`
	LikelyHasAudience          bool                  `json:"likely_has_audience"`
	MediaCount                 int                   `json:"media_count"`
	MutualFollowersCount       int                   `json:"mutual_followers_count"`
	NumberOfLinksInBio         int                   `json:"number_of_links_in_bio"`
	PageIDForAds               *string               `json:"page_id_for_ads"`
	PageNumber                 int                   `json:"page_number"`
	ProfileContext             *string               `json:"profile_context"`
	ProfilePicID               string                `json:"profile_pic_id"`
	ProfilePicUrl              string                `json:"profile_pic_url"`
	ProfilePicUrlHD            string                `json:"profile_pic_url_hd"`
	ProfileViewerFollowButton  bool                  `json:"profile_viewer_follow_button"`
	PromoShareUserFollowButton bool                  `json:"promo_share_user_follow_button"`
	RequestedByViewer          bool                  `json:"requested_by_viewer"`
	RestrictedByViewer         bool                  `json:"restricted_by_viewer"`
	ShouldShowCategory         bool                  `json:"should_show_category"`
	ShouldShowPublicContacts   bool                  `json:"should_show_public_contacts"`
	SocialContext              *string               `json:"social_context"`
	SocialContextList          []string              `json:"social_context_list"`
	StatusEmojis               []string              `json:"status_emojis"`
	Title                      *string               `json:"title"`
	TotalCLIPCount             int                   `json:"total_clip_count"`
	Username                   string                `json:"username"`
	VerificationData           *VerificationData     `json:"verification_data"`
	WatchedStories             []string              `json:"watched_stories"`
}

// BioLink represents a link in the user's bio
type BioLink struct {
	LinkType string `json:"link_type"`
	LynxURL  string `json:"lynx_url"`
	Title    string `json:"title"`
	URL      string `json:"url"`
}

// BiographyWithEntities represents the biography with parsed entities
type BiographyWithEntities struct {
	Entities []interface{} `json:"entities"` // Can be parsed into specific entity types if needed
	RawText  string        `json:"raw_text"`
}

// EdgeConnection represents a connection with edges (like timeline, collections, etc.)
type EdgeConnection struct {
	Count    int           `json:"count"`
	Edges    []interface{} `json:"edges"`
	PageInfo PageInfo      `json:"page_info"`
}

// PageInfo represents pagination information
type PageInfo struct {
	EndCursor   *string `json:"end_cursor"`
	HasNextPage bool    `json:"has_next_page"`
}

// EdgeCount represents a simple count edge (like followers, following)
type EdgeCount struct {
	Count int `json:"count"`
}

// EdgeMutualFollowedBy represents mutual followed by information
type EdgeMutualFollowedBy struct {
	Count int           `json:"count"`
	Edges []interface{} `json:"edges"`
}

// EdgeRelatedProfiles represents related profiles
type EdgeRelatedProfiles struct {
	Edges []interface{} `json:"edges"`
}

// InstagramLocation represents location information
type InstagramLocation struct {
	AddressJSON *string `json:"address_json"`
	HasBeacon   bool    `json:"has_beacon"`
	LocationID  string  `json:"location_id"`
	Name        string  `json:"name"`
	Slug        string  `json:"slug"`
	Website     *string `json:"website"`
}

// VerificationData represents verification information
type VerificationData struct {
	VerificationMethod *string `json:"verification_method"`
	VerificationStatus *string `json:"verification_status"`
}

// HashTagSearchResponse is the top-level GraphQL hashtag search response
type HashTagSearchResponse struct {
	Data       SearchData `json:"data"`
	Extensions Extensions `json:"extensions,omitempty"`
	Status     string     `json:"status"`
}

// SearchData wraps the top search connection
type SearchData struct {
	TopSearch TopSearchConnection `json:"xdt_api__v1__fbsearch__topsearch_connection"`
}

// TopSearchConnection contains hashtag, user, and place results
type TopSearchConnection struct {
	Hashtags     []HashtagResult  `json:"hashtags,omitempty"`
	Places       []any            `json:"places,omitempty"`
	Users        []any            `json:"users,omitempty"`
	RankToken    string           `json:"rank_token"`

	// Nullable / schema-flexible fields
	SeeMore      json.RawMessage  `json:"see_more,omitempty"`
	InformModule json.RawMessage  `json:"inform_module,omitempty"`
}

// HashtagResult represents a ranked hashtag entry
type HashtagResult struct {
	Position int     `json:"position"`
	Hashtag  Hashtag `json:"hashtag"`
}

// Hashtag contains hashtag metadata
type Hashtag struct {
	ID                   string `json:"id"`
	Name                 string `json:"name"`
	MediaCount           int    `json:"media_count"`
	Subtitle             string `json:"search_result_subtitle"`
}

// Extensions contains GraphQL execution metadata
type Extensions struct {
	IsFinal        bool           `json:"is_final"`
	ServerMetadata ServerMetadata `json:"server_metadata"`
}

// ServerMetadata contains request timing details
type ServerMetadata struct {
	RequestStartTimeMS int64 `json:"request_start_time_ms"`
	TimeAtFlushMS      int64 `json:"time_at_flush_ms"`
}