# go-translator

# About

It's a simple "translator" service wrote by using go-lang. 
Possible only EN-RU direction of translation.

## Using

The interface to "translator" is a REST HTTP API. 

To translate the list of words, necessary make a POST request with json in following format:

```
{"words":["first", "second", "third", "etc"]}
```

To url:

```
scheme://address:8080
```
And get the following response:

{"first":"translate", "second":"translate", "third":"translate", "etc":"translate"}

Example:

```
curl  --header "Content-Type: application/json" \
      -X POST -d '{"words": ["what", "sitting", "world", "whose", "mine",  "here", "know", "ololo"]}' \
      http://localhost:8080
      
{"here":"здесь","know":"знать","mine":"шахта","ololo":"none","sitting":"заседание","whose":"чья","world":"мир"}
```

## Translator engine

Words are translated using "Yandex dictionary API" https://tech.yandex.ru/dictionary/.
So necessary pass the API key to the "translator" service using environment variable "TOKEN_YA".

## Docker way

For convenience of running application created Dockerfile, so possible to build and run docker image. 
