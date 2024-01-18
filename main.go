package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
	"github.com/schollz/progressbar/v3"
	"log"
	"net/url"
	"os"
	"regexp"
	"strings"
)

var (
	blockSize       = int64(4 * 1024 * 1024) // 4 MB
	concurrency     = uint16(8)
	pathPrefixRegex = regexp.MustCompile(`[a-zA-Z0-9/]`)
)

func upload(client *azblob.Client, container string, pathPrefix string, path string) (*azblob.UploadFileResponse, error) {
	fh, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	stat, err := os.Stat(path)
	if err != nil {
		return nil, err
	}

	fileSize := stat.Size()

	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Printf("")
			log.Panic(err)
		}
	}(fh)

	pb := progressbar.DefaultBytes(fileSize, "uploading")

	resp, err := client.UploadFile(context.TODO(), container, pathPrefix+path, fh, &azblob.UploadFileOptions{
		BlockSize:   blockSize,
		Concurrency: concurrency,
		Progress: func(bytesTransferred int64) {
			err := pb.Set64(bytesTransferred)
			if err != nil {
				log.Printf("error setting progress bar: %v", err)
			}
		},
	})
	if err != nil {
		return nil, err
	}

	return &resp, nil
}

func main() {
	var account string
	var container string
	var pathPrefix string

	flag.StringVar(&account, "a", "", "account")
	flag.StringVar(&container, "b", "", "bucket/container")
	flag.StringVar(&pathPrefix, "p", "", "path prefix")

	flag.Usage = func() {
		_, err := fmt.Fprintf(os.Stderr, "Usage: %s [options] <file1> <file2>...\n", os.Args[0])
		if err != nil {
			log.Fatal(err)
		}
		flag.PrintDefaults()
		_, err = fmt.Fprintf(os.Stderr, "\nExample: %s -a myaccount -b mycontainer file1 file1.md5sum\n", os.Args[0])
		if err != nil {
			log.Fatal(err)
		}
	}

	flag.Parse()
	values := flag.Args()

	if len(values) == 0 || len(strings.TrimSpace(account)) == 0 || len(strings.TrimSpace(container)) == 0 ||
		(len(pathPrefix) > 0 && len(strings.TrimSpace(pathPrefix)) == 0) {
		flag.Usage()
		os.Exit(1)
	}

	if len(pathPrefix) > 0 && strings.HasPrefix(pathPrefix, "/") {
		pathPrefix = strings.TrimPrefix(pathPrefix, "/")
	}
	if len(pathPrefix) > 0 && !strings.HasSuffix(pathPrefix, "/") {
		pathPrefix = pathPrefix + "/"
	}
	if len(pathPrefix) > 0 && len(pathPrefixRegex.ReplaceAllString(pathPrefix, "")) > 0 {
		log.Fatal("path prefix can only contain alphanumeric characters and path separator (/)")
	}

	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatal(err)
	}

	serviceURL := fmt.Sprintf("https://%s.blob.core.windows.net/", account)
	client, err := azblob.NewClient(serviceURL, cred, nil)
	if err != nil {
		log.Fatal(err)
	}

	for _, path := range values {
		_, err := upload(client, container, pathPrefix, path)
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("%s%s/%s", serviceURL, container, pathPrefix+url.PathEscape(path))
	}
}
