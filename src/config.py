# encoding: utf-8


class Config(object):
    # Number of results to fetch from API
    RESULT_COUNT = 9
    # How long to cache results for
    CACHE_MAX_AGE = 20  # seconds
    ICON = "2B939AF4-1A27-4D41-96FE-E75C901C780F.png"
    GOOGLE_ICON = "google.png"
    # Algolia credentials
    ALGOLIA_APP_ID = "BH4D9OD16A"
    ALGOLIA_SEARCH_ONLY_API_KEY = "1d8534f83b9b0cfea8f16498d19fbcab"
    ALGOLIA_SEARCH_INDEX = "material-ui"
