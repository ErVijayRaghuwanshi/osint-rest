package jaco

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
	Error     string       `json:"error" example:"Success"`
	Data      UserInfoData `json:"data"`
	ErrorCode int          `json:"error_code" example:"10000"`
	RequestID string       `json:"request_id" example:"01f1ee42-a913-4ad2-832e-95e6fb6c8dec"`
}

type UserInfoData struct {
	User           User         `json:"user"`
	UserInfoStatus int          `json:"user_info_status" example:"0"`
	ShowTab        ShowTab      `json:"show_tab"`
	CountInfo      CountInfo    `json:"count_info"`
	BadgeGallery   BadgeGallery `json:"badge_gallery"`
	UserBusiness   UserBusiness `json:"user_business_info"`
	UserMeta       UserMeta     `json:"user_meta"`
	Extra          Extra        `json:"extra"`
}

type User struct {
	ID               string   `json:"id" example:"3004531333"`
	Nickname         string   `json:"nickname" example:"شمسة"`
	Username         string   `json:"username" example:"Uaegirl"`
	Country          string   `json:"country" example:""`
	Province         string   `json:"province" example:""`
	City             string   `json:"city" example:""`
	Region           string   `json:"region" example:"0"`
	AvatarHDPID      string   `json:"avatar_hd_pid" example:"b3158285la1i906ie78z1j20yi0yidj2"`
	CoverPID         string   `json:"cover_pid" example:"b3158285la1i906ie78z1j20yi0yidj2_p,,,,"`
	Gender           int      `json:"gender" example:"2"`
	Birthday         string   `json:"birthday" example:""`
	Description      string   `json:"description" example:""`
	CreatedAt        string   `json:"created_at" example:"1691843311000"`
	UserType         int      `json:"user_type" example:"1"`
	RegSource        string   `json:"reg_source" example:"0"`
	UpdateTime       string   `json:"update_time" example:"1767800276000"`
	FollowerCount    int      `json:"follower_count" example:"32827"`
	FollowingCount   int      `json:"following_count" example:"1155"`
	StatusesCount    int      `json:"statuses_count" example:"0"`
	FriendCount      int      `json:"friend_count" example:"1068"`
	Privacy          string   `json:"privacy" example:"345"`
	TagList          []string `json:"tag_list"`
	RegIP            string   `json:"reg_ip" example:""`
	InvitationCode   string   `json:"invitation_code" example:""`
	Language         int      `json:"language" example:"0"`
	UserMark         string   `json:"user_mark" example:"0"`
	LastDayFollowers string   `json:"last_day_followers" example:"0"`
	MutedCount       int      `json:"muted_count" example:"1"`
	MutedStartTime   string   `json:"muted_start_time" example:"1698075740608"`
	MutedEndTime     string   `json:"muted_end_time" example:"1698162140608"`

	AvatarHD  string   `json:"avatar_hd" example:"https://img2.jacocdn.com/orj1080/b3158285la1i906ie78z1j20yi0yidj2.jpg"`
	CoverHD   []string `json:"cover_hd" example:"https://img2.jacocdn.com/orj1080/b3158285la1i906ie78z1j20yi0yidj2.jpg"`
	BadgeList []Badge  `json:"badge_list"`

	AvatarSmall string `json:"avatar_small" example:"https://img2.jacocdn.com/orj360/b3158285la1i906ie78z1j20yi0yidj2.jpg"`
	AvatarLarge string `json:"avatar_large" example:"https://img2.jacocdn.com/large/b3158285la1i906ie78z1j20yi0yidj2.jpg"`

	Following       bool `json:"following" example:"false"`
	FollowMe        bool `json:"follow_me" example:"false"`
	InBlack         bool `json:"in_black" example:"false"`
	BlackMe         bool `json:"black_me" example:"false"`
	VerifiedType    int  `json:"verified_type" example:"0"`
	Level           int  `json:"level" example:"32"`
	SpecialReminder int  `json:"special_reminder" example:"0"`
	Administrator   int  `json:"administrator_type" example:"0"`
	SubType         int  `json:"sub_type" example:"0"`

	AvatarWebp180 string `json:"avatar_webp180" example:""`
	AvatarWebp360 string `json:"avatar_webp360" example:""`
	AvatarWebp480 string `json:"avatar_webp480" example:""`
	AvatarWebp960 string `json:"avatar_webp960" example:""`

	VerifiedCode int    `json:"verified_code" example:"0"`
	VerifiedInfo string `json:"verified_info" example:""`
	UserStatus   int    `json:"user_status" example:"1"`
	AuthBit      string `json:"auth_bit" example:"0"`
	VIPInfo      []any  `json:"vip_info"`
}

type ShowTab struct {
	ShowVideoLiked    int `json:"show_video_liked" example:"0"`
	ShowVideoFavorite int `json:"show_video_favorite" example:"0"`
}

