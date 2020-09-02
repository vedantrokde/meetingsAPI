# Meetings API

The task is to develop a basic version of meeting scheduling API. You are only required to develop the API for the system. Below are the details.

Meetings have the following Attributes. All fields are mandatory unless marked optional:

>>Id;
>>Title;
>>Partipants;
>>Start Time;
>>End Time;
>>Creation Timestamp;

Participants have the follwoing Attributes.

>>Name;
>>Email;
>>RSVP

The developed HTTP JSON API is capable of the following operations.

## Schedule a meeting 
```txt
Should be a POST request
Use JSON request body
URL should be '/meetings'
Must return the meeting in JSON format
```

## Get a meeting using id
```txt
Should be a GET request
Id should be in the url parameter
URL should be ‘/meeting/<id here>’
Must return the meeting in JSON format
```
# List all meetings within a time frame
```txt
Should be a GET request
URL should be ‘/meetings?start=<start time here>&end=<end time here>’
Must return a an array of meetings in JSON format that are within the time range
```

## List all meetings of a participant
```txt
Should be a GET request
URL should be ‘/articles?participant=<email id>’
Must return a an array of meetings in JSON format that have the participant received in the email within the time range
```

## Constraints
```txt
THE API IS DEVELOPED USING GO
MONGODB IS USED FOR STORAGE
```


## Errors 

##### Meetings may get overlapped i.e.one participant (uniquely identified by email) may have 2 or more meetings with RSVP Yes with any overlap between their times.

##### The server thread is not safe i.e. it may have race conditions especially when two meetings are being booked simultaneously for the same participant with overlapping time.

##### NO pagination is provided to the listpoint

##### NO unit tests are included

## Basic Model

##### Prone to errors. And the above mentioned errors are being nullified. The code will be updated soon. 

