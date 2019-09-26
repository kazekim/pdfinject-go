/*
  GoLang code created by Jirawat Harnsiriwatanakit https://github.com/kazekim
*/

package pdfinject

import "os/user"

func HomeDirectory() (*string, error){
	usr, err := user.Current()
	if err != nil {
		return nil, err
	}

	return &usr.HomeDir, nil
}