type CountInfo struct {
	FollowingCount string `json:"following_count" example:"1155"`
	FollowerCount  string `json:"follower_count" example:"32.8K"`
}

type BadgeGallery struct {
	BadgeCount int            `json:"badge_count" example:"10"`
	BadgeList  []GalleryBadge `json:"badge_list"`
}

type GalleryBadge struct {
	ID          string `json:"id" example:"5"`
	BadgeName   string `json:"badge_name" example:"Creator"`
	BadgeImage  string `json:"badge_image" example:"https://img.jacocdn.com/large/example.jpg"`
	Scheme      string `json:"scheme" example:"jaco://discover/receivedRanking?..."`
	Type        string `json:"type" example:"0"`
	SeriesType  int    `json:"series_type" example:"3"`
	BadgeNameEN string `json:"badge_name_en" example:"Creator"`
	BadgeNameAR string `json:"badge_name_ar" example:"صانع المحتوى"`
}

type UserBusiness struct {
	PowerSign string `json:"power_sign" example:"wqHqfYQv..."`
	Website   []any  `json:"website_list"`
	Age       string `json:"age" example:"-1"`
	Country   string `json:"country" example:""`
}

type UserMeta struct {
	UID               string `json:"uid" example:"3004531333"`
	MetaTitleEN       string `json:"meta_title_en" example:"Uaegirl Live: Stream Fun"`
	MetaDescriptionEN string `json:"meta_description_en" example:"Join Uaegirl's captivating streams"`
	KeywordsEN        string `json:"keywords_en" example:"live streaming content creator"`
	MetaTitleAR       string `json:"meta_title_ar" example:"بث Uaegirl بجاكو"`
	MetaDescriptionAR string `json:"meta_description_ar" example:"عش تجربة بث مباشر"`
	KeywordsAR        string `json:"keywords_ar" example:"صانع محتوى للبث المباشر"`
}

type Extra struct {
	SPLActive bool `json:"spl_active" example:"false"`
}

// ========================= comment tree ============================
type CommentTreeResponse struct {
	Error     string          `json:"error" example:"Success"`
	Data      CommentTreeData `json:"data"`
	ErrorCode int             `json:"error_code" example:"10000"`
	RequestID string          `json:"request_id" example:"82874cc5-8319-423f-8e6f-92e9c94598be"`
}

type CommentTreeData struct {
	MaxID       string    `json:"max_id" example:"0"`
	SinceID     string    `json:"since_id" example:"5191643520499712"`
	Total       int       `json:"total" example:"28"`
	Count       int       `json:"count" example:"20"`
	NextSinceID int64     `json:"next_since_id" example:"0"`
	RootCmtID   string    `json:"root_cmt_id" example:"0"`
	List        []Comment `json:"list"`
	NextMaxID   int64     `json:"next_max_id" example:"-1"`
	MaxIDType   int       `json:"max_id_type" example:"1"`
}

type Comment struct {
	CmtID            string           `json:"cmt_id" example:"5180692040351744"`
	RID              string           `json:"rid" example:"23194424305516544"`
	UID              string           `json:"uid" example:"3004631493"`
	User             CommentUser      `json:"user"`
	ReplyCmtID       string           `json:"reply_cmt_id" example:"0"`
	VFlag            string           `json:"vflag" example:"0"`
	MFlag            string           `json:"mflag" example:"0"`
	IsLiked          bool             `json:"is_liked" example:"false"`
	ChildCommentList ChildCommentList `json:"child_comment_list"`
	ParseEntityList  []any            `json:"parse_entity_list"`
	CreateAt         string           `json:"create_at" example:"1750656698573"`
	Counters         CommentCounters  `json:"counters"`
	Lat              float64          `json:"lat" example:"0"`
	Lng              float64          `json:"lng" example:"0"`
	ReplyUID         string           `json:"reply_uid" example:""`
	IsPinned         bool             `json:"is_pinned" example:"false"`
	RootCmtID        string           `json:"root_cmt_id" example:"0"`
	Text             string           `json:"text" example:"قلبييييييييي هالجمال بسم الله عليك 🥺🤍🤍🤍"`
	IsAuthorComment  bool             `json:"is_author_comment" example:"false"`
	IsAuthorLiked    bool             `json:"is_author_liked" example:"false"`
	ReelUID          string           `json:"reel_uid" example:"3004008323"`
	MentionUIDList   []string         `json:"mention_uid_list"`
	RootUID          string           `json:"root_uid" example:"0"`
}

