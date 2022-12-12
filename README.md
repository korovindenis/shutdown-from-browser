![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/korovindenis/shutdown-from-browser)
![GitHub code size in bytes](https://img.shields.io/github/languages/code-size/korovindenis/shutdown-from-browser)

The service for linux allowing you to manage the power of the system (server or PC) from the web browser or through the rest requests to api

#### The application allows:
- to restart the system
- to turn off the system
- to turn off the system according to the scheduler (for example, after n-hours)
- through the rest api to know or set the time for the automatic shutdown or restart the system immediately 

The configuration file configures the port on which the application will run and the logging level
With the help of the Makefile you can install the application as a linux service

#### Installation
```sh
Make build-app
Make install
```
