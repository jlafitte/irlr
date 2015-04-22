# IIS Redirect Loop Rectifier
Tries to watch the http status code and add a comment to web.config
if it comes up with a 30x error.  This attempts to be a simple work around for
[this bug](http://forums.iis.net/t/1178961.aspx?Redirect+loop+bug).
