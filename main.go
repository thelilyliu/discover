package main

import (
	"bytes"
	"encoding/json"
	"html/template"
	"image"
	"image/jpeg"
	"log"
	"math"
	"net/http"
	"os"
	"os/exec"
	"path"
	"strconv"
	"strings"
	"time"

	"github.com/disintegration/imaging"
	"github.com/gorilla/context"
	"github.com/husobee/vestigo"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type User struct {
	UserID    string  `json:"userID"`
	Name      string  `json:"name"`
	Option    int     `json:"option"`
	FirstLat  float64 `json:"firstLat"`
	FirstLong float64 `json:"firstLong"`
	CurrLat   float64 `json:"currLat"`
	CurrLong  float64 `json:"currLong"`
	Points    []Point `json:"points"`
}

type Point struct {
	PointID   string   `json:"pointID"`
	Latitude  float64  `json:"latitude"`
	Longitude float64  `json:"longitude"`
	Address   string   `json:"address"`
	Keywords  []string `json:"keywords"`
	Results   []Result `json:"results"`
	Outcome   float64  `json:"outcome"`
}

type Response struct {
	Images []Image `json:"images"`
}

type Image struct {
	Classifiers []Classifier `json:"classifiers"`
}

type Classifier struct {
	Classes []Result `json:"classes"`
}

type Result struct {
	Class string  `json:"class"`
	Score float64 `json:"score"`
}

type Person struct {
	Name string `json:"name"`
}

func main() {
	router := vestigo.NewRouter()

	// set up router global CORS policy
	router.SetGlobalCors(&vestigo.CorsAccessControl{
		AllowOrigin:      []string{"*"},
		AllowCredentials: false,
		MaxAge:           3600 * time.Second,
	})

	fileServerAssets := http.FileServer(http.Dir("assets"))
	router.Get("/assets/*", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Vary", "Accept-Encoding")
		w.Header().Set("Cache-Control", "public, max-age=86400")
		w.Header().Set("Server", "GWS")
		r.URL.Path = strings.TrimPrefix(r.URL.Path, "/assets")
		fileServerAssets.ServeHTTP(w, r)
	})

	// User
	router.Get("/insertUser", insertUser)
	router.Get("/getUser/:userID", getUser)
	router.Post("/updateUser/:userID", updateUser)
	router.Get("/userInRange/:userID", userInRange)

	// Point
	router.Get("/insertPoint", insertPoint)
	router.Get("/getPoint/:pointID", getPoint)
	router.Post("/updatePoint/:pointID", updatePoint)
	router.Post("/pointCheckResult/:pointID", pointCheckResult)

	// Other
	router.Get("/sendOption/:userID/:option", sendOption)
	router.Post("/postImage", postImage)
	router.Get("/", viewIndex)

	log.Println("Listening...")

	if err := http.ListenAndServe(":2323", context.ClearHandler(router)); err != nil {
		log.Println(err)
	}
}

/*
  ========================================
  Other
  ========================================
*/

func viewIndex(w http.ResponseWriter, r *http.Request) {
	returnCode := 0
	var err error

	setHeader(w)
	person := Person{"Lily"}

	index := path.Join("assets/html", "index.html")
	content := path.Join("assets/html", "content.html")

	t, err := template.ParseFiles(index, content)
	if returnCode == 0 {
		if err != nil {
			returnCode = 1
		}
	}

	if returnCode == 0 {
		if err := t.ExecuteTemplate(w, "my-template", person); err != nil {
			returnCode = 2
		}
	}

	if returnCode != 0 {
		handleError("viewAdmin", returnCode, w)
	}
}

