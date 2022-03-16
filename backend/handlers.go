package main

import (
	"cloud.google.com/go/firestore"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"google.golang.org/api/iterator"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

var DiscordChannels = map[string]string{
	"staff":      "A",
	"gompy":      "B",
	"general":    "C",
	"eleusinion": "D",
	"bot-dev":    "E",
}

// TODO All DB integration into separate file, in case we need to migrate
// TODO All ErrXX: need to log actual err

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	var from string
	if len(r.Header.Get("X-Forwarded-For")) >= 7 {
		from = r.Header.Get("X-Forwarded-For")
	} else {
		from = r.RemoteAddr
	}
	log.Println("404: Request to ", r.RequestURI, " from ", from)
	w.WriteHeader(404)
	w.Header().Set("Content-Type", "application/json")
	resp := make(map[string]string)
	resp["message"] = "These are not the droids you are looking for."
	jsonResp, _ := json.Marshal(resp)
	_, _ = w.Write(jsonResp)
	return
}

func GlxInfoHandler(w http.ResponseWriter, r *http.Request) {
	glxList := make(GlxList)
	w.Header().Set("Content-Type", "application/json")

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println("Err26: Invalid request body")
		w.WriteHeader(400)
		w.Header().Set("Content-Type", "application/json")
		resp := make(map[string]string)
		resp["message"] = "Invalid request body."
		jsonResp, _ := json.Marshal(resp)
		_, _ = w.Write(jsonResp)
		return
	}

	var glxListSlice []string
	err = json.Unmarshal(body, &glxListSlice)
	if err != nil {
		log.Println("Err35: Cannot unmarshall body ", body)
		w.WriteHeader(400)
		w.Header().Set("Content-Type", "application/json")
		resp := make(map[string]string)
		resp["message"] = "Cannot unmarshall body"
		jsonResp, _ := json.Marshal(resp)
		_, _ = w.Write(jsonResp)
		return
	}

	// For glx in list...
	for _, glx := range glxListSlice {
		glxRecord := Galaxiator{}

		dsnap1, err := db.Collection("glx").Doc(glx).Get(ctx)
		if err != nil {
			log.Println(err)
			continue
		} else {
			dsnap1.DataTo(&glxRecord)
			glxList[glx] = glxRecord
		}

	}
	jsonResp, err := json.Marshal(glxList)
	if err != nil {
		w.WriteHeader(500)
		w.Header().Set("Content-Type", "application/json")
		resp := make(map[string]string)
		resp["message"] = "Cannot marshall response."
		jsonResp, _ := json.Marshal(resp)
		_, _ = w.Write(jsonResp)
		return
	}
	_, _ = w.Write(jsonResp)
	return
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println("Err27: Invalid request body")
		w.WriteHeader(400)
		w.Header().Set("Content-Type", "application/json")
		resp := make(map[string]string)
		resp["message"] = "Invalid request body."
		jsonResp, _ := json.Marshal(resp)
		_, _ = w.Write(jsonResp)
		return
	}

	sigJson := sigStruct{}
	jsonErr := json.Unmarshal(body, &sigJson)
	if jsonErr != nil {
		log.Println("Err36: Cannot unmarshall body ", body)
		w.WriteHeader(400)
		w.Header().Set("Content-Type", "application/json")
		resp := make(map[string]string)
		resp["message"] = "Cannot unmarshall request."
		jsonResp, _ := json.Marshal(resp)
		_, _ = w.Write(jsonResp)
		return
	}

	addr, validateErr := addressFromSig(sigJson.Sig, []byte("Log in to begin your Galaxiators story."))
	if validateErr != nil {
		log.Println("Err31: Validation failed")
		w.WriteHeader(400)
		w.Header().Set("Content-Type", "application/json")
		resp := make(map[string]string)
		resp["message"] = "Sig not validated."
		jsonResp, _ := json.Marshal(resp)
		_, _ = w.Write(jsonResp)
		return
	}

	// Check for existing tokens and delete
	iter := db.Collection("tokens").Where("user", "==", addr).Documents(ctx)
	for {
		doc, err1 := iter.Next()
		if err1 == iterator.Done {
			break
		}
		if err1 != nil {
			log.Println("Err1 deleting existing tokens: ", err1)
		}
		_, err2 := doc.Ref.Delete(ctx)
		if err2 != nil {
			log.Println("Err2 deleting existing tokens: ", err2)
		}
	}

	// Create token
	now := time.Now()
	sec := now.Unix()
	sec += 43200 // Token valid for 12 hours
	doc, _, err := db.Collection("tokens").Add(ctx, map[string]interface{}{
		"user": addr,
		"exp":  sec,
	})
	if err != nil {
		w.WriteHeader(500)
		w.Header().Set("Content-Type", "application/json")
		resp := make(map[string]string)
		resp["message"] = "Cannot store token."
		jsonResp, _ := json.Marshal(resp)
		_, _ = w.Write(jsonResp)
		return
	}

	// Return user and token
	w.WriteHeader(200)
	resp := make(map[string]string)
	resp["user"] = addr
	resp["token"] = doc.ID
	jsonResp, err := json.Marshal(resp)
	if err != nil {
		w.WriteHeader(500)
		w.Header().Set("Content-Type", "application/json")
		resp = make(map[string]string)
		resp["message"] = "Cannot marshall response."
		jsonResp, _ = json.Marshal(resp)
		_, _ = w.Write(jsonResp)
		return
	}
	_, _ = w.Write(jsonResp)
	return
}

