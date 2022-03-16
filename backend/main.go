package main

/*fmt.Printf("%+v\n", glx)*/
/*Current Err 45*/

import (
	"cloud.google.com/go/firestore"
	"context"
	firebase "firebase.google.com/go"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
)

/*


	TYPES AND DECLARATIONS


*/

type sigStruct struct {
	Sig string `json:"sig"`
}

type tokenStruct struct {
	Expiry int64  `firestore:"exp"`
	User   string `firestore:"user"`
}

type levelJson struct {
	Level string `json:"level"`
}

type IMXGalaxiator struct {
	User string `json:"user"`
	Name string `json:"name"`
	// Can add other stuff here if we ever need to get it from IMX
}

type Galaxiator struct {
	GivenName     string `firestore:"givn" json:"givenname"`
	NickName      string `firestore:"nick" json:"nickname"`
	Level         string `firestore:"levl" json:"level"`
	SharedDiscord bool   `firestore:"disc"`
}

type Answer struct {
	NextLevel string `firestore:"x,omitempty" json:"next"`
	Text      string `firestore:"t,omitempty" json:"text"`
}

type Level struct {
	Text       string            `firestore:"t,omitempty" json:"text"`
	Outcome    string            `firestore:"o,omitempty" json:"outc"`
	Locked     bool              `firestore:"lock,omitempty"`
	Answers    map[string]Answer `firestore:"a" json:"ans"`
	NameEarned string            `firestore:"n" json:"name"`
	CyberSand  int               `firestore:"v,omitempty" json:"$glxr"`
	ID         string            `json:"level_id"`
}

type GlxList map[string]Galaxiator

var db *firestore.Client
var ctx context.Context

var discordClient *discordgo.Session
var discordToken = "SECURE"

func main() {
	fmt.Println("Loading Firestore...")
	ctx = context.Background()
	app, err := firebase.NewApp(ctx, nil)
	if err != nil {
		log.Fatalln(err)
	}
	db, err = app.Firestore(ctx)
	if err != nil {
		log.Fatalln(err)
	}
	defer db.Close()
	fmt.Println("Done!")

	fmt.Println("Activating Discord bot..")
	discordClient, err = discordgo.New("Bot " + discordToken)
	if err != nil {
		log.Fatalln(err)
	}
	err = discordClient.Open()
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println("Done!")

	fmt.Println("Initialising http multiplexer...")
	r := mux.NewRouter()

	r.HandleFunc("/login", LoginHandler).Methods("POST")
	// Takes a signs message, and returns an auth token

	r.HandleFunc("/savegiven", SaveGivenNameHandler).Methods("POST")
	// Updates a Galaxiator's given name

	r.HandleFunc("/glxinfo", GlxInfoHandler).Methods("POST")
	// Takes a JSON list of glx ID strings and returns data

	r.HandleFunc("/setl", SetLevelHandler).Methods("POST")
	// Takes a {"level":"newlevel"} JSON and returns data of the new level

	r.HandleFunc("/getl", GetGlxLevelDataHandler).Methods("POST")
	// Takes a request with glx headers and returns data of that Galaxiator's level

	r.HandleFunc("/story/{glxid}", StoryHandler).Methods("GET")
	// Returns the full story text for a given glxid

	r.HandleFunc("/discoshare", ShareOnDiscordHandler).Methods("GET")
	// Shares the nickname to the Discord server - only once per toon!

	r.PathPrefix("/").HandlerFunc(HomeHandler)

	fmt.Println("Done!")

	port := os.Getenv("PORT")
	if port == "" {
		port = "8001"
	}

	fmt.Println("Listening on http://127.0.0.1:0" + port)
	log.Println(http.ListenAndServe(":"+port, corsHandler(r)))
}
