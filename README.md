# webapp
### **📌 Shortened Introduction for README**
The **Health Check API** monitors application instances, preventing traffic to unhealthy ones and enabling automatic recovery. With each successful check, a record is added to the database.

To run this application on a VM do the following:

## Prerequisites

- Create a VM on your preferred platform with the following specifications:- 
  - Server Operating System: Ubuntu 24.04 LTS 
  - At least 1 GB / 1 CPU 25 GB SSD Disk 1000 GB transfer
- Create a .env file with database credentials

## Instructions to Run the Application

- Download the zip of this repository and download the `app_setup.sh` script from this repository
- Copy the script, the zip and .env from your local path to the VM using scp
 
`scp -i <key> ~/Path/To/.zip app_setup.sh  ubuntu@ec2-XX-XXX-XXX-XXX.compute-1.amazonaws.com:/remote/path/`
- Run `chmod +x app_setup.sh` to make the script executable
- Run `. app_setup.sh` to execute the script. This will start the application
- Use Postman or cURL to test the application by hitting `http://<vm_public_ip>:8080/healthz`
