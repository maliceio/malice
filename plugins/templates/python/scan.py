#!/usr/bin/env python
# -*- coding: utf-8 -*-
"""
scan.py
~~~~~~~~
{{ description }}.
:copyright: (c) 2016 by {{ creator }}.
:license: MIT
"""

import argparse


def main():
    parser = argparse.ArgumentParser(prog='{{ plugin_name }}')
    parser.add_argument("-v", "--verbose", help="Display verbose output message", action="store_true", required=False)
    parser.add_argument('hash', metavar='MD5', type=str, nargs='+', help='a md5 hash to search for.')
    args = parser.parse_args()

    if args.hash:
        print "Do something awesome with hash: " + args.hash
    return


if __name__ == "__main__":
    try:
        main()
    except Exception as e:
        print "Error: %s" % e
