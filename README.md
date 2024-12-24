# Device Management ⚡️

## Let's Run the Application 

Rename .env.example to .env and then fill it out with the real values.

```
docker-compose up
or
make migrate-up
make run
```
## Unit testing
```
make test
```

### Rest Api
Payload: DeviceManagement.postman_collection.json

### Domain Model
![Image](./assets/event-modeling.png?raw=true)

### Vertical Slice Architecture
[View Documentation](assets/vertical_slice_and_event_modeling.pdf)