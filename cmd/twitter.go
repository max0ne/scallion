/*
Copyright © 2021 Mingfei Huang <himax1023@gmail.com>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"go.uber.org/zap"

	"github.com/max0ne/scallion/pkg/twitter"
)

var (
	credentialsFile string
	image           string
)

// twitterCmd represents the twitter command
var twitterCmd = &cobra.Command{
	Use:   "tweet",
	Short: "",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		if err := run(args); err != nil {
			zap.S().Panicf("%+v", err)
		}
	},
}

func init() {
	twitterCmd.PersistentFlags().StringVarP(&credentialsFile, "credentials", "c", "",
		"A JSON file containing Twitter credentials, must include keys: consumerKey, consumerSecret, accessToken, accessSecret",
	)
	twitterCmd.PersistentFlags().StringVarP(&image, "image", "i", "",
		"Path to tweet image file",
	)
	rootCmd.AddCommand(twitterCmd)
}

func run(args []string) error {
	file, err := os.Open(credentialsFile)
	if err != nil {
		return errors.Wrapf(err, "unable to open credentials file %s", credentialsFile)
	}

	var creds struct {
		ConsumerKey    string `json:"consumerKey"`
		ConsumerSecret string `json:"consumerSecret"`
		AccessToken    string `json:"accessToken"`
		AccessSecret   string `json:"accessSecret"`
	}
	if err := json.NewDecoder(file).Decode(&creds); err != nil {
		return errors.Wrapf(err, "unable to parse credentials file %s", credentialsFile)
	}

	if len(creds.ConsumerKey) == 0 {
		return fmt.Errorf("consumerKey missing in credentials file %s", credentialsFile)
	}
	if len(creds.ConsumerSecret) == 0 {
		return fmt.Errorf("consumerSecret missing in credentials file %s", credentialsFile)
	}
	if len(creds.AccessToken) == 0 {
		return fmt.Errorf("accessToken missing in credentials file %s", credentialsFile)
	}
	if len(creds.AccessSecret) == 0 {
		return fmt.Errorf("accessSecret missing in credentials file %s", credentialsFile)
	}

	twee := twitter.New(
		creds.ConsumerKey,
		creds.ConsumerSecret,
		creds.AccessToken,
		creds.AccessSecret,
	)

	message := strings.Join(args, " ")
	zap.S().Infof(`Tweeting "%s" with image file %s`, message, image)

	url, err := twee.Post(message, image)
	if err != nil {
		return err
	}

	// Write tweet url to stdout
	fmt.Println(url)
	return nil
}
