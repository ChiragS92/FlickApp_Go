package main

import(
	"fmt"
	"net/http"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"os"
)

type Image struct {
	Id bson.ObjectId 'bson:"_id"'
	Title string 'bson:"Title"'
	Url string 'bson:"URL"'
	Upvotes int 'bson:"Upvotes"'
	Downvotes int 'bson:"Downvotes"'	
}

func main(){
	http.HandleFunc('/getImages',getImages)
	http.HandleFunc('/upvote',updateUpvote)
	http.HandleFunc('/downvote',updateDownvote)
	http.ListenAndServe(':7001',nil);
}

func getImages()
{
	var connection_url = "mongodb://chirag:flickr@ds036648.mongolab.com:36648/flickrapp"
	session,err = mgo.Dial(connection_url);
	if err!=null {
		fmt.printf("Cannot connect to database!!");
		os.Exit(1);
	}
	defer.session.Close();

}