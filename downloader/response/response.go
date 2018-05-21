package response

import "encoding/xml"

type Tumblr struct {
	XMLName   xml.Name `xml:"tumblr"` // 指定最外层的标签为 tumblr
	Tumblelog string   `xml:"tumblelog"`
}

type BasePosts struct {
	Type  string `xml:"type,attr"`
	Start string `xml:"start,attr"`
	Total string `xml:"total,attr"`
}

type BasePost struct {
	Id        string    `xml:"id,attr"`
	ReblogKey string    `xml:"reblog-key,attr"`
	Type      string    `xml:"type,attr"`
	Tag       []string  `xml:"tag"`
	Tumblelog Tumblelog `xml:"tumblelog"`
	PhotoCaption string `xml:"photo-caption"`
}

type Tumblelog struct {
	Title string `xml:"title,attr"`
	Name  string `xml:"name,attr"`
	Value string `xml:",chardata"`
}
