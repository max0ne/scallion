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
package camera

import (
	"io/ioutil"
	"time"

	"github.com/dhowden/raspicam"
	"github.com/pkg/errors"
)

func Capture() (string, error) {
	tempFile, err := ioutil.TempFile("", time.Now().Format("scallion-2006-01-02-15-04-05-*.jpg"))
	if err != nil {
		return "", errors.Wrapf(err, "unable to create temp file")
	}
	defer tempFile.Close()

	still := raspicam.NewStill()
	still.Quality = 100

	errCh := make(chan error)
	raspicam.Capture(still, tempFile, errCh)
	for err := range errCh {
		return "", errors.Wrapf(err, "unable to capture image")
	}
	return tempFile.Name(), nil
}
