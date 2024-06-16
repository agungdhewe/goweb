package main

import (
	"flag"
	"fmt"
	"net/http"
	"path/filepath"
	"runtime"

	"github.com/fgtago/fgweb"
	"github.com/fgtago/fgweb/appsmodel"
	"github.com/fgtago/fgweb/defaulthandlers"
	"github.com/go-chi/chi/v5"
)

func main() {
	var err error

	fmt.Println("Starting Program ...")

	// baca parameter dari cli
	var cfgFileName string
	flag.StringVar(&cfgFileName, "conf", "config.yml", "nama file konfigurasi yang akan di load")
	flag.Parse()

	// set root direktori ke current working direktori
	// rootDir, err := os.Getwd()
	_, filename, _, _ := runtime.Caller(0)
	rootDir := filepath.Dir(filename)

	// ambil file konfigurasi
	cfgpath := filepath.Join(rootDir, cfgFileName)

	// start jalankan web
	ws, err := fgweb.New(rootDir, cfgpath)
	if err != nil {
		// ada error saat inisiasi webservice, halt
		panic(err.Error())
	}

	// info: memulai service
	port := ws.Configuration.Port
	fmt.Println("service running on port:", port)
	err = fgweb.StartService(port, func(mux *chi.Mux) error {
		return Router(mux)
	})
	if err != nil {
		// ada error saat service start, halt
		panic(err.Error())
	}

}

func Router(mux *chi.Mux) error {

	fgweb.Get(mux, "/favicon.ico", defaulthandlers.FaviconHandler)
	fgweb.Get(mux, "/asset/*", defaulthandlers.AssetHandler)
	fgweb.Get(mux, "/template/*", defaulthandlers.TemplateHandler)

	fgweb.Get(mux, "/", Home)
	fgweb.Get(mux, "/about", About)
	fgweb.Post(mux, "/subscribe", Subscribe)

	return nil
}

type User struct {
	Nama    string
	Alamat  string
	Message string
}

func Home(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	pv := ctx.Value(appsmodel.PageVariableKeyName).(*appsmodel.PageVariable)
	pv.PageName = "home"
	pv.Data = &User{
		Nama:   "Agung",
		Alamat: "Tangerang Raya",
	}
	defaulthandlers.SimplePageHandler(pv, w, r)
}

func About(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	pv := ctx.Value(appsmodel.PageVariableKeyName).(*appsmodel.PageVariable)
	pv.PageName = "about"
	defaulthandlers.SimplePageHandler(pv, w, r)
}

func Subscribe(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	pv := ctx.Value(appsmodel.PageVariableKeyName).(*appsmodel.PageVariable)
	pv.PageName = "home"
	data := &User{
		Nama:   "Agung",
		Alamat: "Tangerang Raya",
	}

	email := r.FormValue("email")
	if email != "" {
		data.Message = "Thank you for subscribing " + email
	}

	pv.Data = data
	defaulthandlers.SimplePageHandler(pv, w, r)
}
