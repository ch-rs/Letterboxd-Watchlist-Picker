package film

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gocolly/colly/v2"
	"github.com/gocolly/colly/v2/extensions"
)

//Film struct for http response
type film struct {
	Slug       string `json:"slug"`         //url of film
	Image      string `json:"image_url"`    //url of image
	Year       string `json:"release_year"`
	Name       string `json:"film_name"`
	Length     string `json:"film_length"`
	OriginalIndex int    `json:"original_index"` //original position in the list before shuffling (-1 if not in top)
}

// Add a new response struct that includes URLs
type filmResponse struct {
	Films []film  `json:"films"`
	URLs  []string `json:"urls"`
}

//struct for channel to send film and whether is has finshed a user
type filmSend struct {
	film film //film to be sent over channel
	done bool //if user is done
}

type toIgnore struct {
	unreleased bool
	short bool
	feature bool
}

type nothingReason int

const (
	INTERSECT = iota
	UNION
)

const FEATURELENGTH = 60

type nothingError struct {
	reason nothingReason
}

func (e *nothingError) ToString() string {
	switch e.reason {
	case INTERSECT:
		return "empty intersect"
	case UNION:
		return "empty union"
	default:
		return "big error"
	}

}

func (e *nothingError) Error() string {
	return e.ToString()
}

const urlscrape = "https://letterboxd.com/ajax/poster" //first part of url for getting full info on film
const urlEnd = "std/125x187/"            // second part of url for getting full info on film
const site = "https://letterboxd.com"


// func main() {
// 	getFilmHandler := http.HandlerFunc(Handler)
// 	http.Handle("/film", getFilmHandler)
// 	log.Println("serving at :8080")
// 	port := os.Getenv("PORT")
// 	if port == "" {
// 		port = "8080"
// 		log.Printf("Defaulting to port %s", port)
// 	}

// 	log.Printf("Listening on port %s", port)
// 	http.ListenAndServe(":"+port, nil)
// }

var year int

func init() {
	year = time.Now().Year()
}



//Main handler func for request
func Handler(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	log.Println(year)
	query := r.URL.Query() //Get URL Params(type map)
	users, ok := query["users"]
	log.Println(len(users))
	if !ok || len(users) == 0 {
		http.Error(w, "no users", 400)
		return
	}
	_, inter := query["intersect"]
	ignore, _ := query["ignore"]
	var ignoreing = toIgnore{}

	if len(ignore) > 0 {
		ignoreing = whatToIgnore(ignore[0])
	}
	log.Println(ignoreing)

	var userFilm []film
	var err error
	var scrapedURLs []string

	if inter {
		if len(users) == 1 {
			userFilm, scrapedURLs, err = scrapeMain(users, false, ignoreing)
		} else {
			userFilm, scrapedURLs, err = scrapeMain(users, true, ignoreing)
		}
	} else {
		userFilm, scrapedURLs, err = scrapeMain(users, false, ignoreing)
	}
	if err != nil {
		var e *nothingError
		if errors.As(err, &e) {
			switch e.reason {
			case INTERSECT:
				http.Error(w, "Intersect error", 406)
				return
			case UNION:
				http.Error(w, "Union error", 404)
				return
			}
		}
	}

	// Create response with films and URLs
	response := filmResponse{
		Films: userFilm,
		URLs:  scrapedURLs,
	}

	js, err := json.Marshal(response)
	if err != nil {
		http.Error(w, "internal error", 500)
		return
	}
	w.Write(js)
}


