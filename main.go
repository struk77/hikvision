package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"gopkg.in/alecthomas/kingpin.v2"
	"gopkg.in/yaml.v2"
)

// var opts struct {
// 	Cameras string `short:"c" long:"cameras" description:"List cameras file" value-name:"PATH" default:"/etc/hikvision/cameras.yml"`
// 	Listen  string `short:"l" long:"listen" description:"Listen address" value-name:"[HOST]:PORT" default:":19101"`
// 	Period  uint   `short:"p" long:"period" description:"Period in seconds, should match Prometheus scrape interval" value-name:"SECS" default:"60"`
// }

var config struct {
	Cameras string `short:"c" long:"cameras" description:"List cameras file" value-name:"PATH" default:"cameras.yml"`
	Listen  string `short:"l" long:"listen" description:"Listen address" value-name:"[HOST]:PORT" default:":19101"`
	Period  uint   `short:"p" long:"period" description:"Period in seconds, should match Prometheus scrape interval" value-name:"SECS" default:"60"`
}

//Camera is a camera
type Camera struct {
	Host     string `yaml:"host"`
	IP       string `yaml:"ip"`
	Login    string `yaml:"login"`
	Password string `yaml:"password"`
}

type ymlCameras struct {
	Cameras []Camera `yaml:"cameras"`
}

func getCameras(filePath string) (map[string]Camera, error) {
	yamlFile, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer yamlFile.Close()
	jsonData, err := ioutil.ReadAll(yamlFile)
	if err != nil {
		return nil, err
	}
	var jCameras ymlCameras
	err = yaml.Unmarshal(jsonData, &jCameras)
	if err != nil {
		return nil, err
	}
	cameras := make(map[string]Camera, len(jCameras.Cameras))
	for _, el := range jCameras.Cameras {
		cameras[el.IP] = el
	}
	return cameras, nil
}

func probeHandler(w http.ResponseWriter, r *http.Request) {
	host := r.URL.Query().Get("target")
	if host == "" {
		w.Header().Set("Content-Type", "text/html")
		_, _ = w.Write([]byte(`<html>
		    <head><title>Hikvision Exporter</title></head>
			<body>
			<b>ERROR: missing target parameter</b>
			</body>`))
		return
	}
	cameras, err := getCameras(config.Cameras)
	if err != nil {
		log.Fatal("cameras, err := getCameras(opts.Cameras) ", err)
	}
	camera := cameras[host]
	if err != nil {
		log.Fatal("camera := cameras[host]", err)
	}
	target := GetTarget(
		WorkerSpec{
			login:    camera.Login,
			password: camera.Password,
			host:     camera.IP,
		},
	)

	h := promhttp.HandlerFor(target.registry, promhttp.HandlerOpts{})
	h.ServeHTTP(w, r)
}

func main() {
	godotenv.Load()

	a := kingpin.New("hikvision_exporter", "Hikvision Cameras Exporter for Prometheus")
	a.HelpFlag.Short('h')

	a.Flag("cameras", "path to the file where storing cameras specs").
		Envar("CAMERAS").
		Default("cameras.yml").
		StringVar(&config.Cameras)

	a.Flag("listen", "The address the hikvision_exporter listens on for incoming webhooks").
		Envar("LISTEN").
		Default(":19101").
		StringVar(&config.Listen)

	a.Flag("period", "Period in seconds, should match Prometheus scrape interval").
		Envar("PERIOD").
		Default("60").
		UintVar(&config.Period)

	// if _, err := flags.Parse(&opts); err != nil {
	// 	os.Exit(0)
	// }
	_, err := a.Parse(os.Args[1:])
	if err != nil {
		log.Printf("error parsing commandline arguments: %v\n", err)
		a.Usage(os.Args[1:])
		os.Exit(2)
	}
	http.Handle("/metrics", promhttp.Handler())
	http.HandleFunc("/probe", probeHandler)
	log.Fatal(http.ListenAndServe(config.Listen, nil))
}
