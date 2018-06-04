Topic
-----

### Publish a topic:

* URL: POST http://localhost:8080/topic/{topic}/publish

* BODY:

Note: As the data is a byte slice, you can POST base64 encoded data into data param

```json
{
  "messages": [
    {
      "data": "aG9sYSBtdW5kbyEhIQ==",
      "attributes": {
        "attrKey1": "attrValue1",
        "attrKey2": "attrValue2"
      }
    },
    {
      "data": "YWRpb3MgbXVuZG8h"
    }
  ]
}
```

Subscription
------------

### Create a subscription:

* URL: PUT http://localhost:8080/subscription/{subscription_name}

* BODY:

```json
{
	"topic": "my-topic",
	"ack_deadline_seconds": 10
}
```

### Pull a subscription:

* URL: POST http://localhost:8080/subscription/{topic_name}/pull

* BODY:

```json
{
	"max_messages": 100
}
```

## Messages ack:

* URL: POST http://localhost:8080/subscription/{subscription_name}/ack

* BODY:

```json
{
  "ack_ids": [
    "ceed47ec-6415-11e8-ae84-a860b6301a14",
    "ceed4aa7-6415-11e8-ae84-a860b6301a14"
  ]
}
```

## Messages modify ack:

* URL: POST http://localhost:8080/subscription/{subscription_name}/ack/modify

* BODY:

```json
{
  "ack_ids": [
    "ceed47ec-6415-11e8-ae84-a860b6301a14",
    "ceed4aa7-6415-11e8-ae84-a860b6301a14"
  ],
  "ack_deadline_seconds": 600
}
```

## Messages modify push:

TODO