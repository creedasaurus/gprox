package main

import (
	"crypto/tls"
	"embed"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"

	"github.com/jessevdk/go-flags"
	"github.com/rs/zerolog"
)

var (
	version         = "0.0.0"
	virtualCertPath = "cert/localhost.cert"
	virtualKeyPath  = "cert/localhost.key"
	savedCertName   = "gprox-localhost.cert"
	savedKeyName    = "gprox-localhost.key"
	log             = zerolog.New(os.Stdout).
			Output(zerolog.ConsoleWriter{Out: os.Stdout}).
			With().
			Timestamp().
			Logger()
)

var opts struct {
	Hostname  string `short:"n" long:"hostname" description:"The hostname to be used for the local proxy" default:"localhost"`
	Source    int    `short:"s" long:"source" description:"The source port that you will hit to go through the proxy" default:"9001"`
	Target    int    `short:"t" long:"target" description:"The port you are targeting" default:"9000"`
	Cert      string `short:"c" long:"cert" description:"Path to a .cert file"`
	Key       string `short:"k" long:"key" description:"Path to a .key file"`
	Config    string `short:"o" long:"config"`
	DropCerts bool   `short:"d" long:"dropcerts" description:"Save the built-in cert/key files to disk"`
	Version   bool   `long:"version"`
}

//go:embed cert/*
var certFS embed.FS

func main() {
	_, err := flags.Parse(&opts)
	if err != nil {
		parsedErr, ok := err.(*flags.Error)
		if !ok {
			log.Fatal().Err(err).Msg("error parsing flags")
			return
		}
		switch parsedErr.Type {
		case flags.ErrHelp, flags.ErrCommandRequired:
			return
		default:
			return
		}
	}

	if opts.Version {
		fmt.Printf("gprox %s\n", version)
		return
	}

	var certFile io.Reader
	if opts.Cert == "" {
		certFile, err = certFS.Open(virtualCertPath)
	} else {
		certFile, err = os.Open(opts.Cert)
	}
	if err != nil {
		log.Fatal().Err(err).Msg("failed to open cert file")
	}

	var keyFile io.Reader
	if opts.Key == "" {
		keyFile, err = certFS.Open(virtualKeyPath)
	} else {
		keyFile, err = os.Open(opts.Key)
	}
	if err != nil {
		log.Fatal().Err(err).Msg("failed to open key file")
	}

	if opts.DropCerts {
		outCert, err := os.OpenFile(savedCertName, os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			log.Fatal().Err(err).Msg("failed to open local cert file")
		}
		defer outCert.Close()

		_, err = io.Copy(outCert, certFile)
		if err != nil {
			log.Fatal().Err(err).Msg("failed to copy bytes to local cert file")
		}

		outKey, err := os.OpenFile(savedKeyName, os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			log.Fatal().Err(err).Msg("failed to open local key file")
		}

		_, err = io.Copy(outKey, keyFile)
		if err != nil {
			log.Fatal().Err(err).Msg("failed to copy bytes to local key file")
		}
		return
	}

	certBytes, err := ioutil.ReadAll(certFile)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to read bytes from cert file")
	}
	keyBytes, err := ioutil.ReadAll(keyFile)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to read bytes from key file")
	}
	certificate, err := tls.X509KeyPair(certBytes, keyBytes)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to create certificate")
	}

	origin, _ := url.Parse(fmt.Sprintf("http://%s:%d", opts.Hostname, opts.Target))
	reverseProxy := httputil.NewSingleHostReverseProxy(origin)

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Info().Str("method", r.Method).Msg("proxying")
		reverseProxy.ServeHTTP(w, r)
	})

	cfg := &tls.Config{Certificates: []tls.Certificate{certificate}}

	srv := &http.Server{
		TLSConfig: cfg,
		Addr:      fmt.Sprintf(":%d", opts.Source),
		Handler:   mux,
	}

	log.Info().
		Str("from", fmt.Sprintf("https://%s:%d", opts.Hostname, opts.Source)).
		Str("to", fmt.Sprintf("http://%s:%d", opts.Hostname, opts.Target)).
		Msg("Running proxy!")

	err = srv.ListenAndServeTLS("", "")
	if err != nil {
		log.Fatal().Err(err).Msg("proxy serve failure")
	}
}
