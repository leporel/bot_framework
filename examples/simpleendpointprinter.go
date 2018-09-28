/*
MIT License

Copyright (c) 2017 MichiVIP

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
 */
package skypeapiexamples

import (
	"github.com/leporel/bot_framework/bfmodels"
	"encoding/json"
	"fmt"
)

/*
In this example I am setting up a basic skype APi endpoint and print activity objects.
We will just use the given http.Server to listen to incoming requests.
 */

func startSimpleEndpointPrinter() {
	// bad practice. In real production you should better request the token via skypeapi.RequestAccessToken
	// WARNING: when using a static authorization token it could expire. In future the will be an automatic refresher
	authorizationBearerToken := "YOUR-AUTH-TOKEN"

	// Endpoint is going to listen on 0.0.0.0:8080
	endpoint := NewEndpoint(":8080")

	// we define our own handle function
	srv := endpoint.SetupServer(*NewEndpointHandler(func(activity *bfmodels.Activity) {
		bytes, _ := json.MarshalIndent(activity, "", "  ")
		fmt.Println(string(bytes))
	}, authorizationBearerToken, "YOUR-APP-ID"))
	// finally we just use the default method to start the server
	srv.ListenAndServeTLS("certs/fullchain.pem", "certs/privkey.pem")
}