func sendOption(w http.ResponseWriter, r *http.Request) {
	returnCode := 0
	var err error

	user := new(User)
	user.UserID = vestigo.Param(r, "userID")

	if returnCode == 0 {
		if err = user.get(); err != nil {
			returnCode = 1
		}
	}

	if returnCode == 0 {
		if user.Option, err = strconv.Atoi(vestigo.Param(r, "option")); err != nil {
			returnCode = 2
		}
	}

	if returnCode == 0 {
		switch user.Option {
		case 2:
			user.Points = []Point{
				{
					PointID:   "2-1",
					Latitude:  43.661298141724274,
					Longitude: -79.4013948957703,
					Address:   "266 College St",
					// Keywords:  []string{"black building", "Burger King", "blue tiles"},
					Keywords: []string{"harbor", "shelter", "blue color"},
				},
				{
					PointID:   "2-2",
					Latitude:  43.6579468,
					Longitude: -79.4001475,
					Address:   "58 Willcocks St",
					Keywords:  []string{"red building", "streetlights", "houses"},
				},
			}

		case 4:
			user.Points = []Point{
				{
					PointID:   "4-1",
					Latitude:  43.666189501459584,
					Longitude: -79.39313348623045,
					Address:   "7 Queens Park",
					Keywords:  []string{"bus stop", "beige buildings", "iron gate"},
				},
				{
					PointID:   "4-2",
					Latitude:  43.6557259,
					Longitude: -79.3837337,
					Address:   "91 Dundas St W",
					Keywords:  []string{"glass building", "Canadian Tire", "brown building"},
				},
				{
					PointID:   "4-3",
					Latitude:  43.6541387,
					Longitude: -79.39227089999997,
					Address:   "299-317 Dundas St W",
					Keywords:  []string{"townhouses", "red building", "caution sign"},
				},
				{
					PointID:   "4-4",
					Latitude:  43.6577481,
					Longitude: -79.40007659999998,
					Address:   "464 Spadina Ave",
					Keywords:  []string{"convenience store", "bus stop", "newspaper stands"},
				},
			}

		case 6:
			user.Points = []Point{
				{
					PointID:   "6-1",
					Latitude:  43.6564647,
					Longitude: -79.40770120000002,
					Address:   "85 St George St",
					Keywords:  []string{"iron gates", "bike rack", "streetlights"},
				},
				{
					PointID:   "6-2",
					Latitude:  43.6519345,
					Longitude: -79.40436799999998,
					Address:   "81 St Mary St",
					Keywords:  []string{"grey building", "university", "orange buildings"},
				},
				{
					PointID:   "6-3",
					Latitude:  43.6544597,
					Longitude: -79.39072090000002,
					Address:   "80 Dr Emily Stowe Way",
					Keywords:  []string{"church", "glass building", "restaurant"},
				},
				{
					PointID:   "6-4",
					Latitude:  43.66051820000001,
					Longitude: -79.3874169,
					Address:   "290-292 Dundas St W",
					Keywords:  []string{"haircutter", "restaurant", "art store"},
				},
				{
					PointID:   "6-5",
					Latitude:  43.66679603006645,
					Longitude: -79.39099288518065,
					Address:   "686 Dundas St W",
					Keywords:  []string{"graffiti", "red windows", "trailers"},
				},
				{
					PointID:   "6-6",
					Latitude:  43.6639237,
					Longitude: -79.39837090000003,
					Address:   "440 College St",
					Keywords:  []string{"shops", "streetlights", "bus stop"},
				},
			}
		}

		if err := json.NewEncoder(w).Encode(user); err != nil {
			returnCode = 3
		}
	}

	if returnCode != 0 {
		log.Println(err)
		handleError("postImage", returnCode, w)
	}
}

/*
func postImage(w http.ResponseWriter, r *http.Request) {
	returnCode := 0
	var err error

	data, err := ioutil.ReadAll(r.Body)
	r.Body.Close()

	if returnCode == 0 {
		if err != nil {
			returnCode = 1
		}
	}

	if returnCode == 0 {
		if err = ioutil.WriteFile("moment.jpg", data, 0666); err != nil {
			returnCode = 2
		}
	}

	if returnCode == 0 {
		if err := json.NewEncoder(w).Encode(true); err != nil {
			returnCode = 3
		}
	}

	if returnCode != 0 {
		log.Println(err)
		handleError("postImage", returnCode, w)
	}
}
*/

func postImage(w http.ResponseWriter, r *http.Request) {
	returnCode := 0
	var err error

	file, _, err := r.FormFile("uploadFile")
	defer file.Close()

	if returnCode == 0 {
		if err != nil {
			returnCode = 1
		}
	}

	originalImage, _, err := image.Decode(file)

	if returnCode == 0 {
		if err != nil {
			returnCode = 2
		}
	}

	imageCompressed := imaging.Resize(originalImage, 600, 0, imaging.Linear)

	imageFile, err := os.Create("moment.jpg")
	defer imageFile.Close()

	if returnCode == 0 {
		if err != nil {
			returnCode = 3
		}
	}

	if returnCode == 0 {
		if err = jpeg.Encode(imageFile, imageCompressed, &jpeg.Options{90}); err != nil {
			returnCode = 4
		}
	}

	if returnCode == 0 {
		if err = json.NewEncoder(w).Encode(true); err != nil {
			returnCode = 5
		}
	}

	if returnCode != 0 {
		handleError("postImage", returnCode, w)
	}
}

/*
  ========================================
  Special
  ========================================
*/

