# AldimVerdim

This project aims to fulfill the requirements of an assessment given by Teknasyon in [this](https://github.com/Teknasyon/assessments-backend) repository. The project is also developed for educational purposes. 

## Requirements for the project
- A .env file is required for the configuration of the application and dependencies. Each service has its own .env file and each template .env file can be found in respective services.
- Docker is required to run the project, all of the project is containerized.

## To run the project
```
docker-compose up --build -d
```

# Remarks
- Because this is a simple project, multiple frontends are neglected. It will still work when there are multiple frontends but might have errors due to misalignment of queued events.
- Addition to the remark above: it is not a good idea to sort the requested events, however, there is no way of knowing whether they are sorted or not. Financial transaction systems probably have better handling for such case.
- It is assumed to have only one worker for any scale, that means, the system works as much as the worker can handle. Therefore, vertical scaling should be chosen for larger scale of events.
- The locker in Worker will consume a lot of memory if it stays up for a long time. An optimization for the events is required but it depends on the Wallet ID of event itself.
- Better logging is needed.
- Unit tests are needed, will work on that soon.
