/*
Copyright © 2022 GRISARD Dimitri dimitri.grisard03@gmail.com

*/

package cmd

import (
	"fmt"

	"golang.org/x/exp/slices"
)

var (
	host, user, password, typeOf, cfgFile string
	allowedType                           []string
)

func checkUcceType(allowedType []string, typeOf string) error {
	if !slices.Contains(allowedType, typeOf) {
		err := fmt.Errorf("this ucce instance type doesn't exist %s", typeOf)
		return err
	}

	return nil
}
