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
	"strings"

	"github.com/skitta/goquery/api/openweather"
	"github.com/skitta/goquery/api/slack"
	"github.com/spf13/cobra"
)

const (
	weatherHook = "https://hooks.slack.com/services/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"
)

// weatherCmd represents the weather command
var weatherCmd = &cobra.Command{
	Use:   "weather",
	Short: "Weather report daily",
	Long:  `Get today's weather form openweather, and post to your slack webhook.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		city := strings.Join(args, "")
		info, result, err := openweather.GetDaily(city)
		if err != nil {
			return err
		}

		m := slack.NewMessage(fmt.Sprintf("Today's weather in %s", city), "")

		if info == "Rain" {
			m.AddColor("#D00000")
		}
		m.AddField(info, result, false)
		slack.Post(weatherHook, m)

		return nil
	},
}

func init() {
	RootCmd.AddCommand(weatherCmd)
}
