package models

type GithubEmail struct {
	Email      string `json:"email"`
	Verified   bool   `json:"verified"`
	Primary    bool   `json:"primary"`
	Visibility string `json:"visibility"`
}

type GithubEmails []GithubEmail

type GithubUser struct {
	Login string `json:"login"`
	ID int `json:"id"`
	NodeId string `json:"node_id"`
	AvatarUrl string `json:"avatar_url"`
	GravatarId string `json:"gravatar_id"`
	Url string `json:"url"`
	HtmlUrl string `json:"html_url"`
	FollowerUrl string `json:"follower_url"`
	Type string `json:"type"`
	SideAdmin bool `json:"side_admin"`
	Name string `json:"name"`
	Company string `json:"company"`
	Blog 	string `json:"blog"`
	Location string `json:"location"`
	Email string `json:"email"`
	EmailVerified bool `json:"email_verified"`
	Hireable bool `json:"hireable"`
	Bio string `json:"bio"`
	TwitterUsername string `json:"twitter_username"`
	CreatedAt	string `json:"created_at"`
	UpdatedAt 	string `json:"updated_at"`
}