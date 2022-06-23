# LiveServer

---
---

## Why LiveServer :

LiveServer is made to make webdeveloper's life easier with some useful features : ) 

---
## What is LiveServer:

LiveServer is basically a http server with some useful features to make webdev easier that you can
run at your local machine to test your webpages.

---
## How it Works:

Available soon... (Information about the architecture and working of LiveServer will be documented soon).

---
## KeyFeatures: 
    
### LiveReload:
    
    The server will listen to your webpages and it's resources.
    Automatically refreshs the weppage in browser if there is any changes in 
    the webages (html files and it's resource files like css) in the local machine.
(More details are available down).

### Dynamic Resource Linking:
    
    The resources can be linked to a url_request dynamically.
    No need to restart the server to link resource.
    and the linked resource and url_request_path can be also modified dynamically.
(More details are available down).
    
---

## How to build :

    Step 1 -> Install GoLang/Go in your local machine.
    Step 2 -> Clone the repository in your local machine.    
    Step 4 -> run cmd "go get github.com/gorilla/websocket".
    Step 5 -> run cmd "go build" at the directory where "LiveServer.go" is placed.    
now an Executable Binary will be produced in your current working directory.

---

## How to use :

Before starting make sure that the file "injectable_code" and the "Executable Binary" produced are in the same directory.
    
    Step 1 -> Start the Executable Binary (Just double click the executable binary).
    
Now you can see a terminal/console appearing which is ready to take your commands.
Currently available commands are 
                            
### 1.setmain:

    setmain can be used to specify the "MainHtmlPath" which will be served as "HomePage".
    Syntax: setmain MainHtmlPath

#### 2.setfavicon:

    setfavicon can be used to specify the "FavIconPath" which is automatically requested by the browser to set the icon of your webpage.
    Syntax: setfavicon FavIconPath

#### 3.link:
    
    link can be used to specify a resource (File) to be served for a request with particular url_path
    for example: If you wanna serve (respond) "logo.png" for the request with url_path /get/logo.png 
                 then you can use this command "link /get/logo.png PathFor_logo.png".
    Synatx: link Url_Path Resource_Path    
                

#### 4.start:
    
    start starts the HttpServer at port 80(default http port) or at a different port if specified.
    for example: If you wanna start the server at default (port 80) just Enter "start" as command.
                 else if you wanna start it in coustom port Enter "start PortNumber"
    Synatx: start PortNumber
    
### 5.help:
    
    help provides info about all available commands.    
Note: Commands are case sensitive only use lowercase letters.
---

## First TestRun:
    
Lets run this server and see what it can do. For this use the "example_webpage" folder provided with the repository.
    
    Step 1 -> Move the Executable Binary to the "example_webpage" folder for easy testing.
    Step 2 -> Start the Executable Binary.
    Step 3 -> Enter the cmd "start" to start the http server at your localmachine.
    Step 4 -> Now head to your browser and browse to this url "http://localhost/"
Now you should see a blank webpage it is blank because mainhtml is not seted.
    
    Step 5 -> Now Enter the cmd "setmain index.html" and refresh the page in browser.
You can see a simple webpage.

    Step 6 -> Open the "index.html" placed in the "example_webpage" folder and just save it with any changes. 
The LiveServer will automatically sense the changes and refreshs the webpage in the browser (LiveReloading).

The webpage is plain and doesn't have any css styling. So lets add it.
    
    Step 7 -> Enter the cmd "link /css/main.css css/main.css" and refresh the page (LiveReload for link command will be available soon).
Now you will get some styling.
The LiveReload will also work if you modify the "main.css" file or any other linked file which is served as resource by the server.
    
    Step 8 -> Use the cmd "setfavicon fav.png" and refresh the page 
You can see a icon in the tab of your webpage.

Note:The LiveReload will work only after a file is served for a webrequest.
---
