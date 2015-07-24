package main

import(
	"fmt"
	"net/http"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"os"
	"io/ioutil"
	"encoding/json"
	//"strings"
	//"strconv"
)

type Image struct {
    Image_ID string 	`bson:"Image_id"` 	
	Title string 		`bson:"Title"`
	Text string   		`bson:"Text"`
	Url string 			`bson:"URL"`
	Upvotes int 		`bson:"Upvotes"`
	Downvotes int 		`bson:"Downvotes"`	
}

type Images struct {
		Photos struct {
		Page int 				`json: "page"`
		Pages int 				`json: "pages"`
		PerPage int 			`json: "perpage"`
		Total string 			`json: "total"`
		Photo []struct {
				Id string 		`json: "id"`
				Owner string 	`json: "owner"`
				Secret string 	`json: "secret"`
				Server string 	`json: "server"`
				Farm int 		`json: "farm"`
				Title string 	`json: "title"`
				IsPublic int 	`json: "ispublic"`
				IsFriend int 	`json: "isfriend"`
				IsFamily int 	`json: "isfamily"`
				Url_l string	`json: "url_l"`
				height_l int	`json: "height_l"`
				width_l int		`json: "width_l"`
				} 				`json: "photo"`
		} 						`json: "photos"`
	Status string 				`json: "stat"`
}


func main(){
	http.HandleFunc(`/getImages`,getImages)
	http.HandleFunc(`/upvote`,updateUpvote)
	http.HandleFunc(`/downvote`,updateDownvote)
	http.ListenAndServe(`:7001`,nil)
}


func updateDownvote(w http.ResponseWriter, r *http.Request) {
	var downvotes = r.URL.Query().Get(`downvotes`);
	var id = r.URL.Query().Get(`id`);
	fmt.Println(id)
	fmt.Println(downvotes)
	var connection_url = "mongodb://chirag:flickr@ds036648.mongolab.com:36648/flickrapp"
	session,err := mgo.Dial(connection_url)
	if err!=nil {
		fmt.Printf("Cannot connect to database!!")
		os.Exit(1)
	}
	defer session.Close()

	session.SetMode(mgo.Strong , true)
	collection := session.DB("flickrapp").C("images")

	err = collection.Update(bson.M{"Image_id": id}, bson.M{ "$set": bson.M{ "Downvotes": downvotes}});
	if err != nil {
		fmt.Printf("Could not update!")
		os.Exit(5)
	}else {
		fmt.Println(w,"Success!");
	}
}

func updateUpvote(w http.ResponseWriter, r *http.Request) {
	fmt.Println("In updateUpvote!")
	var upvotes = r.URL.Query().Get(`upvotes`)
	var id = r.URL.Query().Get(`id`)
	fmt.Println(id)
	fmt.Println(upvotes)
	var connection_url = "mongodb://chirag:flickr@ds036648.mongolab.com:36648/flickrapp"
	session,err := mgo.Dial(connection_url)
	if err!=nil {
		fmt.Printf("Cannot connect to database!!")
		os.Exit(1)
	}
	defer session.Close()

	session.SetMode(mgo.Strong, true)
	collection := session.DB("flickrapp").C("images")

	err = collection.Update(bson.M{"Image_id": id}, bson.M{ "$set": bson.M{ "Upvotes": upvotes}});
	if err != nil {
		fmt.Println(err)
		fmt.Printf("Could not update!")
		os.Exit(5)
	}else {
		fmt.Println(w,"Success!")
	}
}

func populateDb(text string)  {
	fmt.Println("In populateDb!")
	var connection_url = "mongodb://chirag:flickr@ds036648.mongolab.com:36648/flickrapp"
	session,err := mgo.Dial(connection_url)
	if err!=nil {
		fmt.Printf("Cannot connect to database!!")
		os.Exit(1)
	}
	defer session.Close()

	session.SetMode(mgo.Strong, true)
	collection := session.DB("flickrapp").C("images")

	var api_url = "https://api.flickr.com/services/rest/?method=flickr.photos.search&api_key=408220e4ee5ccf0f5802e1cd9cd53fea&extras=url_l&per_page=30&format=json&nojsoncallback=1&text="+text
	fmt.Println(api_url)
	res,err := http.Get(api_url)

	fmt.Println(res)
	
	if err != nil {
		fmt.Printf("Error fetching images from Flickr!")
		os.Exit(1)
	}

	defer res.Body.Close()

	body,err := ioutil.ReadAll(res.Body)

	if err != nil {
		fmt.Printf("Error reading response!")
		os.Exit(1)
	}

	var images Images

	err = json.Unmarshal(body,&images)

	if err != nil {
		fmt.Printf("Error Unmarshalling response!")
		os.Exit(1)
	}
	for value := range images.Photos.Photo {
		err := collection.Insert(&Image{Image_ID : images.Photos.Photo[value].Id,Title : images.Photos.Photo[value].Title, Text : text, Url : images.Photos.Photo[value].Url_l, Upvotes: 0, Downvotes: 0})
        if err != nil {
		fmt.Printf("%s", err)
		os.Exit(5)
		}
	}
}

func getImages(w http.ResponseWriter, r *http.Request){
	fmt.Println("Inside getImages!")
	var text = r.URL.Query().Get(`text`);
	var connection_url = "mongodb://chirag:flickr@ds036648.mongolab.com:36648/flickrapp"
	session,err := mgo.Dial(connection_url)
	if err!=nil {
		fmt.Printf("Cannot connect to database!!")
		os.Exit(1)
	}
	defer session.Close()

	session.SetMode(mgo.Strong, true)
	collection := session.DB("flickrapp").C("images")
	var images []Image
	collection.Find(bson.M{"Text" : text}).All(&images);
	count , err := collection.Count()
	if err!=nil{
		fmt.Printf("Cannot connect to database!!")
		os.Exit(1)
	}

	if count==0	|| images == nil {
		populateDb(text)
	}

	response , err := getJson(text)

	if err!=nil {
		fmt.Printf("Cannot connect to database!!")
		os.Exit(1)
	}

	fmt.Fprintf(w,string(response))

}

func getJson(text string)([]byte,error){
	var connection_url = "mongodb://chirag:flickr@ds036648.mongolab.com:36648/flickrapp"
	session,err := mgo.Dial(connection_url)
	if err!=nil {
		fmt.Printf("Cannot connect to database!!")
		os.Exit(1)
	}
	var images []Image
	err = session.DB("flickrapp").C("images").Find(bson.M{"Text" : text}).All(&images)
	if err!=nil{
		fmt.Printf("No data available!");
		os.Exit(1);
	}
	return json.MarshalIndent(images,"","")
}