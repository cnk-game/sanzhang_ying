// Copyright 2009  The "config" Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package config

import (
	"errors"
	"strconv"
	"strings"
)

// Bool has the same behaviour as String but converts the response to bool.
// See "boolString" for string values converted to bool.
func (c *Config) Bool(section string, option string) (value bool, err error) {
	sv, err := c.String(section, option)
	if err != nil {
		return false, err
	}

	value, ok := boolString[strings.ToLower(sv)]
	if !ok {
		return false, errors.New("could not parse bool value: " + sv)
	}

	return value, nil
}

// Float has the same behaviour as String but converts the response to float.
func (c *Config) Float(section string, option string) (value float64, err error) {
	sv, err := c.String(section, option)
	if err == nil {
		value, err = strconv.ParseFloat(sv, 64)
	}

	return value, err
}

// Int has the same behaviour as String but converts the response to int.
func (c *Config) Int(section string, option string) (value int, err error) {
	sv, err := c.String(section, option)
	if err == nil {
		value, err = strconv.Atoi(sv)
	}

	return value, err
}

// RawString gets the (raw) string value for the given option in the section.
// The raw string value is not subjected to unfolding, which was illustrated in
// the beginning of this documentation.
//
// It returns an error if either the section or the option do not exist.
func (c *Config) RawString(section string, option string) (value string, err error) {
	if _, ok := c.data[section]; ok {
		if tValue, ok := c.data[section][option]; ok {
			return tValue.v, nil
		}
	}
	return c.RawStringDefault(option)
}

// RawStringDefault gets the (raw) string value for the given option from the
// DEFAULT section.
//
// It returns an error if the option does not exist in the DEFAULT section.
func (c *Config) RawStringDefault(option string) (value string, err error) {
	if tValue, ok := c.data[DEFAULT_SECTION][option]; ok {
		return tValue.v, nil
	}
	return "", OptionError(option)
}

// String gets the string value for the given option in the section.
// If the value needs to be unfolded (see e.g. %(host)s example in the beginning
// of this documentation), then String does this unfolding automatically, up to
// _DEPTH_VALUES number of iterations.
//
// It returns an error if either the section or the option do not exist, or the
// unfolding cycled.
func (c *Config) String(section string, option string) (value string, err error) {
	value, err = c.RawString(section, option)
	if err != nil {
		return "", err
	}

	var i int

	for i = 0; i < _DEPTH_VALUES; i++ { // keep a sane depth
		vr := varRegExp.FindString(value)
		if len(vr) == 0 {
			break
		}

		// Take off leading '%(' and trailing ')s'
		noption := strings.TrimLeft(vr, "%(")
		noption = strings.TrimRight(noption, ")s")

		// Search variable in default section
		nvalue, _ := c.data[DEFAULT_SECTION][noption]
		if _, ok := c.data[section][noption]; ok {
			nvalue = c.data[section][noption]
		}
		if nvalue.v == "" {
			return "", OptionError(option)
		}

		// substitute by new value and take off leading '%(' and trailing ')s'
		value = strings.Replace(value, vr, nvalue.v, -1)
	}

	if i == _DEPTH_VALUES {
		return "", errors.New("possible cycle while unfolding variables: " +
			"max depth of " + strconv.Itoa(_DEPTH_VALUES) + " reached")
	}

	return value, nil
}
