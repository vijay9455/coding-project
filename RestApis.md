## APIs List
Overall there are 6 endpoints
### 1. Create User
Sample Request
```
  curl -X POST http://localhost:3000/api/v1/users --data '{"first_name":"y","last_name":"x","email":"y@z.com"}'
```
Response
```
{
  "api_version": "1.0.0",
  "data": {
    "user": {
      "availabilities": [
        { "day_of_week": 1, "end_time": "17:00:00", "start_time": "09:00:00" },
        { "day_of_week": 2, "end_time": "17:00:00", "start_time": "09:00:00" },
        { "day_of_week": 3, "end_time": "17:00:00", "start_time": "09:00:00" },
        { "day_of_week": 4, "end_time": "17:00:00", "start_time": "09:00:00" },
        { "day_of_week": 5, "end_time": "17:00:00", "start_time": "09:00:00" }
      ],
      "email": "y@z.com",
      "id": "usr_2jL5ihUk9O775PisxQ6SKgWThJX"
    }
  },
  "success": true
}
```
By default it is assumed and we create availability as part of user create for 9 AM to 5 PM on week days
### 2. My Profile
All the apis other than create user is authenticated, for now we send user_id as header to identify the current user. Ideally some other authentication should be implemented for now user_id is used to reduce dev effort
Sample Request
```
  curl http://localhost:3000/api/v1/my_profile -H 'calendly-user-id: usr_2jL5ihUk9O775PisxQ6SKgWThJX'
```
Response
```
{
  "api_version": "1.0.0",
  "data": {
    "user": {
      "availabilities": [
        { "day_of_week": 1, "end_time": "17:00:00", "start_time": "09:00:00" },
        { "day_of_week": 2, "end_time": "17:00:00", "start_time": "09:00:00" },
        { "day_of_week": 3, "end_time": "17:00:00", "start_time": "09:00:00" },
        { "day_of_week": 4, "end_time": "17:00:00", "start_time": "09:00:00" },
        { "day_of_week": 5, "end_time": "17:00:00", "start_time": "09:00:00" }
      ],
      "email": "y@z.com",
      "id": "usr_2jL5ihUk9O775PisxQ6SKgWThJX"
    }
  },
  "success": true
}
```
### 3. Update Availability
Update availability can be used to update the availability of the user on particular day of week. To mark day of week as unavailable we need to pass `mark_unavailable` as `true`. Else expected availability details need to be sent as array.
Sample Request
```
  curl -X PUT http://localhost:3000/api/v1/users/availabilities -H 'calendly-user-id: usr_2jL5ihUk9O775PisxQ6SKgWThJX' --data '{"day_of_week":1,"mark_unavailable":false,"availabilities":[{"start_time":"09:00:00","end_time":"12:00:00"},{"start_time":"16:00:00","end_time":"20:00:00"}]}'
```
Response
```
{
  "api_version": "1.0.0",
  "data": {
    "availabilities": [
      { "day_of_week": 1, "end_time": "12:00:00", "start_time": "09:00:00" },
      { "day_of_week": 1, "end_time": "20:00:00", "start_time": "16:00:00" }
    ]
  },
  "success": true
}
```
### 4. Available Slots
This endpoint provide the availability details date wise
Sample Request
```
curl "http://localhost:3000/api/v1/available_slots?start_date=2024-07-01&end_date=2024-08-01" -H 'calendly-user-id: usr_2jL5ihUk9O775PisxQ6SKgWThJX'
```
Response
```
{
  "api_version": "1.0.0",
  "data": {
    "availability_data": [
      { "date": "2024-07-01", "available": false, "slots": null },
      { "date": "2024-07-02", "available": false, "slots": null },
      { "date": "2024-07-03", "available": false, "slots": null },
      { "date": "2024-07-04", "available": false, "slots": null },
      { "date": "2024-07-05", "available": false, "slots": null },
      { "date": "2024-07-06", "available": false, "slots": null },
      { "date": "2024-07-07", "available": false, "slots": null },
      { "date": "2024-07-08", "available": false, "slots": null },
      { "date": "2024-07-09", "available": false, "slots": null },
      { "date": "2024-07-10", "available": false, "slots": null },
      { "date": "2024-07-11", "available": false, "slots": null },
      { "date": "2024-07-12", "available": false, "slots": null },
      { "date": "2024-07-13", "available": false, "slots": null },
      { "date": "2024-07-14", "available": false, "slots": null },
      { "date": "2024-07-15", "available": false, "slots": null },
      {
        "date": "2024-07-16",
        "available": true,
        "slots": [{ "start_time": "09:00:00", "end_time": "17:00:00" }]
      },
      {
        "date": "2024-07-17",
        "available": true,
        "slots": [{ "start_time": "09:00:00", "end_time": "17:00:00" }]
      },
      {
        "date": "2024-07-18",
        "available": true,
        "slots": [{ "start_time": "09:00:00", "end_time": "17:00:00" }]
      },
      {
        "date": "2024-07-19",
        "available": true,
        "slots": [{ "start_time": "09:00:00", "end_time": "17:00:00" }]
      },
      { "date": "2024-07-20", "available": false, "slots": [] },
      { "date": "2024-07-21", "available": false, "slots": [] },
      {
        "date": "2024-07-22",
        "available": true,
        "slots": [
          { "start_time": "09:00:00", "end_time": "12:00:00" },
          { "start_time": "16:00:00", "end_time": "20:00:00" }
        ]
      },
      {
        "date": "2024-07-23",
        "available": true,
        "slots": [{ "start_time": "09:00:00", "end_time": "17:00:00" }]
      },
      {
        "date": "2024-07-24",
        "available": true,
        "slots": [{ "start_time": "09:00:00", "end_time": "17:00:00" }]
      },
      {
        "date": "2024-07-25",
        "available": true,
        "slots": [{ "start_time": "09:00:00", "end_time": "17:00:00" }]
      },
      {
        "date": "2024-07-26",
        "available": true,
        "slots": [{ "start_time": "09:00:00", "end_time": "17:00:00" }]
      },
      { "date": "2024-07-27", "available": false, "slots": [] },
      { "date": "2024-07-28", "available": false, "slots": [] },
      {
        "date": "2024-07-29",
        "available": true,
        "slots": [
          { "start_time": "09:00:00", "end_time": "12:00:00" },
          { "start_time": "16:00:00", "end_time": "20:00:00" }
        ]
      },
      {
        "date": "2024-07-30",
        "available": true,
        "slots": [{ "start_time": "09:00:00", "end_time": "17:00:00" }]
      },
      {
        "date": "2024-07-31",
        "available": true,
        "slots": [{ "start_time": "09:00:00", "end_time": "17:00:00" }]
      }
    ]
  },
  "success": true
}
```
### 5. Overlapping Slots
This endpoint takes date as input and provide overlapping slots as output, if there are any
Sample Request
```
curl "http://localhost:3000/api/v1/overlapping_slots?date=2024-07-22&email=a@b.com" -H 'calendly-user-id: usr_2jL5ihUk9O775PisxQ6SKgWThJX'
```
Response
```
{
  "api_version": "1.0.0",
  "data": {
    "slots": [
      { "start_time": "09:00:00", "end_time": "12:00:00" },
      { "start_time": "16:00:00", "end_time": "17:00:00" }
    ]
  },
  "success": true
}
```
### 6. Create Meeting
This endpoint creates the meeting, if free slot is available
Sample Request
```
curl -X POST http://localhost:3000/api/v1/meetings  -H 'calendly-user-id: usr_2jL5ihUk9O775PisxQ6SKgWThJX' --data '{"email":"a@b.com","start_time":"2024-07-22T16:30:00Z","end_time":"2024-07-22T17:00:00Z","title":"first meeting","meeting_description":"dummy"}'
```
Response
```
{
  "api_version": "1.0.0",
  "data": {
    "meeting": {
      "end_time": "2024-07-22T17:00:00Z",
      "id": "meet_2jL6EBITvNd0wszGXo3YfQuY9x5",
      "meeting_description": "dummy",
      "participants": [
        { "accept_status": "ACCEPTED", "email": "y@z.com" },
        { "accept_status": "MAY_BE", "email": "a@b.com" }
      ],
      "start_time": "2024-07-22T16:30:00Z",
      "title": "first meeting"
    }
  },
  "success": true
}
```
Sample error response

HTTP/1.1 400 Bad Request
Date: Tue, 16 Jul 2024 18:17:49 GMT
Content-Length: 95
Content-Type: text/plain; charset=utf-8
```
{
  "api_version": "1.0.0",
  "error": { "code": "bad_request", "message": "no free slot" },
  "success": false
}
```
Sample Overlapping slots response after meeting
```
curl -i "http://localhost:3000/api/v1/overlapping_slots?date=2024-07-22&email=a@b.com" -H 'calendly-user-id: usr_2jL5ihUk9O775PisxQ6SKgWThJX'
```
Response
```
{
  "api_version": "1.0.0",
  "data": {
    "slots": [
      { "start_time": "09:00:00", "end_time": "12:00:00" },
      { "start_time": "16:00:00", "end_time": "16:30:00" }
    ]
  },
  "success": true
}
```