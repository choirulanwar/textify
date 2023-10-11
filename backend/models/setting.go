package models

type Setting struct {
	ID               uint   `gorm:"primarykey" json:"id"`
	BrowserPath      string `gorm:"column:browser_path;type:text" json:"browser_path"`
	BrowserVisible   bool   `gorm:"column:browser_visible;type:bool" json:"browser_visible"`
	SessionGoogle    string `gorm:"column:session_google;type:text" json:"session_google"`
	SessionPinterest string `gorm:"column:session_pinterest;type:text" json:"session_pinterest"`
	ProxyUrl         string `gorm:"column:proxy_url;type:text" json:"proxy_url"`
}