type ChildCommentList struct {
	MaxID       string    `json:"max_id" example:"0"`
	SinceID     string    `json:"since_id" example:"0"`
	Total       int       `json:"total" example:"0"`
	Count       int       `json:"count" example:"0"`
	NextSinceID int64     `json:"next_since_id" example:"-1"`
	RootCmtID   string    `json:"root_cmt_id" example:"5180692040351744"`
	List        []Comment `json:"list"`
	NextMaxID   int64     `json:"next_max_id" example:"-1"`
	MaxIDType   int       `json:"max_id_type" example:"0"`
}

type CommentUser struct {
	ID              string   `json:"id" example:"3004631493"`
	Nickname        string   `json:"nickname" example:"أمال 🤍"`
	Username        string   `json:"username" example:"jaco_er56cd"`
	Icon            string   `json:"icon" example:"https://img2.jacocdn.com/orj1080/b31709c5la1hswjkcqlofj20n00gmjsz.jpg"`
	AvatarHD        string   `json:"avatar_hd" example:"https://img2.jacocdn.com/orj1080/b31709c5la1hswjkcqlofj20n00gmjsz.jpg"`
	AvatarSmall     string   `json:"avatar_small" example:"https://img2.jacocdn.com/orj360/b31709c5la1hswjkcqlofj20n00gmjsz.jpg"`
	FollowerCount   int      `json:"follower_count" example:"0"`
	FollowingCount  int      `json:"following_count" example:"0"`
	FriendCount     int      `json:"friend_count" example:"0"`
	TagList         []string `json:"tag_list"`
	VerifiedType    int      `json:"verified_type" example:"0"`
	BadgeList       []Badge  `json:"badge_list"`
	Level           int      `json:"level" example:"5"`
	Following       bool     `json:"following" example:"false"`
	FollowMe        bool     `json:"follow_me" example:"false"`
	Description     string   `json:"description" example:""`
	UserType        int      `json:"user_type" example:"1"`
	SpecialReminder int      `json:"special_reminder" example:"0"`
}

type Badge struct {
	Code      int    `json:"code" example:"4"`
	Name      string `json:"name" example:"Supporter Level"`
	Desc      string `json:"desc" example:""`
	GrantTime string `json:"grant_time" example:"0"`
	Type      int    `json:"type" example:"0"`
}

type CommentCounters struct {
	LikesCounter    int `json:"likes_counter" example:"4"`
	CommentsCounter int `json:"comments_counter" example:"0"`
}

// ========================= UserReelResponse ============================
type UserReelResponse struct {
	Error     string       `json:"error" example:"Success"`
	Data      UserReelData `json:"data"`
	ErrorCode int          `json:"error_code" example:"10000"`
	RequestID string       `json:"request_id" example:"01f1ee42-a913-4ad2-832e-95e6fb6c8dec"`
}

type UserReelData struct {
	Category             int               `json:"category"`
	ImmersionExperiment  int               `json:"immersion_experiment"`
	ImmersionProperty    ImmersionProperty `json:"immersion_property"`
	ImmersionType        int               `json:"immersion_type"`
	JustWatchSwitch      int               `json:"just_watch_switch"`
	LastVideoPlayCounter int               `json:"last_video_play_counter"`
	List                 []ListItem        `json:"list"`
}

type ImmersionProperty struct {
	Extend       string `json:"extend"`
	ImmersionURL string `json:"immersion_url"`
	RefreshURL   string `json:"refresh_url"`
}

type ListItem struct {
	Content    Content `json:"content"`
	ID         string  `json:"id"`
	ReportData string  `json:"report_data"`
	Type       int     `json:"type"`
	UID        string  `json:"uid"`
}

type Content struct {
	ShareLink string    `json:"share_link"`
	VideoInfo VideoInfo `json:"video_info"`
}

type VideoInfo struct {
	Activity       Activity          `json:"activity"`
	CFlag          string            `json:"cflag"`
	CoCreationAble bool              `json:"co_creation_able"`
	Counters       Counters          `json:"counters"`
	CoverImage     CoverImage        `json:"cover_image"`
	CreateAt       string            `json:"create_at"`
	IsFavorite     bool              `json:"is_favorite"`
	IsLiked        bool              `json:"is_liked"`
	Location       Location          `json:"location"`
	Media          Media             `json:"media"`
	MediaType      int               `json:"media_type"`
	MusicFirstPost bool              `json:"music_first_post"`
	ObjectIndex    []ObjectIndex     `json:"object_index"`
	Objects        map[string]Object `json:"objects"`
	Properties     Properties        `json:"properties"`
	RankObject     RankObject        `json:"rankObject"`
	RemixCodec     RemixCodec        `json:"remix_codec"`
	RemixRID       string            `json:"remix_rid"`
	ReviewState    int               `json:"review_state"`
	ReviewType     int               `json:"review_type"`
	RID            string            `json:"rid"`
	Settings       Settings          `json:"settings"`
	State          string            `json:"state"`
	Text           string            `json:"text"`
	Transcribe     Transcribe        `json:"transcribe"`
	UID            string            `json:"uid"`
	VFlag          string            `json:"vflag"`
	VideoType      int               `json:"video_type"`
}

