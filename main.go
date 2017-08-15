package main

import (
	"net/http"
	"github.com/gorilla/mux"
	"log"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudfront"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"os"
	"fmt"
	"html/template"
	"strings"
	"github.com/prometheus/client_golang/prometheus"
)

var awsCredentials = credentials.NewStaticCredentials(os.Getenv("aws_access_key_id"), os.Getenv("aws_secret_access_key"), "")

type ENV struct {
	CDN []CDN
}

type CDN struct {
	Head string
	Links
}

type Links struct {
	Path []string
}

const PROD = os.Getenv("PRODID")
const STAGE = os.Getenv("STAGEID")
const STAGE2 = os.Getenv("DEVID")

func showList(w http.ResponseWriter, r *http.Request) {
	sess, _ := session.NewSession()

	svc := cloudfront.New(sess, &aws.Config{Region: aws.String("us-east-1"), Credentials: awsCredentials})

	input := &cloudfront.GetDistributionConfigInput{}

	var distribIds = [3]string {STAGE2, STAGE, PROD}

	t := template.Must(template.ParseFiles("/opt/templates/index.tmpl"))
	var env ENV
	cdnArray := make([]CDN, int(len(distribIds)))

	for j, id := range distribIds {
		input.SetId(id)

		result, err := svc.GetDistributionConfig(input)

		if err != nil {
			if aerr, ok := err.(awserr.Error); ok {
				switch aerr.Code() {
				case cloudfront.ErrCodeNoSuchDistribution:
					fmt.Println(cloudfront.ErrCodeNoSuchDistribution, aerr.Error())
				case cloudfront.ErrCodeAccessDenied:
					fmt.Println(cloudfront.ErrCodeAccessDenied, aerr.Error())
				default:
					fmt.Println(aerr.Error())
				}
			} else {
				fmt.Println(err.Error())
			}
		}

		linksArray := make([]string, int64(*result.DistributionConfig.CacheBehaviors.Quantity))

		for i, links := range result.DistributionConfig.CacheBehaviors.Items {
			str := *links.TargetOriginId + " " + *links.PathPattern
			str = strings.Replace(str, "Custom-", "", -1)
			linksArray[i] = str
		}

		cdnArray[j] = CDN{Head: *result.DistributionConfig.Aliases.Items[0], Links: Links{Path: linksArray}}
	}
	env = ENV{CDN:cdnArray}

	w.Header().Set(
		"Content-Type",
		"text/html",
	)

	t.Execute(w, env)
}

func serveResource(w http.ResponseWriter, req *http.Request) {
	path := "/opt/." + req.URL.Path
	http.ServeFile(w, req, path)
}

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/", showList)
	router.HandleFunc("/templates/favicon.ico", serveResource)
	router.HandleFunc("/templates/bootstrap.min.css", serveResource)
	router.Handle("/metrics", prometheus.UninstrumentedHandler())

	http.Handle("/", router)
	log.Println(http.ListenAndServe(":8080", router))
}