package main

import (
	"cloud.google.com/go/firestore"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
)

func addressFromSig(sigHex string, msg []byte) (addr string, err error) {
	sig := hexutil.MustDecode(sigHex)
	// https://github.com/ethereum/go-ethereum/blob/55599ee95d4151a2502465e0afc7c47bd1acba77/internal/ethapi/api.go#L442
	if sig[64] != 27 && sig[64] != 28 {
		return "", errors.New("incorrect sig length")
	}
	sig[64] -= 27

	pubKey, stperr := crypto.SigToPub(signHash(msg), sig)
	if stperr != nil {
		log.Println("Err1: ", stperr)
		return "", stperr
	}

	recoveredAddr := crypto.PubkeyToAddress(*pubKey)
	addr = strings.ToLower(recoveredAddr.String())
	return
}

func signHash(data []byte) []byte {
	msg := fmt.Sprintf("\x19Ethereum Signed Message:\n%d%s", len(data), data)
	return crypto.Keccak256([]byte(msg))
}

func substr(input string, length int) string {
	if length > len(input) {
		return input
	}
	asRunes := []rune(input)
	return string(asRunes[0:length])
}

func validateRequest(r *http.Request) bool {
	authToken := r.Header.Get("x-glx-token")
	glxId := r.Header.Get("x-glx-id")

	// Get user account from token
	dsnap1, err := db.Collection("tokens").Doc(authToken).Get(ctx)
	if err != nil {
		log.Println("Err2: ", err)
		return false
	}

	var tokenRecord tokenStruct
	dsnap1.DataTo(&tokenRecord)

	// Check token time is still valid, if not then delete it
	now := time.Now()
	sec := now.Unix()
	if tokenRecord.Expiry < sec {
		_, err = db.Collection("tokens").Doc(authToken).Delete(ctx)
		log.Println("Err3: ", err)
		return false
	}

	// Check IMX that this user owns that token id
	spaceClient := http.Client{Timeout: time.Second * 5}
	url := "https://api.x.immutable.com/v1/assets/0x6c82e53cbbd8a6afaf9663d58547cfc1a43be7aa/" + glxId
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Println("Err4: ", err)
		return false
	}
	req.Header.Set("User-Agent", "galaxiators-namegame")
	res, getErr := spaceClient.Do(req)
	if getErr != nil {
		log.Println("Err5: ", err)
		return false
	}
	if res.Body != nil {
		defer res.Body.Close()
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Println("Err6: ", err)
		return false
	}

	imxGlx := IMXGalaxiator{}
	err = json.Unmarshal(body, &imxGlx)
	if err != nil {
		log.Println("Err7: ", err)
		return false
	}
	if imxGlx.User != tokenRecord.User {
		return false
	}

	return true
}

func setGivenName(glx string, name string) error {
	// Check for existing tokens and delete
	iter := db.Collection("glx").Where("givn", "==", name).Documents(ctx)
	docArray, err := iter.GetAll()
	if err != nil {
		log.Println("Err46 error checking name duplicates: ", err)
		return err
	}
	if len(docArray) > 0 {
		log.Println("Err47: name already exists")
		return errors.New("name already exists")
	}

	_, err = db.Collection("glx").Doc(glx).Set(ctx, map[string]interface{}{"givn": name, "disc": false}, firestore.MergeAll)
	if err != nil {
		log.Println("Err8: ", err)
		return err
	}
	return nil
}

