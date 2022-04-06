package main

import (
	"flag"
	"log"

	homedir "github.com/mitchellh/go-homedir"

	"github.com/blackironj/go-solana-sdk-test/hello"
)

var (
	keypairPath  string
	programIdB58 string
)

func init() {
	flag.StringVar(&keypairPath, "keypair", "", "keypair path (required)")
	flag.StringVar(&programIdB58, "pid", "", "base58 program id (required)")
	flag.Parse()
}

func main() {
	if keypairPath == "" || programIdB58 == "" {
		log.Fatalf("keypair and program id must be required")
	}

	home, _ := homedir.Dir()
	account := hello.GetAccountFromFile(home + "/" + keypairPath)

	helloWorld := hello.GetHelloWorld(account, programIdB58)
	helloWorld.SayHello()

	helloWorld.ReportGreetings()
}
