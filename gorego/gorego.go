package gorego

import (
	"log"
	"github.com/mrjones/oauth"
	"encoding/xml"
	"io/ioutil"
	"net/http"
	"strconv"
)

var (
	apiRoot = "https://www.goodreads.com/"
	key = "QJ6PQVdkuDiOi8DMGVrFg"
	secret = "BOpRDZvEDt8maSNlIsgZS9cUY589a8m1W6z3TKM4"
	goodreads_consumer = oauth.NewConsumer(
	key,
	secret,
	oauth.ServiceProvider{
		RequestTokenUrl:   "http://www.goodreads.com/oauth/request_token",
		AuthorizeTokenUrl: "http://www.goodreads.com/oauth/authorize",
		AccessTokenUrl:    "http://www.goodreads.com/oauth/access_token",
	})
)

func GetAuthorisationURL() (*oauth.RequestToken, string) {
	_requestToken, auth_url, err := goodreads_consumer.GetRequestTokenAndUrl("")
	if err != nil {
		log.Fatal(err)
	}
	return _requestToken, auth_url
}

func GetAccessToken(verificationCode string, requestToken *oauth.RequestToken) (*oauth.AccessToken) {
	accessToken, err := goodreads_consumer.AuthorizeToken(requestToken, verificationCode)
	if err != nil {
		log.Fatal(err)
	}
	return accessToken
}

func QueryUser(accessToken *oauth.AccessToken) (GRUser) {
	client, err := goodreads_consumer.MakeHttpClient(accessToken)
	if err != nil {
		log.Fatal(err)
	}

	response, err := client.Get(
		"https://www.goodreads.com/api/auth_user")
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()

	bits, err := ioutil.ReadAll(response.Body)
	log.Println(string(bits))
	user := GRUserResponse{}
	xml.Unmarshal(bits, &user)
	log.Println(user)
	return user.User
}

func GetReviewList(userId string, limit int) *Response {
	l := strconv.Itoa(limit)
	uri := apiRoot + "review/list?v=2" + "&id=" + userId + "&key=" + key + "&shelf=read&sort=date_read&order=d&per_page=" + l
	log.Println("Hitting " + uri)
	response := &Response{}
	getData(uri, response)

	return response
}

func getRequest(uri string) []byte {
	res, err := http.Get(uri)
	if err != nil {
		log.Fatal(err)
	}

	body, err := ioutil.ReadAll(res.Body)
	log.Println(string(body))
	res.Body.Close()
	if err != nil {
		log.Fatal(err)
	}

	return body
}


func getData(uri string, i interface{}) {
	data := getRequest(uri)
	xmlUnmarshal(data, i)
}


func xmlUnmarshal(b []byte, i interface{}) {
	err := xml.Unmarshal(b, i)
	if err != nil {
		log.Fatal(err)
	}
}

type GRUserResponse struct {
	User GRUser `xml:"user"`
}
type GRUser struct {
	Id string `xml:"id,attr"`
	Name string `xml:"name"`
}

type GRReview struct {
	GRBook   GRBook   `xml:"book"`
	Rating int    `xml:"rating"`
	ReadAt string `xml:"read_at"`
	Link   string `xml:"link"`
}

type GRBook struct {
	ID       string   `xml:"id"`
	Title    string   `xml:"title"`
	Link     string   `xml:"link"`
	ImageURL string   `xml:"image_url"`
	AvgRating string   `xml:"average_rating"`
	NumPages string   `xml:"num_pages"`
	Format   string   `xml:"format"`
	GRAuthors  []GRAuthor `xml:"authors>author"`
	ISBN     string   `xml:"isbn"`
}

type GRAuthor struct {
	ID   string `xml:"id"`
	Name string `xml:"name"`
	Link string `xml:"link"`
}

type Response struct {
        GRBook    GRBook     `xml:"book"`
        GRReviews []GRReview `xml:"reviews>review"`
}
