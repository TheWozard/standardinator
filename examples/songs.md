# Songs
This examples provides a basic overview of the intended functionality of this application

## Input data
```json
[
    {
        "album": "Abbey Road",
        "cover_art": "<url>",
        "songs": [
            {
                "name":"Come Together",
                "link": "<url>",
                "details": {
                    "duration":"4:19"
                }
            },
            {
                "name":"Something",
                "link": "<url>",
                "details": {
                    "duration":"3:27"
                }
            }
        ]
    }
]
```

## Configuration
```json
{
    "parser": "JSON",
    "outputs": [
        {
            "name":"Song",
            "for": "$[*].songs[*]",
            "target": {
                "song": {
                    "name": "@.name",
                },
                "album": "@Album",
                "duration": "@.details.duration"
            }
        }
        {
            "name":"Album",
            "for": "$[*]",
            "omit": true,
            "target": {
                "name": "@.album",
            }
        }
    ]
}
```