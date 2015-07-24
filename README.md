# FlickApp_Go
This Server interacts with the Flickr Apis and with the Mongo Database.

On using the flickr api, we receive data about various images matching our search criteria. A part of that data is stored in the database and then sent to the node server as response.

For Upvote and Downvote , the image_id is used to identify the image and update its upvote/downvote count.

The api used is : https://api.flickr.com/services/rest/?method=flickr.photos.search&api_key=408220e4ee5ccf0f5802e1cd9cd53fea&extras=url_l&per_page=30&format=json&nojsoncallback=1&text= , with the text changed according to user input.

The database is on the MongoLab cloud.

Go server needs to be up and running before we use the node server.