func SaveGivenNameHandler(w http.ResponseWriter, r *http.Request) {
	// Check auth token etc
	valid := validateRequest(r)
	if !valid {
		log.Println("Err32: Validation failed")
		w.WriteHeader(403)
		w.Header().Set("Content-Type", "application/json")
		resp := make(map[string]string)
		resp["message"] = "Request did not pass validation."
		jsonResp, _ := json.Marshal(resp)
		_, _ = w.Write(jsonResp)
		return
	}

	// Get givenname from the request
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println("Err28: Invalid request body")
		w.WriteHeader(400)
		w.Header().Set("Content-Type", "application/json")
		resp := make(map[string]string)
		resp["message"] = "Cannot read request body."
		jsonResp, _ := json.Marshal(resp)
		_, _ = w.Write(jsonResp)
		return
	}
	glxNewName := Galaxiator{}
	err = json.Unmarshal(body, &glxNewName)
	if err != nil {
		log.Println("Err37: Cannot unmarshall body ", body)
		w.WriteHeader(400)
		w.Header().Set("Content-Type", "application/json")
		resp := make(map[string]string)
		resp["message"] = "Cannot unmarshall request body."
		jsonResp, _ := json.Marshal(resp)
		_, _ = w.Write(jsonResp)
		return
	}

	// Check Galaxiator doesn't already have a givenname
	var glx Galaxiator

	dsnap1, err := db.Collection("glx").Doc(r.Header.Get("x-glx-id")).Get(ctx)
	if !dsnap1.Exists() {
		err = setGivenName(r.Header.Get("x-glx-id"), glxNewName.GivenName)
		if err != nil {
			w.WriteHeader(500)
			w.Header().Set("Content-Type", "application/json")
			resp := make(map[string]string)
			resp["message"] = "A: Cannot set name - " + err.Error()
			jsonResp, _ := json.Marshal(resp)
			_, _ = w.Write(jsonResp)
			return
		}
		w.WriteHeader(204)
		return
	}
	if err != nil {
		w.WriteHeader(500)
		w.Header().Set("Content-Type", "application/json")
		resp := make(map[string]string)
		resp["message"] = "Error getting Galaxiator data: " + err.Error()
		jsonResp, _ := json.Marshal(resp)
		_, _ = w.Write(jsonResp)
		return
	}
	err = dsnap1.DataTo(&glx)
	if err != nil {
		w.WriteHeader(500)
		w.Header().Set("Content-Type", "application/json")
		resp := make(map[string]string)
		resp["message"] = "Cannot parse Galaxiator data."
		jsonResp, _ := json.Marshal(resp)
		_, _ = w.Write(jsonResp)
		return
	}
	if glx.GivenName != "" {
		w.WriteHeader(409)
		w.Header().Set("Content-Type", "application/json")
		resp := make(map[string]string)
		resp["message"] = "You already have name, " + glx.GivenName + ", why would you need a new one? Hit refresh and try again."
		jsonResp, _ := json.Marshal(resp)
		_, _ = w.Write(jsonResp)
		return
	}

	// update the db
	err = setGivenName(r.Header.Get("x-glx-id"), glxNewName.GivenName)
	if err != nil {
		w.WriteHeader(500)
		w.Header().Set("Content-Type", "application/json")
		resp := make(map[string]string)
		resp["message"] = "B: Cannot set name - " + err.Error()
		jsonResp, _ := json.Marshal(resp)
		_, _ = w.Write(jsonResp)
		return
	}

	w.WriteHeader(204)
	return
}

