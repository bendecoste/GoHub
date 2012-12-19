import os
import re
import glob
import sys
import urllib
import urllib2
import json
import argparse
from color import color
from file_manager import Manager

_URL = "https://api.github.com"
_GISTS = "/gists"


def load():
    try:
        f = open('.config')
        return "?access_token=" + f.read()
    except IOError:
        # no config
        f = open('.config', 'w+')
        uname = raw_input('Enter github username: ')

        os.system('curl -i -u ' + uname + ' -d \'{"scopes\":[\"public_repo\",\"gist\"]}\' https://api.github.com/authorizations > .setup')

        g = open('.setup')

        keyMatch = '.*"token": "(\w\d)*.*"'

        token = re.findall(r'.*"token": "((\w|\d)*).*', g.read())
        token = token[0][0]

        g.close()
        os.system('rm .setup')

        f.write(token)

        return "?access_token=" + token

_ACCESS = load()
print _ACCESS


def gist():

    parser = argparse.ArgumentParser(description='Command Line Gists')
    parser.add_argument('files', type=str, nargs='+', help='Name of the file you wish to add')
    parser.add_argument('--private', dest='private', nargs='?', const=False, default=True)
    parser.add_argument('--new', dest='new', nargs='?', const=False, default=True)
    parser.add_argument('--dir', dest='dir', nargs='?', type=str, help='Base directory to upload files from (default to cwd')
    parser.add_argument('--desc', dest='desc', nargs='+', type=str, help='Description of the gist')

    args = parser.parse_args()
    uploadFiles(args)

def uploadFiles(args):

    if args.dir is not None:
        os.chdir(args.dir)

    for file in args.files:
        try:
            f = open(file)
            contents = f.read()
        except IOError:
            # file doesn't exist
            os.system('vi ' + file)
            f = open(file)
            contents = f.read()
            os.remove(file)

        manager.addFile(file, contents)

    _request(args)

def _request(args):
    url = _URL + _GISTS + _ACCESS

    if args.desc is None:
        description = ''
    else:
        description = ""
        for d in args.desc:
            description = description + " " + d

    data = json.dumps({'description': description,
                        'public': args.private,
                        'files': manager.files
                        })

    req = urllib2.Request(url=url, data=data)
    res = urllib2.urlopen(req).read()

    _printRes(res)

def _printRes(res):
    res = json.loads(res)
    url = res[u'html_url']
    _id = res[u'id']
    desc = res[u'description']

    print "\n" + color.Blue + url + color.Reset + " (" + url + ")"

    resp = raw_input("Open gist (y/n)")

    if resp == 'y' or resp == 'yes':
        os.system('open ' + url)


manager = Manager()
gist()
