package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
    "github.com/google/uuid"
	"github.com/gorilla/mux"
)

type Movie struct {
	ID       string    `json:"id"`
	Isbn     string    `json:"isbn"`
	Title    string    `json:"title"`
	Director *Director `json:"director"`
}

type Director struct {
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
}

var movies []Movie

func getMovies(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(movies)
}

func deleteMovies(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)

	for index, item := range movies {
		if item.ID == params["ID"] {
			movies = append(movies[:index], movies[index+1:]...)
			break
		}
	}
	json.NewEncoder(w).Encode(movies)
}

func getMovie(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    params := mux.Vars(r)

    for _, item := range movies {
        if item.ID == params["id"] {
            json.NewEncoder(w).Encode(item)
            return
        }
    }
}


func createMovies(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    var movie Movie
    err := json.NewDecoder(r.Body).Decode(&movie)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
    movie.ID = uuid.New().String()
    movies = append(movies, movie)
    json.NewEncoder(w).Encode(movie) // Renvoie uniquement le film créé
}


func updateMovies(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    params := mux.Vars(r)

    for index, item := range movies {
        if item.ID == params["id"] {
            var updatedMovie Movie
            err := json.NewDecoder(r.Body).Decode(&updatedMovie) // Ajoute un pointeur pour modifier directement le film décodé
            if err != nil {
                http.Error(w, err.Error(), http.StatusBadRequest)
                return
            }
            updatedMovie.ID = item.ID // Conserve l'ID d'origine
            movies[index] = updatedMovie
            json.NewEncoder(w).Encode(updatedMovie) // Renvoie uniquement le film mis à jour
            return
        }
    }

    http.NotFound(w, r)
}


func main() {
	r := mux.NewRouter()

	movies = append(movies, Movie{ID: "1", Isbn: "29383847", Title: "Movie one", Director: &Director{Firstname: "Mamadou", Lastname: "Diagne"}})
	movies = append(movies, Movie{ID: "2", Isbn: "28029398", Title: "Movie two", Director: &Director{Firstname: "Mamadou", Lastname: "Diagne"}})

	r.HandleFunc("/movies", getMovies).Methods("GET")
	r.HandleFunc("/movies/{id}", getMovie).Methods("GET")
	r.HandleFunc("/movies", createMovies).Methods("POST")
	r.HandleFunc("/movies/{id}", updateMovies).Methods("PUT")
	r.HandleFunc("/movies/{id}", deleteMovies).Methods("DELETE")

	fmt.Printf("Starting server at port 8080\n")

	log.Fatal(http.ListenAndServe(":8080", r))
}
