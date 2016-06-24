# Web Sorter

This simple web server will take an array of words and sort them either alphabetically or reverse-alphabtically.

## Usage

```go run main.go```

This webserver only runs on 8080 and you can't do anything about it...=(

Only HTTP POST is allowed. Content-Type must be "application/json". For the request body, use a properly formatted JSON object.

### Request

```
{
    "words": ["cat", "mouse", "dog"],
    "reverse": true
}
```

* **words** - an array of words to sort where each word matches ```^[A-Za-z]$```
* **reverse** - if true, a reverse sort will be done. This value defaults to false

### Response

```
{
    "result": ["mouse", "dog", "cat"],
    "reverse": true
}
```

* **result** - the result of sorting the array of words
* **reverse** - whether or not a reverse sort was done

### Example

```
curl -X POST -H "Content-Type: application/json" -H "Cache-Control: no-cache" -H "Postman-Token: 2625c772-606b-2822-a402-079c0f98b1f8" -d '{
 "words":["mouse","dog","cat"],
 "reverse": true
}' "http://127.0.0.1:8080/sort"

HTTP/1.1 200 OK
Content-Type: application/json
Date: Fri, 24 Jun 2016 04:02:10 GMT
Content-Length: 48

{"result":["mouse","dog","cat"],"reverse":true}
```

## Errors
### 400 Bad Request
If the JSON is malformed, or any of the words in the "words" array contains characters other than [A-Za-z]

### 404 Not Found
Any route other than /sort is undefined

### 405 Method Not Allowed
Only HTTP POST is allowed with this endpoint

### 415 Unsupported Media Type
If the Content-Type of the request is not set to "application/json"
