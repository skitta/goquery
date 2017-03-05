// Copyright © 2017 Skitta Chen <skittachen@gmail.com>
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
	"net/http"
	"strings"

	"github.com/skitta/goquery/api/express"
	"github.com/spf13/cobra"
)

var (
	expressNum string
)

// expressCmd represents the express command
var expressCmd = &cobra.Command{
	Use:   "express",
	Short: "Track your express",
	Long: `Track your express by a given package number, or scan your package database and track each of them.

Run this command without flag will scan database by default.`,
	Run: trackRun,
}

func init() {
	RootCmd.AddCommand(expressCmd)

	expressCmd.Flags().StringVarP(&expressNum, "number", "n", "", "express number")
}

func trackRun(cmd *cobra.Command, args []string) {
	if expressNum == "" {
		trackByScanDB()
	} else {
		trackByNum()
	}
}

func trackByNum() {
	p := express.Track(expressNum)
	if p.Value == "" {
		msg := "No Such A Package!"
		postToHubot("http://127.0.0.1:8080/hubot/webhoods/express", msg)
		fmt.Println(msg)
	}
	postToHubot("http://127.0.0.1:8080/hubot/webhoods/express", p.Value)
	fmt.Println(p.Value)
}

func trackByScanDB() {
	allpackages := express.GetAllPackages()

	if len(allpackages) == 0 {
		fmt.Println("No package in your database.")
		return
	}

	for i := 0; i < len(allpackages); i++ {
		p := express.Track(allpackages[i].Key)
		if p.Value != allpackages[i].Value {
			fmt.Println(p.Value)
			postToHubot("http://127.0.0.1:8080/hubot/webhoods/express", p.Value)
			p.Update()
		}
		if p.Checked() {
			s := fmt.Sprintf("package %s is checked.", allpackages[i].Key)
			fmt.Println(s)
			p.Del()
		}
	}
}

func postToHubot(URL, msg string) {
	data := fmt.Sprintf("{\"message\": \"%s\"}", msg)
	payload := strings.NewReader(data)
	http.Post(URL, "application/json", payload)
}