func setGlxLevel(glxId string, newLevelId string, token string) (levelData Level, err error) {
	var glx Galaxiator
	var oldLevelId string
	var oldLevel Level

	// If newLevel = "" then we're initializing the level game.
	// Give them the starting one for their race.
	if newLevelId == "" {
		// Get race from IMX...
		spaceClient := http.Client{Timeout: time.Second * 5}
		url := "https://api.x.immutable.com/v1/assets/0x6c82e53cbbd8a6afaf9663d58547cfc1a43be7aa/" + glxId
		req, err1 := http.NewRequest(http.MethodGet, url, nil)
		if err1 != nil {
			log.Println("Err9: ", err1)
			return Level{}, err1
		}
		req.Header.Set("User-Agent", "galaxiators-namegame")
		res, err2 := spaceClient.Do(req)
		if err2 != nil {
			log.Println("Err10: ", err2)
			return Level{}, err2
		}
		if res.Body != nil {
			defer res.Body.Close()
		}
		body, err3 := ioutil.ReadAll(res.Body)
		if err3 != nil {
			log.Println("Err11: ", err3)
			return Level{}, err3
		}

		imxGlx := IMXGalaxiator{}
		err = json.Unmarshal(body, &imxGlx)
		if err != nil {
			log.Println("Err12: ", err)
			return Level{}, err
		}

		// ...and use it to set the level
		newLevelId = strings.ToLower(imxGlx.Name[0:1])
	} else {
		// Get the old level...
		glx, err = getGalaxiatorData(glxId)
		if err != nil {
			log.Println("Err13: ", err)
			return
		}
		oldLevelId = glx.Level

		// Make sure the newLevel one is a progression forwards...
		if len(newLevelId) <= len(oldLevelId) {
			err = errors.New("invalid new level: cannot progress backwards")
			log.Println("Err14: ", err)
			return
		}

		// ...and check it is in the available options for answers
		oldLevel, err = getLevelData(oldLevelId)
		if err != nil {
			log.Println("Err15: ", err)
			return
		}
		invalidAnswer := true
		for _, b := range oldLevel.Answers {
			if b.NextLevel == newLevelId {
				invalidAnswer = false
			}
		}
		if invalidAnswer {
			err = errors.New("invalid new level: not an option from previous level")
			log.Println("Err16: ", err)
			return
		}
	}

	// Update Galaxiator level
	_, err = db.Collection("glx").Doc(glxId).Set(ctx, map[string]string{"levl": newLevelId}, firestore.MergeAll)
	if err != nil {
		log.Println("Err17: ", err)
		return
	}

	// Get the new level
	dsnap, err := db.Collection("story").Doc(newLevelId).Get(ctx)
	if err != nil {
		log.Println("Err18: ", err)
		return
	}
	dsnap.DataTo(&levelData)
	levelData.ID = newLevelId

	if levelData.NameEarned != "" {
		// This is a final level, so set the Galaxiator earned nickname...
		_, err = db.Collection("glx").Doc(glxId).Set(ctx, map[string]string{"nick": levelData.NameEarned}, firestore.MergeAll)
		if err != nil {
			log.Println("Err19: ", err)
			return
		}

		//...and add the Payout to the bank transactions
		tokenData := tokenStruct{}

		dsnap, err = db.Collection("tokens").Doc(token).Get(ctx)
		if err != nil {
			log.Println("Err25: ", err)
			return
		}
		dsnap.DataTo(&tokenData)

		_, _, err = db.Collection("banktxns").Add(ctx, map[string]interface{}{
			"in":   levelData.CyberSand,
			"desc": "Initial Galaxiator recruitment stipend for #" + glxId,
			"ts":   firestore.ServerTimestamp,
			"user": tokenData.User,
		})
		if err != nil {
			log.Println("Err20: ", err)
			return
		}
	}
	return
}

func getGalaxiatorData(glxId string) (glx Galaxiator, err error) {
	dsnap1, err := db.Collection("glx").Doc(glxId).Get(ctx)
	if err != nil {
		log.Println("Err21: ", err)
		return
	}
	err = dsnap1.DataTo(&glx)
	if err != nil {
		log.Println("Err22: ", err)
		return
	}
	return
}

func getLevelData(levelId string) (level Level, err error) {
	dsnap1, err := db.Collection("story").Doc(levelId).Get(ctx)
	if err != nil {
		log.Println("Err23: ", err)
		return
	}
	dsnap1.DataTo(&level)
	if err != nil {
		log.Println("Err24: ", err)
		return
	}
	level.ID = levelId

	return
}

func getWholeStory(levelString string) (story string, err error) {
	var subLevelString string
	for i := len(levelString); i > 0; i-- {
		subLevelString = substr(levelString, i)
		dsnap2, err := db.Collection("story").Doc(subLevelString).Get(ctx)
		if err != nil {
			return "", err
		}
		subLevel := Level{}
		dsnap2.DataTo(&subLevel)
		if subLevel.Outcome == "" {
			nextLevelString := substr(levelString, i+1)
			var answer string
			for _, element := range subLevel.Answers {
				if element.NextLevel == nextLevelString {
					answer = element.Text
					break
				}
			}
			story = subLevel.Text + " \n\n" + answer + " \n\n" + story
		} else {
			story = subLevel.Outcome /*+
			"\n\nYou have been awarded an initial Galaxiator recruitment stipend of $GLXR " + strconv.Itoa(subLevel.CyberSand) + " CyberSand."*/
		}
	}
	return
}
