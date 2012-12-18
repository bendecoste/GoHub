import os
import glob
import sys
import urllib
import urllib2

_URL = "https://api.github.com"
_GISTS = "/gists"
_ACCESS = "?access_token=10520ec3d313a78e387fe9cf54da6f2af9acdbf4"
def gist():
    arg = sys.argv[1]

    # No file specified
    if arg is None:
        _noFile()
    else:
        _files(arg)

def _files(arg):
    os.mkdir("./.gists")

    os.system("vi ./.gists/" + arg)

    f = open("./.gists/" + arg, 'r')
    contents = f.read()

    os.unlink("./.gists/" + arg)
    os.rmdir("./.gists")

    # TODO: pub/private
    _request(arg, contents)

def _request(arg, contents):

    url = _URL + _GISTS + _ACCESS
    print url
    data = urllib.urlencode({'description': 'API created gist!', 'public': 'true', 'files': { arg: { 'content': contents}}})
    req = urllib2.urlopen(url, data)
    #res = urllib2.urlopen(req)

    print res.read()

gist()
  
