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
package twitter

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"mime"
	"mime/multipart"
	"net/http"
	"net/url"
	"path"
	"strings"

	"github.com/dghubble/oauth1"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

const (
	MediaUploadEndpoint  = "https://upload.twitter.com/1.1/media/upload.json"
	StatusUpdateEndpoint = "https://api.twitter.com/1.1/statuses/update.json"

	// stepSize is amount of image bytes to upload in each HTTP request
	stepSize = 500 * 1024
)

type Twitter struct {
	httpClient *http.Client
}

func New(consumerKey, consumerSecret, accessToken, accessSecret string) *Twitter {
	config := oauth1.NewConfig(consumerKey, consumerSecret)
	token := oauth1.NewToken(accessToken, accessSecret)
	httpClient := config.Client(oauth1.NoContext, token)
	return &Twitter{
		httpClient: httpClient,
	}
}

func getMediaType(imagePath string) (string, error) {
	comps := strings.Split(imagePath, ".")
	extention := mime.TypeByExtension("." + comps[len(comps)-1])
	if extention == "" {
		return "", fmt.Errorf("unable to resolve mime type for file %s", imagePath)
	}
	return extention, nil
}

func (twitter *Twitter) mediaInit(length int, mediaType string) (string, error) {
	form := url.Values{}
	form.Add("command", "INIT")
	form.Add("media_type", mediaType)
	form.Add("total_bytes", fmt.Sprint(length))

	req, err := http.NewRequest("POST", MediaUploadEndpoint, strings.NewReader(form.Encode()))
	if err != nil {
		return "", err
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	res, err := twitter.httpClient.Do(req)
	if err != nil {
		return "", errors.Wrapf(err, "unable to make request")
	}

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", errors.Wrapf(err, "unable to load response body")
	}

	zap.L().With(zap.String("response_body", string(body))).Debug("response body")

	if res.StatusCode >= 400 {
		return "", fmt.Errorf("response code %d", res.StatusCode)
	}

	var mediaInitResponse struct {
		MediaId          uint64 `json:"media_id"`
		MediaIdString    string `json:"media_id_string"`
		ExpiresAfterSecs uint64 `json:"expires_after_secs"`
	}
	err = json.Unmarshal(body, &mediaInitResponse)
	if err != nil {
		return "", errors.Wrapf(err, "unable to parse json string %s", string(body))
	}

	return mediaInitResponse.MediaIdString, nil
}

func (twitter *Twitter) mediaAppend(fileName string, mediaID string, media []byte) error {
	for s := 0; s*stepSize < len(media); s++ {
		var reqBody bytes.Buffer
		rangeBegining := s * stepSize
		rangeEnd := (s + 1) * stepSize
		if rangeEnd > len(media) {
			rangeEnd = len(media)
		}

		zap.S().Debugf("try to append %d - %d", rangeBegining, rangeEnd)

		writer := multipart.NewWriter(&reqBody)
		writer.WriteField("command", "APPEND")
		writer.WriteField("media_id", mediaID)
		writer.WriteField("segment_index", fmt.Sprint(s))

		part, err := writer.CreateFormFile("media", fileName)
		if err != nil {
			return errors.Wrapf(err, "unable to load media writer")
		}
		_, err = part.Write(media[rangeBegining:rangeEnd])
		if err != nil {
			return errors.Wrapf(err, "unable to load media segment %d:%d", rangeBegining, rangeEnd)
		}
		// Important: this need to be closed here, not defer
		// Closing this writer triggers the form writer to write the terminating token into request body buffer
		if err = writer.Close(); err != nil {
			return errors.Wrapf(err, "unable to close")
		}

		req, err := http.NewRequest("POST", MediaUploadEndpoint, &reqBody)
		if err != nil {
			return errors.Wrapf(err, "unable to create request")
		}

		req.Header.Add("Content-Type", writer.FormDataContentType())
		req.Header.Add("Content-Length", fmt.Sprint(len(reqBody.Bytes())))
		res, err := twitter.httpClient.Do(req)
		if err != nil {
			return errors.Wrapf(err, "unable to make request")
		}

		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return errors.Wrapf(err, "unable to load response")
		}
		zap.L().With(zap.String("response_body", string(body))).Debug("response body")

		if res.StatusCode >= 400 {
			return fmt.Errorf("response code %d", res.StatusCode)
		}
	}

	return nil
}

func (twitter *Twitter) mediaFinilize(mediaID string) error {
	form := url.Values{}
	form.Add("command", "FINALIZE")
	form.Add("media_id", mediaID)

	req, err := http.NewRequest("POST", MediaUploadEndpoint, strings.NewReader(form.Encode()))
	if err != nil {
		return errors.Wrapf(err, "unable to construct request")
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	res, err := twitter.httpClient.Do(req)
	if err != nil {
		return errors.Wrapf(err, "unable to finalize")
	}
	if res.StatusCode >= 400 {
		return fmt.Errorf("response code %d", res.StatusCode)
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return errors.Wrapf(err, "unable to read response body")
	}

	zap.S().Debug("finalize response", string(body))
	return nil
}

func (twitter *Twitter) updateStatusWithMedia(text, mediaID string) (string, error) {
	form := url.Values{}
	form.Add("status", text)
	form.Add("media_ids", fmt.Sprint(mediaID))

	req, err := http.NewRequest("POST", StatusUpdateEndpoint, strings.NewReader(form.Encode()))
	if err != nil {
		return "", errors.Wrapf(err, "unable to construct request")
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	res, err := twitter.httpClient.Do(req)
	if err != nil {
		return "", errors.Wrapf(err, "unable to post status")
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", errors.Wrapf(err, "unable to read response body")
	}
	zap.L().With(zap.String("response_body", string(body))).Debug("response body")

	if res.StatusCode >= 400 {
		return "", fmt.Errorf("response code %d", res.StatusCode)
	}

	var tweetResponse struct {
		ID int64 `json:"id"`
	}
	if err = json.Unmarshal(body, &tweetResponse); err != nil {
		return "", err
	}
	if tweetResponse.ID == 0 {
		zap.L().Error("Unable to find tweet id from api response", zap.String("body", string(body)))
		return "", fmt.Errorf("unable to find tweet id from api response")
	}
	return fmt.Sprintf("https://twitter.com/scallionfriends/status/%d", tweetResponse.ID), nil
}

func (twitter *Twitter) Post(text, imagePath string) (string, error) {
	bytes, err := ioutil.ReadFile(imagePath)
	if err != nil {
		return "", errors.Wrapf(err, "unable to open image file %s", imagePath)
	}

	mediaType, err := getMediaType(imagePath)
	if err != nil {
		return "", errors.Wrapf(err, "unable to get mime type")
	}

	mediaID, err := twitter.mediaInit(len(bytes), mediaType)
	if err != nil {
		return "", errors.Wrapf(err, "unable to init media upload")
	}

	imageFileName := path.Base(imagePath)
	if err := twitter.mediaAppend(imageFileName, mediaID, bytes); err != nil {
		return "", errors.Wrapf(err, "unable to append upload upload")
	}

	if err := twitter.mediaFinilize(mediaID); err != nil {
		return "", errors.Wrapf(err, "unable to finalize media upload")
	}

	url, err := twitter.updateStatusWithMedia(text, mediaID)
	if err != nil {
		return "", errors.Wrapf(err, "unable to post tweet")
	}

	return url, nil
}
