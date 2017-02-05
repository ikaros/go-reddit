package reddit

import "encoding/json"

// Thing is the reddit API's base class
type Thing struct {
	// This item's identifier, e.g. "8xwlg"
	ID string

	// Fullname of comment, e.g. "t1_c3v7f8u"
	Name string

	// All `thing`s have a `kind`.  The kind is a String identifier
	// that denotes the object's type.
	// Some examples: `Listing`, `more`, `t1`, `t2`
	Kind string

	// Data is a custom data structure used to hold valuable information.
	// This object's format will follow the data structure respective
	// of its kind.
	//Data json.RawMessage
}

// Listing is used to paginate content that is too long to display in one go.
// Add the query argument `before` or `after` with the value given to get
// the previous or next page. This is usually used in conjunction
// with a `count` argument.
//
// NOTE: Unlike the other classes documented on this page, Listing is not a
// thing subclass, as it inherits directly from the Python base class, Object.
// Although it does have data and kind as parameters, it does not have
// id and name. A listing's kind will always be Listing and its data
// will be a List of things.
type Listing struct {

	// The fullname of the listing that follows before this page.
	// `null` if there is no previous page.
	Before *string `json:"before"`

	// The fullname of the listing that follows after this page.
	// `null` if there is no next page.
	After *string `json:"after"`

	// This modhash is not the same modhash provided upon login.
	// You do not need to update your user's modhash every time you
	// get a new modhash.  You can reuse the modhash given upon login.
	ModHash string `json:"modhash"`

	// A list of `thing`s that this Listing wraps
	Children []ListingThing `json:"children"`
}

type ListingThing struct {
	// All `thing`s have a `kind`.  The kind is a String identifier
	// that denotes the object's type.
	// Some examples: `Listing`, `more`, `t1`, `t2`
	Kind string

	Data json.RawMessage
}

type Votable struct {
	// The number of upvotes. (includes own)
	Ups int `json:"ups"`

	// The number of downvotes. (includes own)
	Downs int `json:"downs"`

	// `true` if thing is liked by the user, `false` if thing is disliked,
	// `null` if the user has not voted or you are not logged in.
	Likes *bool `json:"likes"`
}

type Created struct {
	// The time of creation in local epoch-second format. ex: `1331042771.0`
	Created float64 `json:"created"`

	// The time of creation in UTC epoch-second format.
	// Note that neither of these ever have a non-zero fraction.
	CreatedUTC float64 `json:"created_utc"`
}

