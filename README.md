# crawler

Super simple go web page crawler. Has no manners, is not concurrent.

```
$ go get github.com/adrianduke/crawler
```

```
$ go run cmd/main.go http://adeduke.com
```

```
/
	Static Assets:
		CSS:
			http://adeduke.com/index.xml
			http://adeduke.com/css/bootstrap.min.css
			http://adeduke.com/css/hc.css
			http://netdna.bootstrapcdn.com/font-awesome/4.0.3/css/font-awesome.css
		JS:
			http://adeduke.com/js/jquery-1.10.2.min.js
			http://adeduke.com/js/bootstrap.min.js
			http://adeduke.com/js/bootstrap.js
			http://adeduke.com/js/hc.js
	Internal Links:
		http://adeduke.com/2015/08/how-to-create-a-private-ethereum-chain/
		http://adeduke.com/projects
		http://adeduke.com/
		http://adeduke.com/2015/09/test-driven-development---a-guided-tour/

/
	Static Assets:
		CSS:
			http://adeduke.com/index.xml
			http://adeduke.com/css/bootstrap.min.css
			http://adeduke.com/css/hc.css
			http://netdna.bootstrapcdn.com/font-awesome/4.0.3/css/font-awesome.css
		JS:
			http://adeduke.com/js/jquery-1.10.2.min.js
			http://adeduke.com/js/bootstrap.min.js
			http://adeduke.com/js/bootstrap.js
			http://adeduke.com/js/hc.js
	Internal Links:
		http://adeduke.com/2015/09/test-driven-development---a-guided-tour/
		http://adeduke.com/2015/08/how-to-create-a-private-ethereum-chain/
		http://adeduke.com/projects
		http://adeduke.com/

/projects
	Static Assets:
		CSS:
			http://adeduke.com/css/bootstrap.min.css
			http://adeduke.com/css/hc.css
			http://netdna.bootstrapcdn.com/font-awesome/4.0.3/css/font-awesome.css
		JS:
			http://adeduke.com/js/jquery-1.10.2.min.js
			http://adeduke.com/js/bootstrap.min.js
			http://adeduke.com/js/bootstrap.js
			http://adeduke.com/js/hc.js
	Internal Links:
		http://adeduke.com/
		http://adeduke.com/projects

/2015/09/test-driven-development---a-guided-tour/
	Static Assets:
		CSS:
			http://adeduke.com/css/bootstrap.min.css
			http://adeduke.com/css/hc.css
			http://netdna.bootstrapcdn.com/font-awesome/4.0.3/css/font-awesome.css
		JS:
			http://adeduke.com/js/jquery-1.10.2.min.js
			http://adeduke.com/js/bootstrap.min.js
			http://adeduke.com/js/bootstrap.js
			http://adeduke.com/js/hc.js
	Internal Links:
		http://adeduke.com/
		http://adeduke.com/2015/09/test-driven-development---a-guided-tour/
		http://adeduke.com/2015/08/how-to-create-a-private-ethereum-chain/
		http://adeduke.com/projects

/2015/08/how-to-create-a-private-ethereum-chain/
	Static Assets:
		CSS:
			http://adeduke.com/css/bootstrap.min.css
			http://adeduke.com/css/hc.css
			http://netdna.bootstrapcdn.com/font-awesome/4.0.3/css/font-awesome.css
		JS:
			http://adeduke.com/js/jquery-1.10.2.min.js
			http://adeduke.com/js/bootstrap.min.js
			http://adeduke.com/js/bootstrap.js
			http://adeduke.com/js/hc.js
	Internal Links:
		http://adeduke.com/projects
		http://adeduke.com/
		http://adeduke.com/2015/09/test-driven-development---a-guided-tour/
```