//main scraping function
func scrapeMain(users []string, intersect bool, ignoreList toIgnore) ([]film, []string, error) {
	var user int = 0          //conuter for number of users increses by one when a users page starts being scraped decreses when user has finished think kinda like a semaphore
	var totalFilms []film     //final list to hold all film
	ch := make(chan filmSend) //channel to send films over
	
	// Track URLs being scraped
	var scrapedURLs []string
	
	// start go routine to scrape each user
	for _, a := range users {
		log.Println(a)
		user++
		var url string
		
		if strings.Contains(a, "/") {
			if (strings.Contains(a,"actor/") || strings.Contains(a,"director/")) {
				url = site + "/" + a
				if ignoreList.short || ignoreList.feature {
					go scrapeActorWithLength(a, ch)
				} else {
					go scrapeActor(a, ch)
				}
			} else {
				if strings.Contains(a, "/list/") {
					url = site + "/" + a + "/by/added-earliest/"
				} else {
					strslice := strings.Split(a, "/")
					url = site + "/" + strslice[0] + "/list/" + strslice[1] + "/by/added-earliest/"
				}
				
				if ignoreList.short || ignoreList.feature {
					go scrapeListWithLength(a, ch)
				} else {
					go scrapeList(a, ch)
				}
			}
		} else {
			url = site + "/" + a + "/watchlist"
			if ignoreList.short || ignoreList.feature {
				go scrapeUserWithLength(a, ch)
			} else {
				go scrapeUser(a, ch)
			}
		}
		
		scrapedURLs = append(scrapedURLs, url)
	}
	for {
		userFilm := <-ch
		if userFilm.done { //if users channel is don't then the scapre for that user has finished so decrease the user count
			user--
			if user == 0 {
				break
			}
		} else {
			totalFilms = append(totalFilms, userFilm.film) //append feilm recieved over channel to list
		}

	}

	//chose random film from list
	if len(totalFilms) == 0 {
		// Return empty film list and error
		return nil, nil, &nothingError{reason: UNION}
	}
	log.Print("results")

	var filmList []film
	if intersect {
		intersectList := getintersect(totalFilms,len(users))
		length := len(intersectList)
		if length == 0 {
			return nil, nil, &nothingError{reason: INTERSECT}
		}
		filmList = intersectList
	} else {
		filmList = totalFilms
	}
	if ignoreList.unreleased {
		filmList = ignoreUnrelased(filmList)
	}
	if ignoreList.short {
		filmList = ignoreShorts(filmList)
	}
	if ignoreList.feature {
		filmList = ignoreFeature(filmList)
	}

	rand.Shuffle(len(filmList), func(i, j int) { filmList[i], filmList[j] = filmList[j], filmList[i] })

	numFilms := len(filmList)
	if numFilms > 20 {
		numFilms = 20
	}

	// Return the URLs along with the films
	return filmList[:numFilms], scrapedURLs, nil
}


func scrapeUserWithLength(userName string, ch chan filmSend) {
	url := site + "/" + userName + "/watchlist"
	scrapeWithLength(url, ch)
}

func scrapeUser(userName string, ch chan filmSend) {
	url := site + "/" + userName + "/watchlist"
	scrape(url, ch)
}

func scrapeListWithLength(listNameIn string, ch chan filmSend) {
	url := ""
	listname := strings.ToLower(listNameIn)

	if strings.Contains(listname, "/list/") {
		url = site + "/" + listname + "/by/added-earliest/"
	} else {
		strslice := strings.Split(listname, "/") //strslice[0] is user name strslice[1] is listname
		url = site + "/" + strslice[0] + "/list/" + strslice[1] + "/by/added-earliest/"

	}
	scrapeWithLength(url, ch)
}

func scrapeList(listNameIn string, ch chan filmSend) {
	url := ""
	listname := strings.ToLower(listNameIn)

	if strings.Contains(listname, "/list/") {
		url = site + "/" + listname + "/by/added-earliest/"
	} else {
		strslice := strings.Split(listname, "/") //strslice[0] is user name strslice[1] is listname
		url = site + "/" + strslice[0] + "/list/" + strslice[1] + "/by/added-earliest/"
	}

	scrape(url, ch)
}


