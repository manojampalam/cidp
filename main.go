package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"fmt"
	"os"

	getopt "github.com/pborman/getopt/v2"
)

var (
	argsConfigFile = ""
	testMode       bool
	help           bool
)

func init() {
	getopt.FlagLong(&argsConfigFile, "config", 'c', "path to k8s config file")
	getopt.Flag(&testMode, 't', "run in test mode")
	getopt.Flag(&help, '?', "display help")
}

func main() {
	getopt.SetParameters("")
	getopt.Parse()

	if help == true {
		getopt.PrintUsage(os.Stdout)
		return
	}

	key, err := rsa.GenerateKey(rand.Reader, 2096)
	if err != nil {
		panic(err)
	}
	privateKeyBytes := x509.MarshalPKCS1PrivateKey(key)
	fmt.Println(base64.StdEncoding.EncodeToString(privateKeyBytes))
	_, err = x509.ParsePKCS1PrivateKey(privateKeyBytes)
	if err != nil {
		panic(err)
	}

	run(argsConfigFile)
}
