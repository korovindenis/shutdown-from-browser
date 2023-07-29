<p align="center">
  <img height="450px" alt="Shows an illustrated sun in light mode and a moon with stars in dark mode." src="https://github.com/korovindenis/shutdown-from-browser/raw/master/.github/img/homer.png">
</p>

![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/korovindenis/shutdown-from-browser)
![GitHub code size in bytes](https://img.shields.io/github/languages/code-size/korovindenis/shutdown-from-browser)
![GitHub](https://img.shields.io/github/license/korovindenis/shutdown-from-browser)


## Description

**Shutdown from Browser** is a simple web server built in Go that allows users to shut down their computer remotely using a web browser. With this tool, you can conveniently power off or restart your computer without needing to be physically present.

The web server uses React to provide a modern and responsive user interface, making it easy and intuitive for users to interact with the shutdown and restart functionality.

## Installation and Running

To build and run the **Shutdown from Browser** web server, you can use the included Makefile:

1.  Clone the repository:
    
    `git clone https://github.com/korovindenis/shutdown-from-browser.git
    cd shutdown-from-browser` 
    
2.  Initialize the `CONFIG_PATH` environment variable to point to the location of your configuration file:
    
    `export CONFIG_PATH=/path/to/your/configs/config.prod.yaml` 
    
3.  Use the Makefile to build the server:
       
    `make build` 
    
4.  Run the server:
    
    `make run` 
    
    The web server will be available at `http://localhost:8081/`.

### Install as a Service

Optionally, you can install **Shutdown from Browser** as a service to ensure it runs automatically on system startup. The Makefile provides a `make install` target to do this for you.

`make install` 

## Configuration

The **Shutdown from Browser** web server allows you to customize its behavior by modifying the `config.yaml` file. This file contains various settings that you can adjust according to your preferences. 

The `config.prod.yaml` file is straightforward to use, and it provides comments to guide you through the available settings. After making changes to the configuration, remember to restart the server for the modifications to take effect.

## License

This project is licensed under the [Apache License 2.0](https://github.com/korovindenis/shutdown-from-browser/blob/master/LICENSE.md).

## Contributions and Feedback

If you have any suggestions, find any issues, or would like to contribute to the project, feel free to create an Issue or Pull Request. We value your contributions and strive to improve the project continuously.

Thank you for using **Shutdown from Browser**! We hope this tool makes it more convenient for you to remotely control your computer's power state.