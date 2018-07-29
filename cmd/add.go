// Copyright © 2018 Martijn Heuvelink <martijnheuvelink@hotmail.com>
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"io/ioutil"
	"log"
	"net/url"
	"os/exec"
	"regexp"
	"strings"
)

var Template = "<VirtualHost *:80> \n" +
	"	ServerName %subDomain%.%domain% \n" +
	"	ServerAlias www.%subDomain%.%domain% \n\n" +

	"	<Proxy *> \n" +
	"		Order allow,deny \n" +
	"		Allow from all \n" +
	"	</Proxy> \n" +
	"	ProxyPass / %redictURL% \n" +
	"	ProxyPassReverse / %redictURL% \n\n" +

	"	ErrorLog ${APACHE_LOG_DIR}/%subDomain%-error.log \n" +
	"	CustomLog ${APACHE_LOG_DIR}/%subDomain%-access.log combined \n" +
	"</VirtualHost>"

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:   "add",
	Short: "a",
	Run: func(cmd *cobra.Command, args []string) {
		checkArgs(args)
		formattedString := interpolateExpressions(args[0], args[1], args[2])
		checkApacheInstlled()
		writeVirtualHostFile(formattedString, args[0], args[1])
		enableSite(args[0], args[1])

		fmt.Printf("\n##########")
		fmt.Printf("\n# SUCCES #")
		fmt.Printf("\n##########\n")
	},
}

func checkArgs(args []string) {
	if len(args) != 3 {
		log.Fatalf("There are not the right amount arguments given. \n\nExpected 3 arguments: \n  * Domain\n  * subDomain\n  * redictURL")
	}
	if len(args[0]) < 3 {
		log.Fatalf("The first argument is to short. need to be at least 3 characters.")
	}
	if len(args[1]) < 3 {
		log.Fatalf("The second argument is to short. need to be at least 3 characters.")
	}
	if len(args[2]) < 3 {
		log.Fatalf("The third argument is to short. need to be at least 3 characters.")
	}
	validDomain(args[0])
	validSubdomain(args[1])
	validRedictURL(args[2])
}

func init() {
	rootCmd.AddCommand(addCmd)
}

func writeVirtualHostFile(formattedString string, domain string, subDomain string) {
	d1 := []byte(formattedString)
	location := "/etc/apache2/sites-available/" + subDomain + "." + domain + ".conf"
	err := ioutil.WriteFile(location, d1, 0644)
	checkError(err)
	fmt.Printf("wrote VirtualHost file '" + location + "'\n")
}

func interpolateExpressions(domain string, subDomain string, redictURL string) string {
	result := strings.Replace(Template, "%domain%", domain, -1)
	result = strings.Replace(result, "%redictURL%", redictURL, -1)
	result = strings.Replace(result, "%subDomain%", subDomain, -1)
	return result
}

func enableSite(domain string, subDomain string) {
	_, err := exec.Command("bash", "-c", "sudo a2ensite "+subDomain+"."+domain+".conf").Output()
	checkError(err)
	fmt.Printf("Enabled site '" + subDomain + "." + domain + "'\n")
}

func checkError(err error) {
	if err != nil {
		log.Fatal("error occured : %n", err)
	}
}

func checkApacheInstlled() {
	_, err := exec.Command("bash", "-c", "service --status-all | grep apache2").Output()
	if err != nil {
		fmt.Printf("Apache isnt installed yet!\n")
		out, err := exec.Command("bash", "-c", "apt install apache2 -y").Output()
		fmt.Printf(string(out))
		checkError(err)
	}
	fmt.Printf("Apache is installed.\n")
}

func validDomain(domain string) {
	r, err := regexp.Compile("^[a-zA-Z0-9][a-zA-Z0-9-]{1,61}[a-zA-Z0-9]\\.[a-zA-Z]{2,}$")
	checkError(err)
	if !r.MatchString(domain) {
		log.Fatal("'" + domain + "' isnt a valid domain\n")
	}
	fmt.Printf("'" + domain + "' is a valid domain\n")
}

func validSubdomain(subdomain string) {
	r, err := regexp.Compile("^[a-zA-Z0-9][a-zA-Z0-9-]{1,61}[a-zA-Z0-9]$")
	checkError(err)
	if !r.MatchString(subdomain) {
		log.Fatal("'" + subdomain + "' isnt a valid subdomain\n")
	}
	fmt.Printf("'" + subdomain + "' is a valid subdomain\n")
}

func validRedictURL(redictURL string) {
	_, err := url.ParseRequestURI(redictURL)
	if err != nil {
		log.Fatal("'" + redictURL + "' isnt a valid redictURL\n")
	}
	fmt.Printf("'" + redictURL + "' is a valid redictURL\n")
}