func scrape(url string, ch chan filmSend) {
	siteToVisit := url
	posterCount := 0  // Track the number of posters processed

	ajc := colly.NewCollector()
	ajc.OnHTML("div.film-poster", func(e *colly.HTMLElement) { //secondard cleector to get main data for film
		name := e.Attr("data-film-name")
		slug := e.Attr("data-film-link")
		img := e.ChildAttr("img", "src")
		year := e.Attr("data-film-release-year")
		
		// Set original index for the first 3 films, -1 for the rest
		originalIndex := -1
		if posterCount < 3 {
			originalIndex = posterCount
		}
		posterCount++
		
		tempfilm := film{
			Slug:  (site + slug),
			Image: makeBigger(img),
			Year: year,
			Name:  name,
			OriginalIndex: originalIndex,
		}
		ch <- ok(tempfilm)
	})
	c := colly.NewCollector()
	c.Limit(&colly.LimitRule{DomainGlob: "*", Parallelism: 100})
	c.OnHTML(".poster-container", func(e *colly.HTMLElement) { //primary scarer to get url of each film that contian full information
		e.ForEach("div.film-poster", func(i int, ein *colly.HTMLElement) {
			slug := ein.Attr("data-target-link")
			ajc.Visit(urlscrape + slug + urlEnd) //start go routine to collect all film data
		})

	})
	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		if strings.Contains(link, "/page") {
			e.Request.Visit(e.Request.AbsoluteURL(link))
		}
	})

	c.Visit(siteToVisit)
	c.Wait()
	ajc.Wait()
	ch <- done() // users has finished so send done through channel

}


func scrapeWithLength(url string, ch chan filmSend) { //is slower so is own function
	siteToVisit := url
	posterCount := 0  // Track the number of posters processed
	
	ajc := colly.NewCollector()
	extensions.RandomUserAgent(ajc)
	ajc.OnHTML("div#film-page-wrapper", func(e *colly.HTMLElement) {
		name := e.ChildText("span.frame-title")
		slug := e.ChildAttr("div.film-poster","data-film-link")
		img := e.ChildAttr("img", "src")
		year := e.ChildAttr("div.film-poster","data-film-release-year")
		lenght := e.ChildText("p.text-footer")
		
		// Set original index for the first 3 films, -1 for the rest
		originalIndex := -1
		if posterCount < 3 {
			originalIndex = posterCount
		}
		posterCount++
		
		tempfilm := film{
			Slug:  (site + slug),
			Image: img,
			Year: year,
			Name:  name,
			Length: strings.TrimSpace(before(lenght,"mins")),
			OriginalIndex: originalIndex,
		}
		ch <- ok(tempfilm)
	})

	c := colly.NewCollector()
	c.Limit(&colly.LimitRule{DomainGlob: "*", Parallelism: 100})
	extensions.RandomUserAgent(c)
	c.OnHTML(".poster-container", func(e *colly.HTMLElement) { //primary scarer to get url of each film that contian full information
		e.ForEach("div.film-poster", func(i int, ein *colly.HTMLElement) {
			slug := ein.Attr("data-target-link")
			ajc.Visit(site + slug) //start go routine to collect all film data
		})

	})
	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		if strings.Contains(link, "/page") {
			e.Request.Visit(e.Request.AbsoluteURL(link))
		}
	})

	c.Visit(siteToVisit)
	c.Wait()
	ajc.Wait()
	ch <- done()

}

func scrapeActor(actor string, ch chan filmSend) {
	siteToVisit := site + "/" + actor
	fmt.Println(siteToVisit)
	posterCount := 0  // Track the number of posters processed

	c := colly.NewCollector()
	c.Limit(&colly.LimitRule{DomainGlob: "*", Parallelism: 100})
	c.OnHTML("div.film-poster", func(e *colly.HTMLElement) { //primary scarer to get url of each film that contian full information
		name := e.Attr("data-film-name")
		slug := e.Attr("data-film-link")
		img := e.ChildAttr("img", "src")
		year := e.Attr("data-film-release-year")
		
		// Set original index for the first 3 films, -1 for the rest
		originalIndex := -1
		if posterCount < 3 {
			originalIndex = posterCount
		}
		posterCount++
		
		tempfilm := film{
			Slug:  (site + slug),
			Image: makeBiggerActor(img),
			Year: year,
			Name:  name,
			OriginalIndex: originalIndex,
		}
		ch <- ok(tempfilm)
	})

	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		if strings.Contains(link, "/page") {
			e.Request.Visit(e.Request.AbsoluteURL(link))
		}
	})

	c.Visit(siteToVisit)
	c.Wait()
	ch <- done() // users has finished so send done through channel

}

