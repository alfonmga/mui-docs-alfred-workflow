#!/usr/bin/python
# encoding: utf-8

from __future__ import print_function, unicode_literals, absolute_import

import functools
import re
import sys
from textwrap import wrap
from urllib import quote_plus
from config import Config
from workflow import Workflow3, ICON_INFO
import json


# log
log = None


def cache_key(query):
    """Make filesystem-friendly cache key"""
    key = "{}".format(query)
    key = key.lower()
    key = re.sub(r"[^a-z0-9-_;.]", "-", key)
    key = re.sub(r"-+", "-", key)
    # print("Cache key : {!r} {!r}".format(query, key))
    return key


def handle_result(api_dict):
    """Extract relevant info from API result"""
    result = {}

    for key in {
        "objectID",
        "hierarchy",
        "_snippetResult",
        "anchor",
        "url"
    }:
        result[key] = api_dict[key]

    return result


def search(query=None, limit=Config.RESULT_COUNT):
    import requests
    if query:
        url = "https://{}-dsn.algolia.net/1/indexes/*/queries?x-algolia-application-id={}&x-algolia-api-key={}".format(
            Config.ALGOLIA_APP_ID, Config.ALGOLIA_APP_ID, Config.ALGOLIA_SEARCH_ONLY_API_KEY)
        r = requests.post(
            url, json={"requests": [{
                "indexName": Config.ALGOLIA_SEARCH_INDEX,
                "params": "query={}&hitsPerPage={}&{}".format(query, Config.RESULT_COUNT, "&facetFilters=%5B%22version%3Amaster%22%2C%22language%3Aen%22%5D")
            }],
            })

        response = r.json()
        # print(json.dumps(response))
        if response is not None and "results" in response:
            return response["results"][0]["hits"]

    return []


def main(wf):
    if wf.update_available:
        # Add a notification to top of Script Filter results
        wf.add_item(
            "New version available",
            "Action this item to install the update",
            autocomplete="workflow:update",
            icon=ICON_INFO,
        )

    query = wf.args[0].strip()

    if not query:
        wf.add_item("Search the Material UI docs...")
        wf.send_feedback()
        return 0

    # Parse query into query string and tags
    words = query.split(" ")

    query = []
    for word in words:
        query.append(word)

    query = " ".join(query)

    # print("query: {!r}".format(query))

    key = cache_key(query)

    results = [
        handle_result(result)
        for result in wf.cached_data(
            key, functools.partial(search, query), max_age=Config.CACHE_MAX_AGE
        )
    ]

    log.debug("{} results for {!r},".format(
        len(results), query))

    # Show results
    if not results:
        url = "https://www.google.com/search?q={}".format(
            quote_plus('"Material UI" {}'.format(query))
        )
        wf.add_item(
            "No matching answers found",
            "Shall I try and search Google?",
            valid=True,
            arg=url,
            copytext=url,
            quicklookurl=url,
            icon=Config.GOOGLE_ICON,
        )

    for result in results:
        hierarchies = []
        for key in result["hierarchy"]:
            if result["hierarchy"][key] is not None:
                hierarchies.append(result["hierarchy"][key])

        title = " > ".join(hierarchies)
        subtitle = wrap(result["_snippetResult"]
                        ["content"]["value"], width=75)[0]

        wf.add_item(
            uid=result["objectID"],
            title=title,
            subtitle=subtitle,
            arg=result["url"],
            valid=True,
            # largetext=result["content"],
            copytext=result["url"],
            quicklookurl=result["url"],
            icon=Config.ICON,
        )
        # print(result)

    wf.send_feedback()


if __name__ == "__main__":
    wf = Workflow3(
        update_settings={
            "github_slug": "alfonmga/mui-docs-alfred-workflow", "frequency": 7},
        libraries=['./lib']
    )
    log = wf.logger
    sys.exit(wf.run(main))
