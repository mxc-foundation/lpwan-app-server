package data

import "time"

type CodecStruct struct {
	JS struct {
		MaxExecutionTime time.Duration `mapstructure:"max_execution_time"`
	} `mapstructure:"js"`
}
