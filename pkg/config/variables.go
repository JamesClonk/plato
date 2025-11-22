package config

import (
	"path"

	"github.com/spf13/viper"
)

var (
	dirSource           = "templates"
	dirTarget           = "rendered"
	dirGeneratedSecrets = "rendered/secrets"
	delimiterLeft       = "{{{"
	delimiterRight      = "}}}"
)

func DirRoot() string {
	return path.Dir("/")
}

func DirSource() string {
	if len(viper.GetString("plato.source")) > 0 {
		return viper.GetString("plato.source")
	}
	return dirSource
}

func DirTarget() string {
	if len(viper.GetString("plato.target")) > 0 {
		return viper.GetString("plato.target")
	}
	return dirTarget
}

func DirGeneratedSecrets() string {
	if len(viper.GetString("plato.secrets")) > 0 {
		return viper.GetString("plato.secrets")
	}
	return dirGeneratedSecrets
}

func DelimiterLeft() string {
	if len(viper.GetString("plato.delimiters.left")) > 0 {
		return viper.GetString("plato.plato.delimiters.left")
	}
	return delimiterLeft
}

func DelimiterRight() string {
	if len(viper.GetString("plato.delimiters.right")) > 0 {
		return viper.GetString("plato.plato.delimiters.right")
	}
	return delimiterRight
}
