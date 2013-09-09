Attempt to decrypt mojibake into readable UTF-8. I can't guarantee that all cases have been covered yet. For now it is up to the user to guess what the source encoding might've been.

Future plans include a heuristic encoding-detection mechanism to bruteforce through a limitable number of encoding combinations and choose the one that looks the least nonsensical.

**Documentation**: [GoPkgDoc](http://go.pkgdoc.org/github.com/moshee/mojibake)

#### TODO

- performance improvements
- find and add more mis-encoding cases
- more tests (like for CP936)
