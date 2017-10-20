# Fission Fuction

This API create a new Ticket in Fresh Desk, it used Fresh Desk API

API ref -- https://developers.freshdesk.com/api/#create_ticket_with_attachment



Request Payload
` {"email":"amitkumarvarman@gmail.com","subject":"TV claim","description":" \n Make of the TV : Samsung \n Model number of the TV : QLED 55\" \n Please upload a detailed image of the damage :  \n When did the incident happen? :  \n Do you still have the TV in your possession? :  \n We have made the following assumptions about your property, you and anyone living with you :  \n Was there a theft or a deliberate event by a 3rd party? :  \n If you have any other insurance or warranties covering your TV, please advise us of the company name. : no \n Are you aware of anything else relevant to your claim that you would like to advise us of at this stage? : no \n In as much detail as possible, please provide a full written account of what has happened to your TV, including where it happened. : Tv broke... smashed to bits \n Serial number of the TV : Don't have it \n What was the purchase price of the TV? : ","status":2,"priority":1,"name":"TEST NAME"}`  

Response JSON

HTTP STATUS : 201 - Created
