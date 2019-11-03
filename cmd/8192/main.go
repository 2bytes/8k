package main

import (
	"8192bytes/internal/flags"
	"8192bytes/pkg/backend"
	"8192bytes/pkg/storage"
	"8192bytes/pkg/storage/inmemory"
	"8192bytes/util"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"runtime"
	"time"
)

func main() {

	flag.Parse()

	if *flags.ShowVersion {
		fmt.Printf("Version %q built using %s.\n", flags.Version, runtime.Version())
		os.Exit(0)
	}

	if *flags.PublicAddress == "" {
		log.Println("public address cannot be empty")
		flag.Usage()
		os.Exit(1)
	}

	s := backend.NewServer(storage.NewMediator(inmemory.New(*flags.DataTTL, time.Second, *flags.MaxItemsStored)))

	http.HandleFunc("/", util.StatusCoder(s.HandleRequest))

	bindAddr := *flags.BindAddress

	if bindAddr == "" {
		bindAddr = "*"
	}

	var err error
	fmt.Printf("UI address set to: %s\n", s.UI.BaseAddress())

	if *flags.BindTLS {
		fmt.Printf("Binding: https://%s:%d\n", bindAddr, *flags.BindPort)
		listenAddr := fmt.Sprintf("%s:%d", *flags.BindAddress, *flags.BindPort)
		err = http.ListenAndServeTLS(listenAddr, *flags.TLSCertFile, *flags.TLSKeyFile, nil)
	} else {
		fmt.Printf("Binding: http://%s:%d\n", bindAddr, *flags.BindPort)
		listenAddr := fmt.Sprintf("%s:%d", *flags.BindAddress, *flags.BindPort)
		err = http.ListenAndServe(listenAddr, nil)
	}

	if err != nil {
		log.Println(err)
	}
}
