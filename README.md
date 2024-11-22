# Kubernetes Microservices Demo
This repository provides a demonstration environment for a Kubernetes cluster, featuring multiple microservices that communicate over various protocols, including HTTP/REST and AMQP/RabbitMQ. Data generated by the microservices is stored in a MongoDB database, and monitoring is managed through Grafana with Loki as the data source.

## Key Features
- Multi-Protocol Communication: Microservices communicate using both HTTP/REST and RabbitMQ (ToDo: RabbitMQ package created and finished but not used yet.), showcasing different integration patterns.
- Centralized Monitoring: Grafana is deployed as the primary monitoring tool, with Loki as the log data source for real-time analysis.
- Automated Deployment: The cmds.sh script automates the deployment of all services using Helm charts, simplifying the setup process.
- Automated Testing: There is a Jenkins service connected to the repository that launches the Unit testing pipeline after every commit. Pull requests cannot be merged until the pipeline passes. 
## Important Note
This project is a proof of concept developed in my free time over the last weeks to demonstrate my expertise in Kubernetes, Go, and related technologies. It is not intended for production use. While this repository showcases some of my skills, I have over 6 years of experience as a software developer, working with various platforms and technologies that are not represented here.
