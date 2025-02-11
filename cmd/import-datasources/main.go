// This is a simple example of usage of Grafana client
// for importing datasources from a bunch of JSON files (current dir used).
// You are can export datasources with backup-datasources utitity.
// NOTE: old datasources with same names will be silently overrided!
//
// Usage:
//   import-datasousces http://sdk.host:3000 api-key-string-here
//
// You need get API key with Admin rights from your Grafana!
package main

/*
   Copyright 2016 Alexander I.Grafov <grafov@gmail.com>
   Copyright 2016-2019 The Grafana SDK authors

   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

	   http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.

   ॐ तारे तुत्तारे तुरे स्व
*/

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/K-Phoen/sdk"
)

func main() {
	var (
		datasources []sdk.Datasource
		filesInDir  []os.FileInfo
		rawDS       []byte
		status      sdk.StatusMessage
		err         error
	)
	if len(os.Args) != 3 {
		fmt.Fprint(os.Stderr, "Usage:  import-datasources http://sdk-host:3000 api-key-string-here\n")
		os.Exit(0)
	}
	ctx := context.Background()
	c, err := sdk.NewClient(os.Args[1], os.Args[2], sdk.DefaultHTTPClient)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create a client: %s\n", err)
		os.Exit(1)
	}
	if datasources, err = c.GetAllDatasources(ctx); err != nil {
		fmt.Fprint(os.Stderr, err)
		os.Exit(1)
	}
	filesInDir, err = ioutil.ReadDir(".")
	if err != nil {
		fmt.Fprint(os.Stderr, err)
	}
	for _, file := range filesInDir {
		if strings.HasSuffix(file.Name(), ".json") {
			if rawDS, err = ioutil.ReadFile(file.Name()); err != nil {
				fmt.Fprint(os.Stderr, err)
				continue
			}
			var newDS sdk.Datasource
			if err = json.Unmarshal(rawDS, &newDS); err != nil {
				fmt.Fprint(os.Stderr, err)
				continue
			}
			for _, existingDS := range datasources {
				if existingDS.Name == newDS.Name {
					if status, err = c.DeleteDatasource(ctx, existingDS.ID); err != nil {
						fmt.Fprintf(os.Stderr, "error on deleting datasource %s with %s", newDS.Name, err)
					}
					break
				}
			}
			if status, err = c.CreateDatasource(ctx, newDS); err != nil {
				fmt.Fprintf(os.Stderr, "error on importing datasource %s with %s (%s)", newDS.Name, err, *status.Message)
			}
		}
	}
}
