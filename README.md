# mgr [![Build Status](https://secure.travis-ci.org/walf443/mgr.png)](http://travis-ci.org/walf443/mgr)

## What is it?

It's a database migration management tool. This project is under developing yet, API may change in future.

## Feature

### No step files

you should only manage database's schema's sql file.

### No DSL

I'd like to use MySQL's specific feature. So you just write create table statement to file.


## USAGE

```
$ mgr

# specify target file manualy.
$ mgr -before before.sql -after after.sql
$ cat before.sql | mgr -before=stdin -after after.sql
```

SEE ALSO
-----------

 - https://metacpan.org/pod/GitDDL
 - https://github.com/winebarrel/ridgepole
 - https://github.com/r7kamura/scheman.git

Author
--------

Copyright (c) 2015 Keiji Yoshimi

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