type Activity struct {
	FailDesc string `json:"fail_desc"`
	ID       string `json:"id"`
	Name     string `json:"name"`
	Status   int    `json:"status"`
	Type     int    `json:"type"`
}

type Counters struct {
	CommentsCounter  int `json:"comments_counter"`
	FavoriteCounter  int `json:"favorite_counter"`
	GoldCounter      int `json:"gold_counter"`
	LikesCounter     int `json:"likes_counter"`
	ShareCounter     int `json:"share_counter"`
	VideoPlayCounter int `json:"video_play_counter"`
}

type CoverImage struct {
	Height      string            `json:"height"`
	PID         string            `json:"pid"`
	Src         int               `json:"src"`
	URLTemplate string            `json:"url_template"`
	URLs        map[string]string `json:"urls"`
	Width       string            `json:"width"`
}

type Location struct {
	City string `json:"city"`
}

type Media struct {
	PlayInfo PlayInfo `json:"play_info"`
}

type PlayInfo struct {
	DefaultCodec   string        `json:"default_codec"`
	DefaultQuality int           `json:"default_quality"`
	Details        []MediaDetail `json:"details"`
}

type MediaDetail struct {
	Duration           float64 `json:"duration"`
	FirstFrameCoverURL string  `json:"first_frame_cover_url"`
	Height             int     `json:"height"`
	PostRollDuration   float64 `json:"post_roll_duration"`
	PostRollURL        string  `json:"post_roll_url"`
	QualityDesc        string  `json:"quality_desc"`
	QualityIndex       int     `json:"quality_index"`
	Size               int     `json:"size"`
	URL                string  `json:"url"`
	VideoCodec         string  `json:"video_codec"`
	Width              int     `json:"width"`
}

type ObjectIndex struct {
	ID   string `json:"id"`
	Type string `json:"type"`
}

type Object struct {
	Type    string         `json:"@type"`
	Content MusicViewModel `json:"content"`
}

type MusicViewModel struct {
	Type         string       `json:"@type"`
	AuthorInfo   AuthorInfo   `json:"author_info"`
	CoverImage   CoverImage   `json:"cover_image"`
	CreateTime   string       `json:"create_time"`
	Features     Features     `json:"features"`
	MID          string       `json:"mid"`
	MusicControl MusicControl `json:"music_control"`
	MusicInfo    MusicInfo    `json:"music_info"`
	State        int          `json:"state"`
	OffsetX      int          `json:"offset_x"`
	OffsetY      int          `json:"offset_y"`
}

type AuthorInfo struct {
	AvatarSmall string `json:"avatar_small"`
	Icon        string `json:"icon"`
	ID          string `json:"id"`
	Nickname    string `json:"nickname"`
	Username    string `json:"username"`
}

type Features struct {
	IsOfficial bool `json:"is_official"`
	IsUGC      bool `json:"is_ugc"`
}

type MusicControl struct {
	IsUse           bool `json:"is_use"`
	MuteChangeName  bool `json:"mute_change_name"`
	MuteShare       bool `json:"mute_share"`
	PreventDownload bool `json:"prevent_download"`
}

type MusicInfo struct {
	AIFlag       int    `json:"ai_flag"`
	AlbumName    string `json:"album_name"`
	AlbumNameAR  string `json:"album_name_ar"`
	AuthorName   string `json:"author_name"`
	AuthorNameAR string `json:"author_name_ar"`
	AuthorUID    string `json:"author_uid"`
	LanguageType int    `json:"language_type"`
	MusicDesc    string `json:"music_desc"`
	MusicDescAR  string `json:"music_desc_ar"`
	MusicName    string `json:"music_name"`
	MusicNameAR  string `json:"music_name_ar"`
}

type Properties struct {
	IsNew bool `json:"is_new"`
	Pin   bool `json:"pin"`
}

type RankObject struct {
	TopRank []interface{} `json:"topRank"`
	Total   int           `json:"total"`
}

type RemixCodec struct {
	RemixDefaultCodec   string `json:"remix_default_codec"`
	RemixDefaultQuality int    `json:"remix_default_quality"`
}

type Settings struct {
	CoCreationAble   bool `json:"co_creation_able"`
	DownloadAble     bool `json:"download_able"`
	ExtractSoundAble bool `json:"extract_sound_able"`
	FullScreenAble   bool `json:"full_screen_able"`
	PlaySound        bool `json:"play_sound"`
	QuoteSoundAble   bool `json:"quote_sound_able"`
	RewardAble       bool `json:"reward_able"`
}

type Transcribe struct {
	LID string `json:"lid"`
	MID string `json:"mid"`
	UID string `json:"uid"`
}
