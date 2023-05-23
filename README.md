# tempest-compression-service
This service takes files from an API request, compresses them, and then returns the file compressed.  
  
This service is called by the `tempest-data-service` and can be configured to call the `tempest-decider-service`  

# How to run  
this application contains a `Dockerfile` - this allows you run build and run the service using Docker console commands   
## Build  
```bash
docker build -t .
 ```
   
 ## Run  
 ```bash
docker run -p 8080:8080 -v . -e ENV_VARIABLE=value .
 ```
   
 ## Stop the container  
 ```bash
 docker stop container-name
 ```
