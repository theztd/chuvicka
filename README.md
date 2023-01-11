# chuvicka

Simple UP time robot implementation using GO and influxdb

## Key concepts

Every user has is own bucket, so the new registration make works as a login to the influzdb bucket
Every measured applicatin is the measurement in the influxdb.

The application has two parts:

### Dashboard
Vizualize status of the application, could be avalialble publicly on the request.

### Server
This is the management interface where user can add and remove his measured application.

### Agent
This is the colletcting app, that list all measurements in the bucket and for each append latest status. This part have to be most performance and run in parallel per bucket. After each run it reports cumulative time spend by measurement and report it to table **__customer_billing**

 
