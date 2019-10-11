package controller

import (
	"fmt"
	"log"
	"encoding/json"
	"context"
	"time"
	"net/http"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"github.com/gorilla/mux"

	"github.com/dymzfp/base-golang-mongo/model"
	conf "github.com/dymzfp/base-golang-mongo/config"
)

const (
	errorConnectDB  = "unable to connect database"
	errorDecodeJson = "decoding json error"
	errorDataEmpty  = "data empty"
	errorInvalidID  = "Invalid id"
)

func PostData(w http.ResponseWriter, r *http.Request) {
	db, err := conf.Connect()
	if err != nil {
		log.Println(errorConnectDB, err.Error())
		return
	}

	resp := model.NewResponseFormat()

	var data model.User
	err = json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		resp.AddError(errorDecodeJson, err.Error())
		sendResponse(http.StatusInternalServerError, resp, w, r)
		return
	}

	// set collection
	collection := db.Collection("user")
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	result, _ := collection.InsertOne(ctx, data)

	resp.SetData(result)
	sendResponse(http.StatusCreated, resp, w, r)
	return
}

func GetData(w http.ResponseWriter, r *http.Request) {
	db, err := conf.Connect()
	if err != nil {
		log.Println(errorConnectDB, err.Error())
		return
	}

	resp := model.NewResponseFormat()

	var data []model.User

	collection := db.Collection("user")
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	cur, err := collection.Find(ctx, bson.M{})
	if err != nil { 
		resp.AddError("error", err.Error())
	   	sendResponse(http.StatusInternalServerError, resp, w, r)
	   	return
	}
	defer cur.Close(ctx)
	
	for cur.Next(ctx) {
		var result model.User
		err := cur.Decode(&result)
		if err != nil { 
	   		resp.AddError("error", err.Error())
	   		sendResponse(http.StatusInternalServerError, resp, w, r)
	   		return
	   	}

	   	data = append(data, result) 
	}

	if err := cur.Err(); err != nil {
	   resp.AddError("error", err.Error())
	   sendResponse(http.StatusInternalServerError, resp, w, r)
	   return
	}

	if len(data) <= 0 {
		resp.AddError(errorDataEmpty, "")
		sendResponse(http.StatusNotFound, resp, w, r)
		return
	}

	resp.SetData(data)
	sendResponse(http.StatusOK, resp, w, r)
	return
}

func GetDataSingle(w http.ResponseWriter, r *http.Request) {
	db, err := conf.Connect()
	if err != nil {
		log.Println(errorConnectDB, err.Error())
		return
	}

	resp := model.NewResponseFormat()

	params := mux.Vars(r)
	id, e := primitive.ObjectIDFromHex(params["id"])
	if e != nil {
		resp.AddError(errorInvalidID , e.Error())
		sendResponse(http.StatusBadRequest, resp, w, r)
		return
	}

	var data model.User

	collection := db.Collection("user")
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	err = collection.FindOne(ctx, model.User{ID: id}).Decode(&data)	
	if err != nil {
    	resp.AddError("Id invalid", err.Error())
		sendResponse(http.StatusInternalServerError, resp, w, r)
		return
	}

	resp.SetData(data)
	sendResponse(http.StatusOK, resp, w, r)
	return
}

// response
func sendResponse(statusCode int, resp *model.ResponseFormat, w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("content-type", "application/json")
	w.WriteHeader(statusCode)
	
	encodedResponse, err := resp.EncodeToJSON()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(nil)
		log.Printf("Source: %v| Destination: %v| ResponseCode: %v| ResponseLen: %v", r.RemoteAddr, r.RequestURI, statusCode, "error while encoding response")
		return fmt.Errorf("unable to encode JSON: %v", err)
	}
	w.Write(encodedResponse)
	if user := r.Header.Get("user"); user != "" {
		log.Printf("| User: %v | Source: %v | Destination: %v | Mehod: %v | ResponseCode: %v | ResponseLen: %v", user, r.RemoteAddr, r.RequestURI, r.Method, statusCode, len(encodedResponse))
	} else {
		log.Printf("| Source: %v | Destination: %v | Mehod: %v | ResponseCode: %v | ResponseLen: %v", r.RemoteAddr, r.RequestURI, r.Method, statusCode, len(encodedResponse))
	}
	return nil
}