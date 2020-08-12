#!/usr/local/bin/python3

import yaml
from urllib.parse import unquote
import sys
import random
import string

inFilename = 'summary.v1.yaml'
outFilename = 'summary.fixed.v1.yaml'

def keyFor(name):
    letters = string.ascii_lowercase
    return name + '-' + ''.join(random.choice(letters) for i in range(10))

def findKey(d, key):
    if isinstance(d, dict):
        if key in d:
            yield d
        for k in d:
            yield from findKey(d[k], key)
    if isinstance(d, list):
        for val in d:
            yield from findKey(val, key)

def findSchema(d, path):
    parts = path.split("/")
    for p in parts[1:-1]:
        d = d[p.replace("~1", "/")]
    return (d, parts[-1])

if __name__ == "__main__":
    key = "$ref"
    with open(inFilename) as file:
        schema = yaml.load(file, Loader=yaml.FullLoader)

        # Get list of refs
        refs = list(findKey(schema, key))

        # Get distinct list of schema paths
        schemaPaths = []
        for ref in refs:
            print(ref)
            if unquote(ref[key]) not in schemaPaths:
                schemaPaths.append(unquote(ref[key]))

        # Move ref to components section
        for path in schemaPaths:
            (schemaSection, last) = findSchema(schema, path)
            title = keyFor(last)
            if not "components" in schema:
                    schema["components"] = {}
            if not last in schema["components"]:
                schema["components"][last] = {}
            schema["components"][last][title] = schemaSection[last]
            newPath = "#/components/schemas/%s" % title

            # update all refs for this path
            for ref in refs:
                if unquote(ref[key]) == path:
                    ref[key] = newPath

            schemaSection[last] = {"$ref": newPath}

        with open(outFilename, 'w') as outfile:
            documents = yaml.dump(schema, outfile)
