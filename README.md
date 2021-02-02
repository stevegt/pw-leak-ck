# pw-leak-ck

A simple Go tool to discover if any of your passwords have leaked into
any published password lists.  

## Usage

To install the `pw-leak-ck` binary in your path (assuming you trust
the `go get` machinery -- use the `git clone` method below for a more
paranoid approach):

`go get -v github.com/stevegt/pw-leak-ck`

By default, passwords aren't echoed to the screen as you enter them, so a session looks like this:

```
$ pw-leak-ck 
enter passwords, one per line:
> no known leaks
> no known leaks
> leaked 5 times
> leaked 2897638 times
```

For more detail, you can use the `-m` option to show yourself a masked
version of each password -- this helps when you're checking several
and need to remember which passwords need changing: 

```
$ pw-leak-ck -m
enter passwords, one per line:
> a*****7 no known leaks
> l*******f no known leaks
> (*******t leaked 5 times
> a****3 leaked 2897638 times
```

## Security

This tool uses the pwnedpasswords.com API.  The tool uses SHA1 to hash
passwords you enter, and then sends only the first 5 bytes of the hash
to the API server.  I encourage you to examine main.go to confirm
this:

```
cd /tmp
git clone https://github.com/stevegt/pw-leak-ck
cd pw-leak-ck
view main.go
```

...after you're satisfied:

```
go build
./pw-leak-ck
```

The passwords you enter are resident in RAM for a short period of time
during execution, meaning local malware might be able to grab a copy.
But if your local machine has already been exploited, then you have
bigger problems anyway.  

## Disclaimers

If you don't already have a working Go environment, none of the above
will work -- see golang.org to get started, or see pwnedpasswords.com
for tools in other languages.  

I have no affiliation with pwnedpasswords.com.
