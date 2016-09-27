Project Creator
==

`project-creator` is a simple dockerised service which generates empty .{pl,pr}project files, file structures, and s3 goodies for Financial Times video projects/

The api takes a payload as per:

```json
{
    "uuid": "2d3f8685-f3b9-49b5-824d-1215a0580a4a",
    "name": "flying-cars",
}
```
