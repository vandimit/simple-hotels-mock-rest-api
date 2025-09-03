# Simple hotels mock rest api

A very simple local test rest-api, ideal to make quick front end examples.

It contains detailed information about ten hotels located in seattle (if you want to add
more entries to this local api please check the _Contributors welcome_ section).

# Why this dev local api?

Hitting real world api's is great, but when a training is on going, it can happen that 
the internet connectivity has some firewall restrictions (or even worse no internet connection available), you need an auth api key per
student, or even worse that the api has been deprecated and  your examples don't work because of that.

The goal of this project is to provide a dumb mock-data api to let trainers / students focus
on the content of their training (e.g. Front End development), and reduce the external noise going around.

If you are looking for other dumb rest-api's available online, check out:

- JSON Place holder: https://jsonplaceholder.typicode.com/
- Star wars api: https://swapi.co/

# How it works

## Getting started

### Option 1: Go

If you want to run the application directly with Go:

1. Make sure you have Go installed (version 1.19 or later)
2. Clone this repository
3. Run the application:

```
go run main.go
```

### Option 2: Docker

If you prefer to use Docker:

1. Make sure you have Docker installed
2. Clone this repository
3. Build and run using Docker Compose:

```
docker-compose up -d
```

Or you can build and run the Docker container manually:

```
# Build the Docker image
docker build -t hotels-mock-api .

# Run the container
docker run -p 8080:8080 hotels-mock-api
```

## Querying the API

Just open your web browser, postman or favourite tool, and type the following url's

- Get the list of hotels available:

```
http://localhost:8080/api/hotels
```

- To get the details of a given hotel (last url chunk is the id of the hotel):

```
http://localhost:8080/api/hotels/0248058a-27e4-11e6-ace6-a9876eff01b3
```

- Get hotels with filters (supports filtering by name, city, countryCode, minRate, maxRate, minRating, maxRating, amenityMask):

```
http://localhost:8080/api/hotels?city=Seattle&minRating=3
```

- Get hotels with pagination (supports limit and offset parameters):

```
http://localhost:8080/api/hotels?limit=5&offset=0
```

- Load thumbnail images from a given hotel (you can find the hotel picture path in the
_./mock-data/hotels-data.json_ file in each hotel entry under the _thumbNailUrl_ field):

```
http://localhost:8080/thumbnails/16950_158_t.jpg
```

# Contributors welcome

It would be great to get more hotels mock-data into this api, if you are keen on adding more entries
please fork this project and start adding them, you can find: 

- Hotels data in the following path:
_./mock-data/hotels-data.json_.

- Thumbnail images are stored under the following folder:
_./public/thumbnails_.

# Acknowledge

Original JSON feed extracted from this [apigee/DevJam](https://github.com/apigee/DevJam/blob/master/Resources/hotels-data.json) Github project.

The API is implemented in Go with the [gorilla/mux](https://github.com/gorilla/mux) router.

# About Basefactor + Lemoncode

We are an innovating team of Javascript experts, passionate about turning your ideas into robust products.

[Basefactor, consultancy by Lemoncode](http://www.basefactor.com) provides consultancy and coaching services.

[Lemoncode](http://lemoncode.net/services/en/#en-home) provides training services.

For the LATAM/Spanish audience we are running an Online Front End Master degree, more info: http://lemoncode.net/master-frontend
