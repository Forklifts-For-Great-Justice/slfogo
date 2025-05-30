package slfogolib

import (
	"context"
	"fmt"

	"gopkg.in/mcuadros/go-syslog.v2/format"
)

type OutputProcessor interface {
	Put(context.Context, format.LogParts) error
	Close() error
}

func getKey(lp format.LogParts, key string) (string, error) {
	v, ok := lp[key]
	if !ok {
		return "", fmt.Errorf("Not value for key \"%s\"", key)
	}

	retVal, ok := v.(string)
	if !ok {
		return "", fmt.Errorf("lp[\"%s\"] cannot be convert to type string", key)
	}

	return retVal, nil
}
