# mygodb

dbms written in golang.   
I created this db based on this site's db : https://cstack.github.io/db_tutorial/

## Usage
```bash
db > insert 1 gopher example1.com
Executed.
db > select
(1, gopher, example1.com)
Executed.
db > 
```

## TODO
- ~~small sql parser and vm (support insert and select)~~
- ~~datastructure (table, page, row) using array~~
- ~~serializer and deserializer~~
- ~~save data in memory~~
- WAL
- persistence to disk
- b-tree
