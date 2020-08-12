#!/usr/local/bin/python3

import sys
from urllib.parse import unquote

import yaml

unnamed = 0

omit=frozenset([ "#", "components", "schemas", "schema", "parameters"])

def keyFor(path):
    global unnamed
    parts = path.split("/")
    name = parts[-1] + "-" + str(unnamed)
    unnamed = unnamed + 1
    print(parts, file=sys.stderr)
    filtered=filter(lambda x: x not in omit, parts)
    return "-".join(filtered)


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
        key = p.replace("~1", "/")
        d = d[key]
    return (d, parts[-1])


def pathOk(path):
    parts = path.split("/")
    if len(parts) != 4:
        return False
    if parts[0] != "#":
        return False
    if parts[1] != "components":
        return False
    if parts[2] != "schemas":
        return False
    return True


if __name__ == "__main__":
    key = "$ref"
    schema = yaml.load(sys.stdin.read(), Loader=yaml.FullLoader)
    if not "components" in schema:
        schema["components"] = {}
    if not "schemas" in schema["components"]:
        schema["components"]["schemas"] = {}

    # Get list of refs
    refs = list(findKey(schema, key))

    # Get distinct list of schema paths
    schemaPaths = []
    for ref in refs:
        if unquote(ref[key]) not in schemaPaths:
            schemaPaths.append(unquote(ref[key]))

    # Move ref to components section
    for path in sorted(schemaPaths, key=len):
        if pathOk(path):
            continue
        (schemaSection, last) = findSchema(schema, path)
        title = keyFor(path)
        schema["components"]["schemas"][title] = schemaSection[last]
        newPath = "#/components/schemas/%s" % title

        # update all refs for this path
        for ref in refs:
            if unquote(ref[key]) == path:
                ref[key] = newPath

        schemaSection[last] = {"$ref": newPath}

    yaml.dump(schema, sys.stdout)
