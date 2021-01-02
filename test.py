#!/opt/airbnb/python/bin/python3

import requests
import plistlib

output = {}
output["machines"] = 1234
output["plugins"] = ["wee", "poo"]
output["install_type"] = "custom"

# get database type
output["database"] = "notsql"

output["version"] = "abc123"
# plist encode output
post_data = plistlib.dumps(output)
# print(post_data)
response = requests.post(
    "https://version.salopensource.com", data={"data": post_data}, timeout=10
)

print(response.text)
print(vars(response))