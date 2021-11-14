universal configuration language improved
=========================================

it does objects with named and anonymous properties, that's it.
inspired by HCL which is inspired by UCL.
supports like 1% of its features, but does have the major advantage that it can be formated tabular
for readability when you have lots of similar entries.


for example

```hcl
rocket falcon heavy { yuge: true        destination: up         taste: "i don't know, please do not lick the rocket" }
rocket new shepard  { yuge: nah         destination: space      taste: "i don't know, but bezos looks tasty" }
rocket electron     { yuge: babyrocket  destination: fishies    taste: "zappy di zap zap"}

tweet {
    hashtag:    falconheavy
    author:     veryRealElonMusk81123311
    text:       "big rocket moon
like my new ponzi chain
click to win free integers
"}
```

equivalent json:

```json
{
    "rocket": {
        "falcon" : {
            "[1]": "heavy",
            "destination": "up",
            "taste": "i don't know, please do not lick the rocket",
            "yuge": "true"
        },
        ...
    ],
    "tweet": {
        "[0]": {
            "author": "veryRealElonMusk81123311",
            "hashtag": "falconheavy",
            "text": "big rocket moon\nlike my new ponzi chain\nclick to win free integers\n"
        }
    },
}

```


you probably want to combine this with https://github.com/mitchellh/mapstructure

```go
type Rocket struct {
    Name2       string  `mapstructure:"[1]"`
    Destination string
}

type Config struct {
    Rocket      map[string]Rocket
    Tweet       map[string]Tweet
}

func main() {
    f, err := os.Open("something.ucli");
    if err != nil { panic(err) }

    m, err := ucli.Parse(f)
    if err != nil { panic(err) }

    var config Fabric
    err = mapstructure.Decode(m, &config)
    if err != nil { panic(err) }
}

```
