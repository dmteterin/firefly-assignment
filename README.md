# Firefly Assignment

# Usage

Run the app using

`make run` or `go run cmd/*.go`

# Configuration

### Validator Parameters

| Parameter           | Default Value |   Enviromental Variable   |
| :------------------ | :-----------: | :-----------------------: |
| Minimal Word Length |      `3`      | VALIDATOR_MIN_WORD_LENGTH |
| Regular Expression  | `^[a-zA-Z]+$` |     VALIDATOR_REG_EXP     |

### File Path Parameters

| Parameter            |   Default Value   |  Enviromental Variable  |
| :------------------- | :---------------: | :---------------------: |
| Word Bank Path       |   `files/words`   |   WORD_BANK_FILE_PATH   |
| Crawler URLs Path    | `files/endg-urls` | CRAWLER_LINK_LIST_PATH  |
| Crawler Proxies Path |  `files/proxies`  | CRAWLER_PROXY_LIST_PATH |

### Crawler Parametersa

| Parameter                   | Default Value |        Enviromental Variable        |
| :-------------------------- | :-----------: | :---------------------------------: |
| Proxy Prefix                |  `socks5://`  |      CRAWLER_PROXY_TYPE_PREFIX      |
| Queue Consumer Thread Count |     `20`      | CRAWLER_QUEUE_CONSUMER_THREAD_COUNT |
| Retry Limit                 |      `5`      |       CRAWLER_URL_RETRY_LIMIT       |
| Enable Exponential Backoff  |    `false`    |     CRAWLER_ENABLE_EXP_BACKOFF      |

# Avoiding Rate Limits

### Proxies

More proxies with more threads should be used for faster scraping without hitting the rate limit

### Exponential backoff

Scraping with a low amount of proxies is slow, but could be done using exponential backoff

# Example Output

```json
[
  {
    "word": "the",
    "count": 726013
  },
  {
    "word": "and",
    "count": 366737
  },
  {
    "word": "that",
    "count": 203224
  },
  {
    "word": "you",
    "count": 140421
  },
  {
    "word": "with",
    "count": 125469
  },
  {
    "word": "The",
    "count": 124262
  },
  {
    "word": "has",
    "count": 67946
  },
  {
    "word": "have",
    "count": 66884
  },
  {
    "word": "from",
    "count": 65903
  },
  {
    "word": "more",
    "count": 58271
  }
]
```
