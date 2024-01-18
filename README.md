# kursRates
This web service created to receive and store exchange rates from the National Bank of Kazakhstan. The API built by using Go and uses the Gorilla Mux router to process HTTP requests.

# Description
Here I used Gorilla router, because it makes it easier to extract the parameters from the URL and use them in further operations. And also used Go MySQL Driver because it has a large audience of users, allows you to work effectively with databases.

# Endpoints
POST /currency/save/{date} - get dates from API National Bank with a given date and save it to Database "Test" <br />
GET /currency/{date}/{code} - return certain exchnge rate with a given date and with a given code in the JSON format <br />
GET /currency/{date} - return certain exchnge rate with a given date in the JSON format <br />
GET /currency/{date}/{code} - delete data from database by date and code <br />
GET /currency/{date} - delete data from database by date <br />
GET /Health - Healthchecker