func GetGlxLevelDataHandler(w http.ResponseWriter, r *http.Request) {
	// Check auth token etc
	valid := validateRequest(r)
	if !valid {
		log.Println("Err33: Validation failed")
		w.WriteHeader(403)
		w.Header().Set("Content-Type", "application/json")
		resp := make(map[string]string)
		resp["message"] = "Request did not pass validation."
		jsonResp, _ := json.Marshal(resp)
		_, _ = w.Write(jsonResp)
		return
	}

	// Get glx data from the request
	glx, err := getGalaxiatorData(r.Header.Get("x-glx-id"))
	if err != nil {
		log.Println("Err29: Cannot get Galaxiator data")
		w.WriteHeader(500)
		w.Header().Set("Content-Type", "application/json")
		resp := make(map[string]string)
		resp["message"] = "Cannot get Galaxiator data."
		jsonResp, _ := json.Marshal(resp)
		_, _ = w.Write(jsonResp)
		return
	}

	level, err := getLevelData(glx.Level)
	if err != nil {
		w.WriteHeader(400)
		w.Header().Set("Content-Type", "application/json")
		resp := make(map[string]string)
		resp["message"] = "Cannot get Level data."
		jsonResp, _ := json.Marshal(resp)
		_, _ = w.Write(jsonResp)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	jsonResp, _ := json.Marshal(level)
	_, _ = w.Write(jsonResp)
	return
}

func SetLevelHandler(w http.ResponseWriter, r *http.Request) {
	// Check auth token etc
	valid := validateRequest(r)
	if !valid {
		log.Println("Err34: Validation failed")
		w.WriteHeader(403)
		w.Header().Set("Content-Type", "application/json")
		resp := make(map[string]string)
		resp["message"] = "Request did not pass validation."
		jsonResp, _ := json.Marshal(resp)
		_, _ = w.Write(jsonResp)
		return
	}

	// Get givenname from the request and update the db
	// (lots of checking logic in the setGlxLevel func)
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println("Err30: Invalid request body")
		w.WriteHeader(400)
		w.Header().Set("Content-Type", "application/json")
		resp := make(map[string]string)
		resp["message"] = "Cannot read request body."
		jsonResp, _ := json.Marshal(resp)
		_, _ = w.Write(jsonResp)
		return
	}
	levelIn := levelJson{}
	err = json.Unmarshal(body, &levelIn)
	if err != nil {
		log.Println("Err38: Cannot unmarshall body ", body)
		w.WriteHeader(400)
		w.Header().Set("Content-Type", "application/json")
		resp := make(map[string]string)
		resp["message"] = "Cannot marshall level JSON."
		jsonResp, _ := json.Marshal(resp)
		_, _ = w.Write(jsonResp)
		return
	}
	level, err := setGlxLevel(r.Header.Get("x-glx-id"), levelIn.Level, r.Header.Get("x-glx-token"))
	if err != nil {
		w.WriteHeader(400)
		w.Header().Set("Content-Type", "application/json")
		resp := make(map[string]string)
		resp["message"] = "Cannot set Galaxiator level."
		jsonResp, _ := json.Marshal(resp)
		_, _ = w.Write(jsonResp)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	jsonResp, err := json.Marshal(level)
	_, _ = w.Write(jsonResp)
	return
}

func StoryHandler(w http.ResponseWriter, r *http.Request) {
	// Get glxId from request
	params := mux.Vars(r)
	glxId := params["glxid"]

	// Get the data for this Galaxiator
	glx, err := getGalaxiatorData(glxId)
	if err != nil {
		log.Println("Err39: Cannot get Galaxiator")
		w.WriteHeader(404)
		w.Header().Set("Content-Type", "application/json")
		resp := make(map[string]string)
		resp["message"] = "Cannot get Galaxiator."
		jsonResp, _ := json.Marshal(resp)
		_, _ = w.Write(jsonResp)
		return
	}

	// Get the full story text for that Galaxiator
	story, err := getWholeStory(glx.Level)
	if err != nil {
		log.Println("Err40: Cannot get story")
		w.WriteHeader(404)
		w.Header().Set("Content-Type", "application/json")
		resp := make(map[string]string)
		resp["message"] = "Cannot get story."
		jsonResp, _ := json.Marshal(resp)
		_, _ = w.Write(jsonResp)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	jsonResp, err := json.Marshal(map[string]string{
		"story":     story,
		"givenname": glx.GivenName,
		"nickname":  glx.NickName,
	})
	_, _ = w.Write(jsonResp)
	return
}

func ShareOnDiscordHandler(w http.ResponseWriter, r *http.Request) {
	// Check auth token etc
	valid := validateRequest(r)
	if !valid {
		log.Println("Err41: Validation failed")
		w.WriteHeader(403)
		w.Header().Set("Content-Type", "application/json")
		resp := make(map[string]string)
		resp["message"] = "Request did not pass validation."
		jsonResp, _ := json.Marshal(resp)
		_, _ = w.Write(jsonResp)
		return
	}

	// Get Galaxiator data
	glx, err := getGalaxiatorData(r.Header.Get("x-glx-id"))
	if err != nil {
		log.Println("Err44: Cannot get Galaxiator data.")
		w.WriteHeader(500)
		w.Header().Set("Content-Type", "application/json")
		resp := make(map[string]string)
		resp["message"] = "Cannot get Galaxiator data."
		jsonResp, _ := json.Marshal(resp)
		_, _ = w.Write(jsonResp)
		return
	}

	// Check has not already been shared
	if glx.SharedDiscord {
		log.Println("Err43: Already shared")
		w.WriteHeader(409)
		w.Header().Set("Content-Type", "application/json")
		resp := make(map[string]string)
		resp["message"] = "You've already shared this on Discord."
		jsonResp, _ := json.Marshal(resp)
		_, _ = w.Write(jsonResp)
		return
	}

	// Share on Discord
	message := fmt.Sprintf(`
:right_facing_fist: Welcome a new Galaxiator to the Arena! :left_facing_fist:

During the recruitment process, %s earned the nickname %s.

Read their story: https://names.galaxiators.com/#/story/%s
Or start your own story: https://names.galaxiators.com/
`, glx.GivenName, glx.NickName, r.Header.Get("x-glx-id"))
	_, err = discordClient.ChannelMessageSend(DiscordChannels["general"], message)
	if err != nil {
		log.Println("Err45: Error sending Discord message")
		w.WriteHeader(500)
		w.Header().Set("Content-Type", "application/json")
		resp := make(map[string]string)
		resp["message"] = "Error sending Discord message."
		jsonResp, _ := json.Marshal(resp)
		_, _ = w.Write(jsonResp)
		return
	}

	// Record this sharing so they only do it once
	_, err = db.Collection("glx").Doc(r.Header.Get("x-glx-id")).Set(ctx, map[string]bool{"disc": true}, firestore.MergeAll)
	if err != nil {
		log.Println("Err42: Error saving Discord status")
		w.WriteHeader(500)
		w.Header().Set("Content-Type", "application/json")
		resp := make(map[string]string)
		resp["message"] = "Error saving Discord status."
		jsonResp, _ := json.Marshal(resp)
		_, _ = w.Write(jsonResp)
		return
	}

	w.WriteHeader(204)
	return
}

func corsHandler(h http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "OPTIONS" {
			fmt.Println("Preflight detected: ", r.Header)
			w.Header().Add("Connection", "keep-alive")
			w.Header().Add("Access-Control-Allow-Origin", "*")
			w.Header().Add("Access-Control-Allow-Methods", "POST, OPTIONS")
			w.Header().Add("Access-Control-Allow-Headers", "Content-Type, Origin, Accept, x-glx-id, x-glx-token")
			w.Header().Add("Access-Control-Max-Age", "86400")
		} else {
			w.Header().Add("Access-Control-Allow-Origin", "*")
			h.ServeHTTP(w, r)
		}
	}
}
