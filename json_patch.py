#!/usr/bin/python3
import json
import sys


def pathfinder(o, p):
    for i in p:
        o = o[i]
    return o


def patch_me(patch_file=None):
    global template
    try:
        changed_objects = patch_file["Changed"]
        added_objects = patch_file["Added"]
        removed_objects = patch_file["Removed"]
    except KeyError as e:
        print("Missing a field in patch", e)
        exit(1)

    def _change(tmplt=None, k=None, new=None, old=None, path=None):
        obj = tmplt
        obj = pathfinder(obj, path)
        if obj[k] == old:
            obj[k] = new
        else:
            print("Woah there boyo something happen. old value of of key %s doesnsnt match new shi %s".format(k,old))
            exit(1)

    def _remove(tmplt=None, k=None, path=None):
        obj = tmplt
        obj = pathfinder(obj, path)
        try:
            obj.pop(k)
        except KeyError as e:
            print("Key no in object for removal.", e)
            exit(1)

    def _add(tmplt=None, k=None, v=None, path=None):
        obj = tmplt
        obj = pathfinder(obj, path)
        obj[k] = v

    for item in changed_objects:
        _change(tmplt=template, k=item["Key"], new=item["newValue"], old=item["oldValue"], path=item["Path"])

    for item in added_objects:
        _add(tmplt=template, k=item["Key"], v=item["Value"], path=item["Path"])

    for item in removed_objects:
        _remove(tmplt=template, k=item["Key"], path=item["Path"])

    print(json.dumps(template, sort_keys=True))


if __name__ == "__main__":
    if len(sys.argv) != 3:
        print("Format: json_patch <patch> <template>")
        exit(1)
    patch = json.loads(open(sys.argv[1]).read())
    template = json.loads(open(sys.argv[2]).read())
    patch_me(patch_file=patch)
