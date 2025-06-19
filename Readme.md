# Task

Develop a task management mechanism accessible via HTTP.
The system should provide the following capabilities:
- Users can create new tasks via an HTTP endpoint.
- Newly created tasks should be added to a queue for asynchronous processing.
- Users can list all tasks and check the status of individual tasks using their task ID.
- The system must include both a task queue and a task management component to handle task execution and status tracking.

## Run the service with:
```bash
make run
```

## Usage (Requests) - (via Postman for ex.)

```
POST - http://localhost:2025/task - making request for mocking data creation, no need post params
GET - http://localhost:2025/tasks  - get the all user's tasks list with statuses
GET - http://localhost:2025/task/3 - get specific task by taskID

```
P.S.:
Could also add additional method for getting tasks by userID. 
In addition, usage of maps probably better solution than slices for tasks storage, currently implemented slices 
