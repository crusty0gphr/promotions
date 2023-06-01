# Promotions (case study task)
- - -
## Task descriptions 
- We receive some records in a CSV file (example promotions.csv attached) every 30 minutes. We would like to store these objects in a way to be accessed by an endpoint.
- Expected Result
```
{
  "id": "172FFC14-D229-4C93-B06B-F48B8C095512",
  "price": 9.68,
  "expiration_date": "2022-06-04 06:01:20"
}
```
- - -
## How to use
- Run postgres
```
docker-compose up
```
- Create table inside the DB using sql file
```
sql/db.sql
```
- Run application
```
go run cmd/main.go
```
- Get promotion request example
```
curl --location '0.0.0.0:7700/promotions/334'
```
- Populate DB with data request example
```
curl --location --request POST '0.0.0.0:7700/update'
```
- - -
## Case study results
I don't know how the file would be inserted into the app, so I assumed the file would be stored in a privae bucket like S3 or Minio. I have designed the application in such a way that it will receive a signal (```POST /update```) which will trigger the update process and the application will request the file from the bucket and update the DB. The file is stored directly in the ```/bucket``` directory, which represents the resource bucket.

Another way is to send the file directly to the application using HTTP and read it as ```multipart form data```. Just stream the file over HTTP with a specific content length, define a buffer size of 5MB for example, and write the parts of the file piece by piece.
- - -
- I created an application and uploaded a CSV file to it. On the first try, I just opened the file and stored the result inside the slice to prepare the SQL for the bulk insert into the database. I soon realized that this is the worst way to handle large files. My laptop has 32 GB of RAM and running locally won't be a problem. But this application will most likely run inside a container in a constrained environment, so applications are expected to use no more than 50 MB of RAM. The current file is 15MB and 50MB is more than enough, but if the next file is more than 1GB, app will simply run out of memory and crash instantly. It's more efficient to open a file and read it line by line until the end of the file. I have been using Postgres as storage for my application, so another question is how am I going to insert so much data into the DB. The 15 MB CSV file contains approximately 250k records. So if I read the file line by line and generate an INSERT query for each line, I get 250k INSERT queries. Doing so many inserts in a minute or less is an unnecessary load on the database. More inserts will just kill Postgres.
- - -
- On my second try, I researched how to insert 500k records into Postgres. One solution is to use an ORM tool. GORM has built-in functionality to migrate data, but it still requires a data structure stored inside the slice. This means that GORM will still iterate (perhaps more efficiently than a simple for loop) over all the objects and insert them into the table, but it still needs to load the entire file into memory.
- - - 
- On my third try, I did more research and came across Postgres' built-in COPY function, which allows you to copy data into tables directly from a file system file. My approach is to create one file that returns ```io.Reader``` and pass that reader to Postgres as STDIN. The query looks like this ```copy table_name from stdin with (format csv)```. This method saves tons of memory and provides great speed because this operation is performed directly inside the database using its own optimizations.
- - -

## Performance and scalability
Most functions in the application have linear constant complexity.
- service.go:25 ```f.getOne - O(1) => 1*N^0 + 1*N^0 + 1*N^0```
- service.go:46 ```f.updateData - O(1) => 1*N^0 + 1*N^0```
---
- storage.go:36 ```f.getById - O(N) => 1*N^0 + 1*N^1```
- storage.go:49 ```f.batchInsert - O(1) => 1*N^0```
- storage.go:60 ```f.copyFrom - O(N) => 1*N^0 + 1*N^0 + 1*N^0 + 1*N^0 + 1*N^1 + 1*N^0```
- storage.go:91 ```f.clearDB - O(1) => 1*N^0```
- - -
To work with this application in a production environment, one solution is to run this application natively inside the operating system, for example on an Ubuntu server, just clone it in the repository, run it with a command, and set up an access point for the HTTP server. In this scenario, we can scale the application vertically. If it's a physical machine, add more RAM and a faster CPU, or if the application is containerized, we can allocate more vRAM and add more CPU vCores according to the needs of the application.
- - -
Another way to run an application in a production environment is to use a container orchestration system. For example, we can use Kubernetes, which can deploy, scale, and manage pods. By configuring autoscaling, Kubernetes will automatically add additional modules to the namespace and automatically balance them using the built-in ingress service. The application must contain a Dockerfile that will be used to containerize it and then deployed inside the Kubernetes namespace as an independent pod.
- - -
You can use Kibana or Grafana to monitor a running application. Both of them provide dashboard for various metrics and logging. You can use Opentelemetry to monitor load and latency. Also, the application should be covered with a decent amount of logs to help the debugging process.
