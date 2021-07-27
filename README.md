# Standardinator

![test](https://github.com/TheWozard/standardinator/actions/workflows/test.yml/badge.svg)

Golang package for converting streams of data into a standardized formats using [JSONPath](https://goessner.net/articles/JsonPath/)

## Target Functionality
Current the target functionality is being worked out and examples can be found in [examples](examples/index.md)

# Config Overview
```go
config := standardinator.Config{
    Parser: standardinator.JSONParser(io.Reader),
    Outputs: []standardinator.Outputs{
        {
            Name: "Name"
            For: "$[*]"
            Target: map[string]interface{}{
                
            }
        }
    }
}
```