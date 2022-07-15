package main

// title Orders API
// @version 1.0
// @description This is a first service
// @termsOfService http://swagger.io/terms/
// @contact.name API Support
// @contact.email dbairamkulow@mail.ru
// @host localhost:8081
// @BasePath /

import (
	"encoding/json"
	_ "encoding/json"
	"fmt"
	_ "fmt"
	"github.com/gorilla/mux"
	_ "github.com/gorilla/mux"
	"io/ioutil"
	"log"
	_ "log"
	"net/http"
	_ "net/http"
	"net/http/httputil"
	"os"
	_ "os"
	"time"
	_ "time"
)

type Uniprot struct {
	Id            string `json:"Id"`
	Protein_name  string `json:"Protein_Name"`
	Protein_descr string `json:"Protein_descr"`
}

type Options struct {
	homePage            string
	returnAllProteins   string
	createNewProtein    string
	deleteProtein       string
	returnSingleProtein string
}

type Option func(*Options)

var uni []Uniprot

// returnAllArticles godoc
// @SAccept json
// @Router /all

func logHandler(fn func(w http.ResponseWriter, r *http.Request) Option) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		x, err := httputil.DumpRequest(r, true)
		if err != nil {
			http.Error(w, fmt.Sprint(err), http.StatusInternalServerError)
			return
		}
		select {
		case <-time.After(10 * time.Second):
			log.Println(fmt.Sprintf("%q", x))
			f, err := os.OpenFile("logs.txt", os.O_APPEND|os.O_WRONLY, 0644)
			if err != nil {
				fmt.Println(err)
				return
			}
			_, err = fmt.Fprintln(f, fmt.Sprintf("%q", x))
			if err != nil {
				fmt.Println(err)
			}

		}
	}
}

func returnAllProteins(w http.ResponseWriter, r *http.Request) Option {
	return func(args *Options) {

		w.WriteHeader(http.StatusOK)
		fmt.Println("Endpoint Hit: returnAllArticles")
		select {
		case <-time.After(10 * time.Second):
			fmt.Println("Save")

		}
		json.NewEncoder(w).Encode(uni)

	}
}

// homePage godoc
// @SAccept json
// @Router /

func homePage(w http.ResponseWriter, r *http.Request) Option {
	return func(args *Options) {
		fmt.Fprintf(w, "Welcom to the HomePage")
		fmt.Println("Endpoint Hit: homePage")
	}
}

// createNewArticle godoc
// @SAccept json
// @Router /article [post]
func createNewProtein(w http.ResponseWriter, r *http.Request) Option {
	return func(args *Options) {
		reqBody, _ := ioutil.ReadAll(r.Body)
		var seq Uniprot
		json.Unmarshal(reqBody, &seq)
		uni = append(uni, seq)
		json.NewEncoder(w).Encode(seq)
	}
}

// deleteArticle godoc
// @SAccept json
// @Router /article/{id} [delete]
func deleteProtein(w http.ResponseWriter, r *http.Request) Option {
	return func(args *Options) {
		vars := mux.Vars(r)
		id := vars["id"]

		for index, article := range uni {
			if article.Id == id {
				uni = append(uni[:index], uni[index+1:]...)
			}
		}
	}
}

// returnSingleArticle godoc
// @SAccept json
// @Router /article/{id}
func returnSingleProtein(w http.ResponseWriter, r *http.Request) Option {
	return func(args *Options) {
		vars := mux.Vars(r)
		key := vars["id"]
		for _, seq := range uni {
			if seq.Id == key {
				json.NewEncoder(w).Encode(seq)
			}
		}
	}
}

func handleRequests() {
	mainRouter := mux.NewRouter().StrictSlash(true)

	mainRouter.HandleFunc("/", logHandler(homePage))
	mainRouter.HandleFunc("/seq", logHandler(returnAllProteins))
	mainRouter.HandleFunc("/data/{id}", logHandler(returnSingleProtein))
	mainRouter.HandleFunc("/data", logHandler(createNewProtein)).Methods("POST")
	mainRouter.HandleFunc("/data/{id}", logHandler(deleteProtein)).Methods("Delete")

	log.Fatal(http.ListenAndServe(":8083", mainRouter))

}

func main() {
	fmt.Println("Rest API v2.0 - Mux Routers")
	uni = []Uniprot{
		Uniprot{Id: "1", Protein_name: "IgA", Protein_descr: ""},
		Uniprot{Id: "2", Protein_name: "BSA", Protein_descr: ""},
		Uniprot{Id: "3", Protein_name: "S-protein", Protein_descr: ""},
	}

	handleRequests()
}
