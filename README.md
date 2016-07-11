Project Creator
==

A simple lambda function/ set of functions to create adobe premiere pro project files on the fly.

The lambda function takes four arguments, as below:

```json
{
    "uuid": "2d3f8685-f3b9-49b5-824d-1215a0580a4a",
    "name": "flying-cars",
    "srcBucket": "my-source-bucket",
    "dstBucket": "my-destination-bucket"
}
```

We take some empty files from `my-src-bucket`, create a directory structure in `my-destination-bucket` and pop correctly named/ formatted `{pr,pl}proj` files in.

The assumption is that these variables come from an AWS api gateway using a mix of stage variables and request variables, as per:

Body Mapping Template
--

```json
#set($inputRoot = $input.path('$'))
{
  "uuid": "$inputRoot.uuid",
  "name": "$inputRoot.name",
  "srcBucket": "$stageVariables.srcBucket",
  "dstBucket": "$stageVariables.dstBucket"
  }
```

Request Payload
--

```json
{
    "uuid": "dc2216c0-98d5-4757-938c-f1b46ea4c855",
    "name": "a-test-project"
}
```

Stage Variables
--

```json
{
    "srcBucket": "my-source-bucket",
    "dstBucket": "my-destination-bucket"
}
```

(etc.)
