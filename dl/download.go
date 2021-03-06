package dl

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"path"
	"strings"
	"syscall"

	pb "gopkg.in/cheggaaa/pb.v1"

	"golang.org/x/crypto/ssh/terminal"

	"github.com/PuerkitoBio/goquery"
	"github.com/tschf/odl/resource"
)

func checkRedirect(req *http.Request, via []*http.Request) error {
	req.Header.Add("User-Agent", "Mozilla/5.0")

	return nil
}

func SaveResource(res *resource.OracleResource, otnUser string, otnPassword string, skipExisting bool) {
	var cookies []*http.Cookie
	cookies = append(cookies, res.AcceptCookie)

	cookieJar, _ := cookiejar.New(nil)

	client := &http.Client{
		// Need to re-add user agent as Go doesn't propagate them
		// see: https://groups.google.com/forum/#!topic/golang-nuts/OwGvopYXpwE
		// and https://github.com/golang/go/issues/4800
		// slated to be fixed in golang 1.8. Done in checkRedirect
		CheckRedirect: checkRedirect,
		Jar:           cookieJar,
	}

	// flag to identify if we need to submit username/password data
	// to https://login.oracle.com/oam/server/sso/auth_cred_submit. Specifically
	// some resources include more than one file (such as Oracle 12c EE).
	// In which case, we only need to submit once.
	// (todo/test: read response headers to see if auth is required?)
	requiresAuth := true

	for _, file := range res.File {
		//Check to see if the file we are requesting to download already exists.
		//Impl borrowed from: http://stackoverflow.com/questions/12518876/how-to-check-if-a-file-exists-in-go
		//err will be nil when getting information about a file that exists
		_, fileStatErr := os.Stat(path.Base(file))
		fileExists := fileStatErr == nil

		if fileExists && skipExisting {
			continue
		} else if fileExists {
			// The skip existing flag was unset (false), so prompt user at run time
			fmt.Printf("This file already exists. Would you like to overwrite %s?\n", path.Base(file))
			fmt.Print("Enter Y to overwrite, or N to skip: ")
			reader := bufio.NewReader(os.Stdin)
			strOverwriteFile, _ := reader.ReadString('\n')
			if strings.TrimSpace(strOverwriteFile) != "Y" {
				continue
			}
		}

		req, _ := http.NewRequest("GET", file, nil)
		req.Header.Add("User-Agent", "Mozilla/5.0")

		u, _ := url.Parse(file)

		//Set initial jar with a cookie accepting the license agreement
		//Example taken from: https://gist.github.com/Rabbit52/a8a44c3c4cd514052952
		cookieJar.SetCookies(u, cookies)

		resp, respErr := client.Do(req)
		if respErr != nil {
			fmt.Println("Couldn't read response")
			log.Fatal(respErr)
		}
		defer resp.Body.Close()

		if !res.SkipAuth && requiresAuth {

			if len(otnPassword) == 0 {
				fmt.Printf("To complete the license acceptance, you must enter valid OTN credentials. Please enter your OTN password (%s):", otnUser)
				consolePass, _ := terminal.ReadPassword(int(syscall.Stdin))
				fmt.Println()
				otnPassword = string(consolePass)
			}

			doc, err := goquery.NewDocumentFromResponse(resp)
			if err != nil {
				log.Fatal(err)
			}

			// https://godoc.org/github.com/PuerkitoBio/goquery
			pageForms := doc.Find("form")

			//POST example: http://stackoverflow.com/questions/19253469/make-a-url-encoded-post-request-using-http-newrequest
			authData := url.Values{}
			pageForms.Find("input").Each(func(index int, el *goquery.Selection) {
				inputName, _ := el.Attr("name")
				inputValue, _ := el.Attr("value")

				authData.Set(inputName, inputValue)
			})

			authData.Set("username", otnUser)
			authData.Set("password", otnPassword)

			req, _ = http.NewRequest("POST", "https://login.oracle.com/oam/server/sso/auth_cred_submit", bytes.NewBufferString(authData.Encode()))
			req.Header.Add("User-Agent", "Mozilla/5.0")
			req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

			resp, _ = client.Do(req)
			requiresAuth = false
		}

		savedFile, err := os.Create(path.Base(file))
		defer savedFile.Close()
		if err != nil {
			log.Fatal(err)
		}

		requestTotalSize := resp.ContentLength
		//Set up progress bar
		progressBar := pb.New64(requestTotalSize)
		progressBar.SetUnits(pb.U_BYTES)
		progressBar.Prefix(fmt.Sprintf("%s:", path.Base(file)))
		progressBar.Start()

		readerWithProgress := progressBar.NewProxyReader(resp.Body)

		_, err = io.Copy(savedFile, readerWithProgress)
		if err != nil {
			log.Fatal(err)
		}
		progressBar.Finish()
	}

}