func scrapeActorWithLength(actor string, ch chan filmSend) {
	siteToVisit := site + "/" + actor
	log.Println(siteToVisit)
	posterCount := 0  // Track the number of posters processed

	c := colly.NewCollector()
	c.Limit(&colly.LimitRule{DomainGlob: "*", Parallelism: 100})
	extensions.RandomUserAgent(c)
	ajc := colly.NewCollector()
	extensions.RandomUserAgent(ajc)
	ajc.OnHTML("div#film-page-wrapper", func(e *colly.HTMLElement) {
		name := e.ChildText("span.frame-title")
		slug := e.ChildAttr("div.film-poster","data-film-link")
		img := e.ChildAttr("img", "src")
		year := e.ChildAttr("div.film-poster","data-film-release-year")
		lenght := e.ChildText("p.text-footer")
		
		// Set original index for the first 3 films, -1 for the rest
		originalIndex := -1
		if posterCount < 3 {
			originalIndex = posterCount
		}
		posterCount++
		
		tempfilm := film{
			Slug:  (site + slug),
			Image: img,
			Year: year,
			Name:  name,
			Length: strings.TrimSpace(before(lenght,"mins")),
			OriginalIndex: originalIndex,
		}
		ch <- ok(tempfilm)
	})

	c.OnHTML(".poster-container", func(e *colly.HTMLElement) {
		e.ForEach("div.film-poster", func(i int, ein *colly.HTMLElement) {
			slug := ein.Attr("data-film-link")
			ajc.Visit(site + slug)
		})

	})

	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		if strings.Contains(link, "/page") {
			e.Request.Visit(e.Request.AbsoluteURL(link))
		}
	})

	c.Visit(siteToVisit)
	c.Wait()
	ajc.Wait()
	ch <- done() // users has finished so send done through channel

}

func ok(f film) filmSend {
	return filmSend{
		film: f,
		done: false,
	}
}

func done() filmSend {
	return filmSend{
		film: film{},
		done: true,
	}
}

func getintersect(filmSlice []film, numOfUsers int) []film {
	keys := make(map[film]int)
	list := []film{}
	for _, entry := range filmSlice {
		i, _ := keys[entry]
		if i < (numOfUsers - 1) {
			keys[entry] ++
		} else {
			list = append(list, entry)
		}
	}
	return list
}

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
}

func makeBiggerActor(url string) string {
	return strings.ReplaceAll(url, "-0-150-0-225-", "-0-230-0-345-")

}
func makeBigger(url string) string {
	return strings.ReplaceAll(url, "-0-125-0-187-", "-0-230-0-345-")
}

func ignoreUnrelased(filmSlice []film) []film {
	list := []film{}
	for _, entry := range filmSlice {
		if entry.Year == "" {
			continue
		}
		filmYear, _ := strconv.Atoi(entry.Year)
		if filmYear < year {
			list = append(list, entry)
		}
	}
	return list
}

func ignoreShorts(filmSlice []film) []film {
	list := []film{}
	for _, entry := range filmSlice {
		length, err := strconv.Atoi(entry.Length)
		if err != nil {
			continue
		}
		if length >= FEATURELENGTH {
			list = append(list, entry)
		}
	}
	return list
}

func ignoreFeature(filmSlice []film) []film {
	list := []film{}
	for _, entry := range filmSlice {
		length, err := strconv.Atoi(entry.Length)
		if err != nil {
			continue
		}
		if length < FEATURELENGTH {
			list = append(list, entry)
		}
	}
	return list
}

func before(value string, a string) string {
    // Get substring before a string.
    pos := strings.Index(value, a)
    if pos == -1 {
        return ""
    }
    return value[0:pos]
}

func contains(s []string, e string) bool {
    for _, a := range s {
        if a == e {
            return true
        }
    }
    return false
}

func whatToIgnore(ignoreString string) toIgnore {
	ignoreList := strings.Split(ignoreString, ",")
	return toIgnore{
		unreleased: contains(ignoreList,"unreleased"),
		short: contains(ignoreList, "shorts"),
		feature: contains(ignoreList, "feature"),
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
