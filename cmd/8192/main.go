package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"runtime"
	"time"

	"github.com/2bytes/8k/internal/config"
	"github.com/2bytes/8k/pkg/backend"
	"github.com/2bytes/8k/pkg/storage"
	"github.com/2bytes/8k/pkg/storage/inmemory"
	"github.com/2bytes/8k/util"
)

func main() {

	flag.Parse()

	if *config.ShowVersion {
		fmt.Printf("Version %q built using %s.\n", config.Version, runtime.Version())
		os.Exit(0)
	}

	if *config.PublicAddress == "" {
		log.Println("public address cannot be empty")
		flag.Usage()
		os.Exit(1)
	}

	s := backend.NewServer(storage.NewMediator(inmemory.New(*config.DataTTL, time.Second, *config.MaxItemsStored)))

	http.HandleFunc("/", util.StatusCoder(s.HandleRequest))

	bindAddrDisplay := *config.BindAddress
	if bindAddrDisplay == "" {
		bindAddrDisplay = "*"
	}

	var err error
	fmt.Printf("UI address set to: %s\n", s.BaseAddress())

	if *config.BindTLS {
		fmt.Printf("Binding: https://%s:%d\n", bindAddrDisplay, *config.BindPort)
		listenAddr := fmt.Sprintf("%s:%d", *config.BindAddress, *config.BindPort)
		err = http.ListenAndServeTLS(listenAddr, *config.TLSCertFile, *config.TLSKeyFile, nil)
	} else {
		fmt.Printf("Binding: http://%s:%d\n", bindAddrDisplay, *config.BindPort)
		listenAddr := fmt.Sprintf("%s:%d", *config.BindAddress, *config.BindPort)
		err = http.ListenAndServe(listenAddr, nil)
	}

	if err != nil {
		log.Println(err)
	}
}
