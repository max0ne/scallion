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
	"os"

	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	verbose bool
)

var rootCmd = &cobra.Command{
	Use: "scallion",
}

func Execute() {
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "")

	consoleEncoder := zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig())
	core := zapcore.NewCore(consoleEncoder, zapcore.Lock(os.Stderr), zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		if verbose {
			return lvl >= zapcore.DebugLevel
		} else {
			return lvl >= zapcore.InfoLevel
		}
	}))
	logger := zap.New(core)
	defer logger.Sync()
	zap.ReplaceGlobals(logger)

	cobra.CheckErr(rootCmd.Execute())
}
