# SANextract

## Archived

This project is archived as other awesome community tools that also support a wider range of functionality came up and we switched internal usage to alternatives like httpx.

## About

TLS certificates carry a field called "Subject Alternative Names" containing a list of names the certificate is valid for.
These names are interesting, as they often reveal DNS names of IP addresses, contain subdomains not previously known, or allow to identify services (e.g. self-generated appliance certificates).

SANextract allows to fetch certificates and extract SANs. 
It is tremendously fast (several hundreds of connections per second) and is suitable for bulk operations.
It integrates well into existing tooling as targets are piped into the tool and stdout is designed to be reused by other programs.

## Usage

~~~
./SANextract -h
Usage of ./SANextract:
  -json
        Output JSON.
  -timeout duration
        Connection timeout as duration, e.g. 2s or 800ms (default 2.5s)
  -workers int
        Number of workers. (default 250)
~~~

~~~
$ echo "wikipedia.org" | ./SANextract
*.wikipedia.org
*.wikimedia.org
*.wmfusercontent.org
*.wikimediafoundation.org
*.wiktionary.org
*.wikivoyage.org
*.wikiversity.org
*.wikisource.org
*.wikiquote.org
*.wikinews.org
*.wikidata.org
*.wikibooks.org
wikimedia.org
*.mediawiki.org
wikipedia.org
wikiquote.org
mediawiki.org
wmfusercontent.org
w.wiki
wikimediafoundation.org
wikibooks.org
wiktionary.org
wikivoyage.org
wikidata.org
wikiversity.org
wikisource.org
wikinews.org
*.m.wikipedia.org
*.m.wiktionary.org
*.m.wikivoyage.org
*.m.wikiquote.org
*.m.wikiversity.org
*.m.wikisource.org
*.m.wikimedia.org
*.m.wikinews.org
*.m.wikidata.org
*.m.wikibooks.org
*.planet.wikimedia.org
*.m.mediawiki.org
~~~

~~~
$ cat techgiants.txt
apple.com
microsoft.com
amazon.com
$ time ./SANextract -json < techgiants.txt
{"target":"microsoft.com:443","SANs":["*.oneroute.microsoft.com","oneroute.microsoft.com"]}
{"target":"amazon.com:443","SANs":["amazon.co.uk","uedata.amazon.co.uk","www.amazon.co.uk","origin-www.amazon.co.uk","*.peg.a2z.com","amazon.com","amzn.com","uedata.amazon.com","us.amazon.com","www.amazon.com","www.amzn.com","corporate.amazon.com","buybox.amazon.com","iphone.amazon.com","yp.amazon.com","home.amazon.com","origin-www.amazon.com","origin2-www.amazon.com","buckeye-retail-website.amazon.com","huddles.amazon.com","amazon.de","www.amazon.de","origin-www.amazon.de","amazon.co.jp","amazon.jp","www.amazon.jp","www.amazon.co.jp","origin-www.amazon.co.jp","*.aa.peg.a2z.com","*.ab.peg.a2z.com","*.ac.peg.a2z.com","origin-www.amazon.com.au","www.amazon.com.au","*.bz.peg.a2z.com","amazon.com.au","origin2-www.amazon.co.jp"]}
{"target":"apple.com:443","SANs":["extensions.apple.com","feedback.apple.com","genserv.apple.com","help.apple.com","helposx.apple.com","helpqt.apple.com","images.apple.com","itunespartner.apple.com","prohelp.apple.com","rebate.apple.com","safari-extensions.apple.com","trackingshipment.apple.com","trailers.apple.com","apple.com","www.apple.com"]}
0.01user 0.01system 0:00.45elapsed 4%CPU (0avgtext+0avgdata 11672maxresident)k
0inputs+0outputs (0major+221minor)pagefaults 0swaps
~~~

## Building

Make sure you have go installed. 

- Option A: Clone the repo and run `go build`. 
- Option B: `go get github.com/hvs-consulting/SANextract`

There are no third-party dependencies. 
Since SANextract is written in pure go, you may cross-compile it for all architectures supported by go. 
Tested with Go 1.14.
