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

	"github.com/skitta/goquery/api/express"
	"github.com/skitta/goquery/api/slack"
	"github.com/spf13/cobra"
)

var expressNum string

const expressHook = "https://hooks.slack.com/services/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"

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

	expressCmd.Flags().StringVarP(&expressNum, "add", "a", "", "express number")
}

func trackRun(cmd *cobra.Command, args []string) {
	if expressNum == "" {
		trackByScanDB()
	} else {
		addExpressNum()
	}
}

func addExpressNum() {
	p := express.Track(expressNum)
	if p.Value != "" {
		m := slack.NewMessage(fmt.Sprintf("Status of express %s", p.ID), "")
		m.AddField("Express Status", p.Value, false)
		slack.Post(expressHook, m)
		p.Update()
	}
}

func trackByScanDB() {
	ps := express.GetAllPackages()

	for i := 0; i < len(ps); i++ {
		p := express.Track(ps[i].ID)
		if p.Value != ps[i].Value {
			m := slack.NewMessage(fmt.Sprintf("Status of express %s", p.ID), "")
			m.AddField("Express Status", p.Value, false)
			slack.Post(expressHook, m)
			p.Update()
		}
		if p.Checked() {
			p.Del()
		}
	}
}