type Comment struct {
	Created
	Votable

	// Who approved this comment. null if nobody or you are not a mod
	ApprovedBy *string `json:"approved_by"`

	// The account name of the poster
	Author string `json:"author"`

	// The CSS class of the author's flair. subreddit specific
	AuthorFlairCSSClass string `json:"author_flair_css_class"`

	// The text of the author's flair. subreddit specific
	AuthorFlairText string `json:"author_flair_text"`

	// Who removed this comment. null if nobody or you are not a mod
	BannedBy *string `json:"banned_by"`

	// The raw text. This is the unformatted text which includes the raw
	// markup characters such as ** for bold. <, >, and & are escaped.
	Body string `json:"body"`

	// The formatted HTML text as displayed on reddit.
	// For example, text that is emphasised by * will now have <em>
	// tags wrapping it. Additionally, bullets and numbered lists
	// will now be in HTML list format.
	// NOTE: The HTML string will be escaped. You must unescape
	// to get the raw HTML.
	BodyHTML string `json:"body_html"`

	// False if not edited, edit date in UTC epoch-seconds otherwise.
	// NOTE: for some old edited comments on reddit.com,
	// this will be set to true instead of edit date.
	Edited string `json:"edited"`

	// The number of times this comment received reddit gold
	Gilded int `json:"gilded"`

	// How the logged-in user has voted on the comment - True = upvoted,
	// False = downvoted, null = no vote
	Likes *bool `json:"likes"`

	// Present if the comment is being displayed outside its thread
	// (user pages, /r/subreddit/comments/.json, etc.).
	// Contains the author of the parent link
	LinkAuthor *string `json:"link_author"`

	// ID of the link this comment is in
	LinkID string `json:"link_id"`

	// Present if the comment is being displayed outside its thread
	// (user pages, /r/subreddit/comments/.json, etc.).
	// Contains the title of the parent link
	LinkTitle *string `json:"link_title"`

	// Present if the comment is being displayed outside its thread
	// (user pages, /r/subreddit/comments/.json, etc.).
	// Contains the URL of the parent link
	LinkURL string `json:"link_url"`

	// How many times this comment has been reported,
	// null if not a mod
	NumReports *int `json:"num_reports"`

	// ID of the thing this comment is a reply to,
	// either the link or a comment in it
	ParentID string `json:"parent_id"`

	// A list of replies to this comment
	Replies []ListingThing `json:"replies"`

	// True if this post is saved by the logged in user
	Saved bool `json:"saved"`

	// The net-score of the comment
	Score int `json:"score"`

	// Whether the comment's score is currently hidden.
	ScoreHidden bool `json:"score_hidden"`

	// Subreddit of thing excluding the /r/ prefix. "pics"
	Subreddit string `json:"subreddit"`

	// The id of the subreddit in which the thing is located
	SubredditID string `json:"subreddit_id"`

	// To allow determining whether they have been distinguished
	// by moderators/admins.
	// null = not distinguished.
	// moderator = the green [M].
	// admin = the red [A].
	// special = various other special distinguishes
	Distinguished *string `json:"distinguished"`
}

type Link struct {
	Votable
	Created

	// The account name of the poster. null if this is
	// a promotional link
	Author string `json:"author"`

	// The CSS class of the author's flair. subreddit specific
	AuthorFlairCSSClass string `json:"author_flair_css_class"`

	// The text of the author's flair. subreddit specific
	AuthorFlairText string `json:"author_flair_text"`

	// probably always returns false
	Clicked bool `json:"clicked"`

	// The domain of this link. Self posts will be self.<subreddit>
	// while other examples include en.wikipedia.org and s3.amazon.com
	Domain string `json:"domain"`

	// true if the post is hidden by the logged in user.
	// false if not logged in or not hidden.
	Hidden bool `json:"hidden"`

	// True if this link is a selfpost
	IsSelf bool `json:"is_self"`

	// How the logged-in user has voted on the link - True = upvoted,
	// False = downvoted, null = no vote
	Likes bool `json:"likes"`

	// The CSS class of the link's flair.
	LinkFlairCSSClass string `json:"link_flair_css_class"`

	// The text of the link's flair.
	LinkFlairText string `json:"link_flair_text"`

	// Whether the link is locked (closed to new comments) or not.
	Locked bool `json:"locked"`

	// Used for streaming video. Detailed information about
	// the video and it's origins are placed here
	Media json.RawMessage `json:"media"`

	// Used for streaming video. Technical embed
	// specific information is found here.
	MediaEmbed json.RawMessage `json:"media_embed"`

	// The number of comments that belong to this link.
	// includes removed comments.
	NumComments int `json:"num_comments"`

	// True if the post is tagged as NSFW. False if otherwise
	Over18 bool `json:"over18"`

	// Relative URL of the permanent link for this link
	Permalink string `json:"permalink"`

	// True if this post is saved by the logged in user
	Saved bool `json:"saved"`

	// The net-score of the link.
	//
	// NOTE: A submission's score is simply the number of upvotes
	// minus the number of downvotes. If five users like the
	// submission and three users don't it will have a score of 2.
	// Please note that the vote numbers are not "real" numbers,
	// they have been "fuzzed" to prevent spam bots etc.
	// So taking the above example, if five users upvoted the
	// submission, and three users downvote it, the upvote/downvote
	// numbers may say 23 upvotes and 21 downvotes, or 12 upvotes,
	// and 10 downvotes. The points score is correct,
	// but the vote totals are "fuzzed".
	Score int `json:"score"`

	// The raw text. this is the unformatted text which includes the
	// raw markup characters such as ** for bold. <, >, and & are escaped.
	// Empty if not present.
	Selftext string `json:"selftext"`

	// The formatted escaped HTML text. this is the HTML formatted version
	// of the marked up text. Items that are boldened by ** or *** will now
	// have <em> or *** tags on them. Additionally, bullets and numbered
	// lists will now be in HTML list format.
	//
	// NOTE: The HTML string will be escaped. You must unescape to get
	// the raw HTML.
	//
	// Null if not present.
	SelftextHTML *string `json:"selftext_html"`

	// Subreddit of thing excluding the /r/ prefix. "pics"
	Subreddit string `json:"subreddit"`

	// The id of the subreddit in which the thing is located
	SubredditID string `json:"subreddit_id"`

	// Full URL to the thumbnail for this link;
	// "self" if this is a // self post;
	// "default" if a thumbnail is not available
	Thumbnail string `json:"thumbnail"`

	// The title of the link. may contain newlines for some reason
	Title string `json:"title"`

	// The link of this post. the permalink if this is a self-post
	URL string `json:"url"`

	// Indicates if link has been edited. Will be the edit timestamp
	// if the link has been edited and return false otherwise.
	Edited json.RawMessage `json:"edited"`

	// To allow determining whether they have been distinguished
	// by moderators/admins.
	// null = not distinguished.
	// moderator = the green [M].
	// admin = the red [A].
	// special = various other special distinguishes
	Distinguished string `json:"distinguished"`

	// True if the post is set as the sticky in its subreddit.
	Stickied bool `json:"stickied"`
}