func (p *Point) checkResult() bool {
	for _, result := range p.Results {
		for _, keyword := range p.Keywords {
			if result.Class == keyword {
				p.Outcome += (result.Score * 10.0)
			}
		}
	}

	return p.Outcome != 0
}

func (p *Point) classifyImage() {
	var out bytes.Buffer
	var stderr bytes.Buffer

	cmd := exec.Command("./script.sh")
	cmd.Stdout = &out
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		log.Println(err)
	}

	responseFile, err := os.Open("response.json")
	if err != nil {
		log.Println(err)
	}

	response := new(Response)
	if err := json.NewDecoder(responseFile).Decode(&response); err != nil {
		log.Println(err)
	}

	p.Results = response.Images[0].Classifiers[0].Classes
}

func (u *User) inRange() bool {
	c := u.nextPoint()
	x := u.CurrLat - c.Latitude
	y := u.CurrLong - c.Longitude
	dist := math.Sqrt(math.Pow(x, 2) + math.Pow(y, 2))

	return dist <= 100
}

func (u *User) nextPoint() *Point {
	for _, point := range u.Points {
		if point.Outcome == 0 {
			return &point
		}
	}

	return &Point{}
}

/*
  ========================================
  HTTP: User
  ========================================
*/

func insertUser(w http.ResponseWriter, r *http.Request) {
	returnCode := 0
	var err error

	user := new(User)

	if returnCode == 0 { // 1: Decode
		if err = json.NewDecoder(r.Body).Decode(user); err != nil {
			returnCode = 1
		}
	}

	if returnCode == 0 { // 2: Insert
		if err = user.insert(); err != nil {
			returnCode = 2
		}
	}

	if returnCode == 0 { // 3: Encode
		if err = json.NewEncoder(w).Encode(user.UserID); err != nil {
			returnCode = 3
		}
	}

	if returnCode != 0 {
		log.Println(err)
		handleError("insertUser", returnCode, w)
	}
}

func getUser(w http.ResponseWriter, r *http.Request) {
	returnCode := 0
	var err error

	user := new(User)
	user.UserID = vestigo.Param(r, "userID")

	if returnCode == 0 { // 1: Get
		if err = user.get(); err != nil {
			returnCode = 1
		}
	}

	if returnCode == 0 { // 2: Encode
		if err = json.NewEncoder(w).Encode(user); err != nil {
			returnCode = 2
		}
	}

	if returnCode != 0 {
		log.Println(err)
		handleError("getUser", returnCode, w)
	}
}

func updateUser(w http.ResponseWriter, r *http.Request) {
	returnCode := 0
	var err error

	user := new(User)
	user.UserID = vestigo.Param(r, "userID")

	if returnCode == 0 { // 1: Decode
		if err = json.NewDecoder(r.Body).Decode(user); err != nil {
			returnCode = 1
		}
	}

	if returnCode == 0 { // 2: Update
		if err = user.update(); err != nil {
			returnCode = 2
		}
	}

	if returnCode == 0 { // 3: Encode
		if err = json.NewEncoder(w).Encode(user); err != nil {
			returnCode = 3
		}
	}

	if returnCode != 0 {
		log.Println(err)
		handleError("updateUser", returnCode, w)
	}
}

// =======================================

func userInRange(w http.ResponseWriter, r *http.Request) {
	returnCode := 0
	var err error

	user := new(User)
	user.UserID = vestigo.Param(r, "userID")

	if returnCode == 0 { // 1: Decode
		if err = json.NewDecoder(r.Body).Decode(user); err != nil {
			returnCode = 1
		}
	}

	inRange := user.inRange() // **** Special ****

	if returnCode == 0 { // 2: Encode
		if err = json.NewEncoder(w).Encode(inRange); err != nil {
			returnCode = 2
		}
	}

	if returnCode != 0 {
		log.Println(err)
		handleError("updateUser", returnCode, w)
	}
}

/*
  ========================================
  HTTP: Point
  ========================================
*/

func insertPoint(w http.ResponseWriter, r *http.Request) {
	returnCode := 0
	var err error

	point := new(Point)

	if returnCode == 0 { // 1: Decode
		if err = json.NewDecoder(r.Body).Decode(point); err != nil {
			returnCode = 1
		}
	}

	if returnCode == 0 { // 2: Insert
		if err = point.insert(); err != nil {
			returnCode = 2
		}
	}

	if returnCode == 0 { // 3: Encode
		if err = json.NewEncoder(w).Encode(point.PointID); err != nil {
			returnCode = 3
		}
	}

	if returnCode != 0 {
		log.Println(err)
		handleError("insertPoint", returnCode, w)
	}
}

