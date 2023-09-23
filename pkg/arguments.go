package pkg

import (
	"github.com/alexflint/go-arg"
)

type Arguments struct {
	Run string `default:"app"`
}

func ReadArguments() Arguments {
	args := Arguments{}
	arg.MustParse(&args)
	return args
}
