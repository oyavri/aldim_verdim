# AldimVerdim

This project aims to fulfill the requirements of an assessment given by Teknasyon in [this](https://github.com/Teknasyon/assessments-backend) repository. It is also developed for educational purposes. 

## Requirements for the project
- A .env file is required for the configuration of the application and dependencies. A template environment file can be found in the following path shared/config/.env.template
- Docker is required to run the project, all of the project is containerized. (At least aimed to be containerized)
- make is required to ease the process, however, one can use the commands listed in the makefile.

## To run the project
With make command:
```bash
make run
```

Without make command:
```bash
docker-compose up --build -d
```
