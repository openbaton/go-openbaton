/*
 *  Copyright (c) 2017 Open Baton (http://openbaton.org)
 *
 *  Licensed under the Apache License, Version 2.0 (the "License");
 *  you may not use this file except in compliance with the License.
 *  You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 *  Unless required by applicable law or agreed to in writing, software
 *  distributed under the License is distributed on an "AS IS" BASIS,
 *  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *  See the License for the specific language governing permissions and
 *  limitations under the License.
 */

/*
Package plugin implements a runtime for OpenBaton plugins.

Currently, only vim-drivers are supported; see the go-vimdriver-test repo for a sample implementation.

A new vim-driver plugin can be created by using the New() function together with a Driver instance:

    var driver plugin.Driver = &myDriver{}

    params := &plugin.Params{
        // insert your config here
    }

    plug, err := plugin.New(driver, params)
    if err != nil {
        panic("error: " + err.Error())
    }


Ensure that your VIMDriver implements the plugin.Driver interface.

The new plugin.Plugin can then be started using its Serve() method, blocking the current goroutine.
Use Stop() to stop the service and quit.


    if err := plug.Serve(); err != nil {
        panic("error while setting up plugin: " + err.Error())
    }
*/
package plugin
