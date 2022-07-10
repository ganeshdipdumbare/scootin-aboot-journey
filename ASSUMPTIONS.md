## Assumptions

### Scooter booking assumption as follows - 
1. User/Mobile client will get the available scooters from his location within given radius. Note that getting scooters within sqare area is confusing.
2. User scans the QR code which send the BE request to book the scooter/begin trip.
3. BE system sends notification to scooter along with user id to unlock it.
4. The scooter sends trip start event and starts sending location update along with user id.
5. User stops the trip by calling BE api with current location. Scooter becomes free and the current scooter location is updated.
6. Scooter sends trip stop event.
7. The sequence of events getting stored in DB doesnt matter as the time of event creation is sent by client.
8. Usually only one user/client will try to scan and book the particular scooter and not more tha n one at a time.
9. User will always move to North by 10m per 3 Secons during trip with scooter.

All the above mentioned assumptions are also considered while implementing test clients. When service receives stop signal, server stops gracefully although clients are stopped abruptly. 

Also use of MongoDB is considered as it supports Geo spatial data and related operations.