type Subreddit struct {

	// Number of users active in last 15 minutes
	AccountsActive int `json:"accounts_active"`

	// Number of minutes the subreddit initially
	// hides comment scores.
	CommentScoreHideMins int `json:"comment_score_hide_mins"`

	// Sidebar text.
	Description string `json:"description"`

	// Sidebar text, escaped HTML format.
	DescriptionHTML string `json:"description_html"`

	// Human name of the subreddit
	DisplayName string `json:"display_name"`

	// Full URL to the header image, or null
	HeaderIMG *string `json:"header_img"`

	// Width and height of the header image as [widt,height], or null
	HeaderSize *[2]int `json:"header_size"`

	// Description of header image shown on hover, or null
	HeaderTitle *string `json:"header_title"`

	// Whether the subreddit is marked as NSFW
	Over18 bool `json:"over18"`

	// Description shown in subreddit search results.
	PublicDescription string `json:"public_description"`

	// Whether the subreddit's traffic page is publicly-accessible
	PublicTraffic bool `json:"public_traffic"`

	// The number of redditors subscribed to this subreddit
	Subscribers int64 `json:"subscribers"`

	// The type of submissions the subreddit allows.
	// One of "any", "link" or "self".
	SubmissionType string `json:"submission_type"`

	// The subreddit's custom label for the submit // link button, if any
	SubmitLinkLabel string `json:"submit_link_label"`

	// The subreddit's custom label for the submit text button, if any
	SubmitTextLabel string `json:"submit_text_label"`

	// The subreddit's type - one of "public", "private", "restricted",
	// or in very special cases "gold_restricted" or "archived"
	SubredditType string `json:"subreddit_type"`

	// Title of the main page
	Title string `json:"title"`

	// The relative URL of the subreddit. Ex: "/r/pics/"
	URL string `json:"url"`

	// Whether the logged-in user is banned from the subreddit.
	UserIsBanned bool `json:"user_is_banned"`

	// Whether the logged-in user is an approved submitter in the subreddit
	UserIsContributor bool `json:"user_is_contributor"`

	// Whether the logged-in user is a moderator of the subreddit
	UserIsModerator bool `json:"user_is_moderator"`

	// Whether the logged-in user is subscribed to the subreddit
	UserIsSubscriber bool `json:"user_is_subscriber"`
}

