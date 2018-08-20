//go:generate env GO111MODULE=on packr2

package main

import (
	"context"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"github.com/sevein/amflow/cmd"
)

func main() {
	if err := cmd.Run(); err != nil {
		if errors.Cause(err) == context.Canceled {
			logrus.Debugln(errors.Wrap(err, "ignore error since context is cancelled"))
		} else {
			logrus.Fatal(err)
		}
	}
}