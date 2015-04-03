Baku
----

Tapir is named after a beast from Chinese mythology, known in Japanese mythology as the Baku. This library is a semi replacement for http://tapirgo.com.

What it does?
-------------

Baku indexes one or many specified RSS feeds and posts them to elasticsearch in 15 minutes intervals. It also exposes a search api which returns results in json.

To get started, checkout the `config.json.sample` that comes with this repo.

```
{
  // rss feeds you want to poll every 15 minutes
  "to_be_indexed" : [
    "http://kdlearn.kd.io/feed.xml"
  ],

  // port to run this application under
  "port" : 5200,

  // location of es you want to post results to
  "elasticsearch" : "http://log0.sjc.koding.com"
}
```

Baku will only index the following fields in the feed.

```
<item>
  <title></title>
  <link></link>
  <guid></guid>
  <author></author>
  <description></description>
</item>
```

This was the structure of xml when Baku was written. If you want to add more fields, add them `Item` struct with appropriate tags.

Search
------

By field:

```
/search?title=<your input>

{
  "total"   : 10,
  "results" : [ ...  ]
}
```

All fields:

/search?query=<your input>

{
  "total"   : 100,
  "results" : [ ...  ]
}
```

Pagination:

/search?title=<your input>&page=2

{
  "total"   : 10,
  "results" : [ ...  ]
}
```

Limit:

/search?title=<your input>&limit=10

{
  "total"   : 10,
  "results" : [ ...  ]
}
```

How to run?
-----------

There's a mac and linux binary zipped distributed with this repo. Baku, by default, starts in port 5200, but you can specify which port you want with in the config file.