type Message struct {
	Created

	Author string `json:"author"`

	// The message itself.
	Body string `json:"body"`

	// The message itself with HTML formatting.
	BodyHTML string `json:"body_html"`

	// If the message is a comment, then the permalink to the comment
	// with ?context=3 appended to the end, otherwise an empty string
	Context string `json:"context"`

	// Either null or the first message's ID represented as base 10 (wtf)
	FirstMessage *Message `json:"first_message"`

	// Either null or the first message's fullname
	FirstMessageName *string `json:"first_message_name"`

	// How the logged-in user has voted on the message.
	// True = upvoted,
	// False = downvoted,
	// null = no vote
	Likes *bool `json:"likes"`

	// If the message is actually a comment, contains the title
	// of the thread it was posted in
	LinkTitle string `json:"link_title"`

	// ex: "t4_8xwlg"
	Name string `json:"name"`

	// unread? not sure.
	New bool `json:"new"`

	// Null if no parent is attached.
	ParentID *string `json:"parent_id"`

	// An empty string if there are no replies.
	Replies string `json:"replies"`

	// Subject of message.
	Subject string `json:"subject"`

	// null if not a comment.
	Subreddit string `json:"subreddit"`

	WasComment bool `json:"was_comment"`
}

// Example of raw account data:
// 	{
// 		"kind": "t2",
// 		"data": {
// 			"has_mail": false,
// 			"name": "fooBar",
// 			"created": 123456789.0,
// 			"modhash": "f0f0f0f0f0f0f0f0...",
// 			"created_utc": 1315269998.0,
// 			"link_karma": 31,
// 			"comment_karma": 557,
// 			"is_gold": false,
// 			"is_mod": false,
// 			"has_verified_email": false,
// 			"id": "5sryd",
// 			"has_mod_mail": false
// 		}
// 	}
type Account struct {
	Created

	// User's comment karma
	CommentKarma int `json:"comment_karma"`

	// User has unread mail? null if not your account
	HasMail *bool `json:"has_mail"`

	// User has unread mod mail? null if not your account
	HasModMail *bool `json:"has_mod_mail"`

	// User has provided an email address and got it verified.
	HasVerifiedEmail bool `json:"has_verified_email"`

	// ID of the account; prepend t2_ to get fullname
	ID string `json:"id"`

	// Number of unread messages in the inbox.
	// Not present if not your account
	InboxCount *int `json:"inbox_count"`

	// Whether the logged-in user has this user set as a friend.
	IsFriend bool `json:"is_friend"`

	// Reddit gold status
	IsGold bool `json:"is_gold"`

	// Whether this account moderates any subreddits
	IsMod bool `json:"is_mod"`

	// User's link karma
	LinkKarma int `json:"link_karma"`

	// Current modhash. not present if not your account
	Modhash *string `json:"modhash"`

	// The username of the account in question. This attribute
	// overrides the superclass's name attribute. Do not confuse
	// an account's name which is the account's username
	// with a thing's name which is the thing's FULLNAME.
	Name string `json:"name"`

	// Whether this account is set to be over 18
	Over18 bool `json:"over18"`
}

// Example of more:
// 	{
// 		"kind": "more",
// 		"data": {
// 			"children": [
// 				"c3y9tyh",
// 				"c3y9ul7",
// 				"c3y9unp",
// 				"c3y9uoi"
// 			],
// 			"name": "t1_c3y9tyh",
// 			"id": "c3y9tyh"
// 		}
// 	}
type More struct {

	// A list of String `id`s that are the additional `thing`s that can
	// be downloaded but are not because there are too many to list.
	Children []string `json:"children"`
}

type kind string

const (
	kindComment       kind = "t1"
	kindAccount       kind = "t2"
	kindLink          kind = "t3"
	kindMessage       kind = "t4"
	kindSubreddit     kind = "t5"
	kindAward         kind = "t6"
	kindPromoCampaign kind = "t8"
)
