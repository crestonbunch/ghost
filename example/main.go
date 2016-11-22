package main

import (
	"context"
	"encoding/json"
	"errors"
	. "github.com/crestonbunch/ghost"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"
)

type SampleExtender struct {
	secret string
}

func (e *SampleExtender) Extend(ctx context.Context) context.Context {
	return context.WithValue(ctx, "secret", e.secret)
}

type User struct {
	Id     int
	Name   string
	Email  string
	Secret string
}

type UserIdGet struct {
	id     int
	secret string
}

// Extracts the :id from the URL and puts it in a model
func (u *UserIdGet) FromRequest(r *http.Request) HttpError {
	vars := mux.Vars(r)
	id := vars["id"]

	if i, err := strconv.Atoi(id); err != nil {
		return NewHttpError(err, 400)
	} else {
		u.id = i
	}

	u.secret = r.Context().Value("secret").(string)

	return nil
}

// Validates that the integer is a positive number
type UserIdGetValidator struct {
}

func (v *UserIdGetValidator) Validate(model interface{}) HttpError {
	id := model.(*UserIdGet).id

	if id < 0 {
		return NewHttpError(errors.New("Please enter a valid user id."), 400)
	}

	return nil
}

// Get the user from the user id
type UserIdGetProcessor struct {
}

func (p *UserIdGetProcessor) Process(model interface{}) (interface{}, HttpError) {
	id := model.(*UserIdGet).id
	secret := model.(*UserIdGet).secret

	if id == 1 {
		return User{
			Id:     1,
			Name:   "Joe",
			Email:  "joe@example.com",
			Secret: secret,
		}, nil
	} else {
		return nil, NewHttpError(errors.New("User not found!"), 404)
	}
}

type UserIdGetWriter struct {
}

func (w *UserIdGetWriter) Serialize(output interface{}) ([]byte, HttpError) {
	if result, err := json.Marshal(output); err != nil {
		return nil, NewHttpError(err, 500)
	} else {
		return result, nil
	}
}

func buildRouter() *Router {

	router := NewRouter()

	sampleExtender := &SampleExtender{"s3cr3tstr!ng"}
	sampleModel := &UserIdGet{}
	sampleValidator := &UserIdGetValidator{}
	sampleProcessor := &UserIdGetProcessor{}
	sampleWriter := &UserIdGetWriter{}

	router.AddRoute("/user/id/{id}").
		Methods("GET").
		Extender(sampleExtender).
		Model(sampleModel).
		Validator(sampleValidator).
		Processor(sampleProcessor).
		Writer(sampleWriter)

	return router
}

func main() {
	router := buildRouter()

	http.Handle("/", router.Mux)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
