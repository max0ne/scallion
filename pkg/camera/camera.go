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