func getPoint(w http.ResponseWriter, r *http.Request) {
	returnCode := 0
	var err error

	point := new(Point)
	point.PointID = vestigo.Param(r, "pointID")

	if returnCode == 0 { // 1: Get
		if err = point.get(); err != nil {
			returnCode = 1
		}
	}

	if returnCode == 0 { // 2: Encode
		if err = json.NewEncoder(w).Encode(point); err != nil {
			returnCode = 2
		}
	}

	if returnCode != 0 {
		log.Println(err)
		handleError("getPoint", returnCode, w)
	}
}

func updatePoint(w http.ResponseWriter, r *http.Request) {
	returnCode := 0
	var err error

	point := new(Point)
	point.PointID = vestigo.Param(r, "pointID")

	if returnCode == 0 { // 1: Decode
		if err = json.NewDecoder(r.Body).Decode(point); err != nil {
			returnCode = 1
		}
	}

	if returnCode == 0 { // 2: Update
		if err = point.update(); err != nil {
			returnCode = 2
		}
	}

	if returnCode == 0 { // 3: Encode
		if err = json.NewEncoder(w).Encode(point); err != nil {
			returnCode = 3
		}
	}

	if returnCode != 0 {
		log.Println(err)
		handleError("updatePoint", returnCode, w)
	}
}

// =======================================

func pointCheckResult(w http.ResponseWriter, r *http.Request) {
	returnCode := 0
	var err error

	point := new(Point)

	if returnCode == 0 { // 1: Decode
		if err = json.NewDecoder(r.Body).Decode(point); err != nil {
			returnCode = 1
		}
	}

	point.classifyImage()              // **** Special ****
	checkResult := point.checkResult() // **** Special ****

	if returnCode == 0 && checkResult { // 2: Update
		if err = point.update(); err != nil {
			returnCode = 2
		}
	}

	if returnCode == 0 { // 3: Encode
		if err = json.NewEncoder(w).Encode(checkResult); err != nil {
			returnCode = 3
		}
	}

	if returnCode != 0 {
		log.Println(err)
		handleError("updatePoint", returnCode, w)
	}
}

/*
  ========================================
  Database: User
  ========================================
*/

func (u *User) insert() error {
	// create new MongoDB session
	collection, session := initMongoDB("user")
	defer session.Close()

	// initialize fields
	mongoID := bson.NewObjectId().String()
	u.UserID = mongoID[13 : len(mongoID)-2]

	return collection.Insert(u)
}

func (u *User) get() error {
	// create new MongoDB session
	collection, session := initMongoDB("user")
	defer session.Close()

	selector := bson.M{"userid": u.UserID}
	return collection.Find(selector).One(u)
}

func (u *User) update() error {
	// create new MongoDB session
	collection, session := initMongoDB("user")
	defer session.Close()

	selector := bson.M{"userid": u.UserID}
	change := bson.M{}

	update := bson.M{"$set": &change}
	return collection.Update(selector, update)
}

/*
  ========================================
  Database: Point
  ========================================
*/

func (p *Point) insert() error {
	// create new MongoDB session
	collection, session := initMongoDB("point")
	defer session.Close()

	// initialize fields
	mongoID := bson.NewObjectId().String()
	p.PointID = mongoID[13 : len(mongoID)-2]

	return collection.Insert(p)
}

func (p *Point) get() error {
	// create new MongoDB session
	collection, session := initMongoDB("point")
	defer session.Close()

	selector := bson.M{"pointid": p.PointID}
	return collection.Find(selector).One(p)
}

func (p *Point) update() error {
	// create new MongoDB session
	collection, session := initMongoDB("point")
	defer session.Close()

	selector := bson.M{"pointid": p.PointID}
	change := bson.M{"results": p.Results, "outcome": p.Outcome}

	update := bson.M{"$set": &change}
	return collection.Update(selector, update)
}

/*
  ========================================
  Basic
  ========================================
*/

func handleError(funcName string, returnCode int, w http.ResponseWriter) {
	w.WriteHeader(555)
	message := funcName + " " + strconv.Itoa(returnCode)

	if err := json.NewEncoder(w).Encode(message); err != nil {
		log.Println(err)
	}
}

func initMongoDB(collectionName string) (*mgo.Collection, *mgo.Session) {
	session, err := mgo.Dial("127.0.0.1")

	if err != nil {
		log.Println(err)
	}

	return session.DB("discover").C(collectionName), session
}

func setHeader(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Cache-control", "no-cache, no-store, max-age=0, must-revalidate")
	w.Header().Set("Expires", "Fri, 01 Jan 1990 00:00:00 GMT")
	w.Header().Set("Server", "GWS")
}
